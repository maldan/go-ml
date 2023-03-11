package ms_handler

import (
	"encoding/json"
	"fmt"
	ms_error "github.com/maldan/go-ml/server/error"
	ms_response "github.com/maldan/go-ml/server/response"
	ml_file "github.com/maldan/go-ml/util/io/fs/file"
	ml_string "github.com/maldan/go-ml/util/string"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"strings"
)

type API struct {
	ControllerList []any
}

func virtualCall(fn reflect.Method, args ...any) reflect.Value {
	function := reflect.ValueOf(fn.Func.Interface())

	// Prepare args
	in := make([]reflect.Value, len(args))
	for i, v := range args {
		in[i] = reflect.ValueOf(v)
	}

	// Call and get response
	r := function.Call(in)
	if len(r) > 0 {
		return r[0]
	}

	return reflect.ValueOf("")
}

func callMethod(method reflect.Method, controller any, params map[string]any) reflect.Value {
	functionType := reflect.TypeOf(method.Func.Interface())

	// Has 0 arg
	if functionType.NumIn() == 1 {
		return virtualCall(method, controller)
	}

	// Has 1 arg
	if functionType.NumIn() == 2 {
		// Get last arg
		arg := reflect.New(functionType.In(1)).Interface()

		// Get type
		argType := reflect.TypeOf(arg).Elem()

		// Create new value
		argValue := reflect.New(argType)

		// Copy json to that value
		b, _ := json.Marshal(params)
		json.Unmarshal(b, argValue.Elem().Addr().Interface())

		// Is struct
		if argType.Kind() == reflect.Struct {
			return virtualCall(method, controller, argValue.Elem().Interface())
		} else {
			panic("Argument must be struct type")
		}
	}

	// Has 2 arg
	if functionType.NumIn() == 3 {
		// Get last arg
		arg := reflect.New(functionType.In(2)).Interface()

		// Get type
		argType := reflect.TypeOf(arg).Elem()

		// Create new value
		argValue := reflect.New(argType)

		// Copy json to that value
		b, _ := json.Marshal(params)
		json.Unmarshal(b, argValue.Elem().Addr().Interface())

		// Is struct
		if argType.Kind() == reflect.Struct {
			return virtualCall(method, controller, &Context{}, argValue.Elem().Interface())
		} else {
			panic("Argument must be struct type")
		}
	}

	panic("Method not found")
}

func (a API) Handle(args Args) {
	// Get authorization
	authorization := args.Request.Header.Get("Authorization")
	authorization = strings.Replace(authorization, "Token ", "", 1)

	// Collect params
	params := map[string]any{
		"accessToken": authorization,
	}

	// Read url params
	for key, element := range args.Request.URL.Query() {
		num, err := strconv.Atoi(element[0])
		if err == nil {
			params[key] = num
		} else {
			params[key] = element[0]
		}
	}

	// Parse multipart body
	if strings.Contains(args.Request.Header.Get("Content-Type"), "multipart/form-data") {
		// Parse multipart body and collect params
		err := args.Request.ParseMultipartForm(0)
		ms_error.FatalIfError(err)
		for key, element := range args.Request.MultipartForm.Value {
			params[key] = element[0]
		}

		// Collect files
		if len(args.Request.MultipartForm.File) > 0 {
			for kk, fileHeaders := range args.Request.MultipartForm.File {
				for _, header := range fileHeaders {
					params[kk] = ml_file.NewWithMime(
						reflect.ValueOf(header).Elem().FieldByName("tmpfile").String(),
						header.Header.Get("Content-Type"),
					)
				}
			}
		}
	} else {
		// Read body as json
		bodyBytes, _ := io.ReadAll(args.Request.Body)
		err := args.Request.Body.Close()
		ms_error.FatalIfError(err)

		if len(bodyBytes) > 0 {
			// Parse json body and
			jsonMap := map[string]any{}
			err2 := json.Unmarshal(bodyBytes, &jsonMap)
			ms_error.FatalIfError(err2)

			// Collect params
			for key, element := range jsonMap {
				params[key] = element
			}
		}
	}

	// Get controller
	path := strings.Split(strings.Replace(args.Request.URL.Path, args.Route, "", 1), "/")
	if len(path) <= 1 {
		panic("controller not specified")
	}
	controllerName := path[1]

	// Get method
	methodName := ""
	if len(path) > 2 {
		methodName = path[2]
	}
	if methodName == "" {
		methodName = "Index"
	}
	methodName = ml_string.Title(strings.ToLower(args.Request.Method)) + ml_string.Title(methodName)

	// Find controller
	var controller any = nil
	for _, c := range a.ControllerList {
		cc := strings.Split(reflect.TypeOf(c).String(), ".")
		// ctrName := ml_string.UnTitle(cc[len(cc)-1])

		if strings.ToLower(cc[len(cc)-1]) == strings.ToLower(controllerName) {
			controller = c
			break
		}
	}
	if controller == nil {
		panic(fmt.Sprintf("controller %v not found", controllerName))
	}

	// Find method
	var method any = nil
	controllerType := reflect.TypeOf(controller)
	for i := 0; i < controllerType.NumMethod(); i++ {
		if strings.ToLower(controllerType.Method(i).Name) == strings.ToLower(methodName) {
			method = controllerType.Method(i)
			break
		}
	}
	if method == nil {
		panic(fmt.Sprintf("method %v not found", methodName))
	}

	// Call method
	reflectValue := callMethod(method.(reflect.Method), controller, params)
	value := reflectValue.Interface()

	switch value.(type) {
	case ms_response.File:
		v := value.(ms_response.File)

		// Copy headers
		for k, v := range v.Headers {
			args.Response.Header().Add(k, v)
		}

		http.ServeFile(args.Response, args.Request, v.Path)
		break
	case ms_response.Custom:
		v := value.(ms_response.Custom)

		// Copy headers
		for k, v := range v.Headers {
			args.Response.Header().Add(k, v)
		}

		_, err := args.Response.Write(v.Body)
		ms_error.FatalIfError(err)
		break
	default:
		// Check to response method
		tr, ok := reflect.TypeOf(value).MethodByName("ToResponse")
		if ok {
			// Call
			ret := tr.Func.Call([]reflect.Value{reflect.ValueOf(value)})
			if len(ret) > 0 {
				value = ret[0].Interface()
			}
		}

		// Convert to json
		data, err := json.Marshal(&value)
		ms_error.FatalIfError(err)

		// Write response
		args.Response.Header().Add("Content-Type", "application/json")
		_, err = args.Response.Write(data)
		ms_error.FatalIfError(err)
		break
	}
}
