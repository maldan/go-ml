package ms_handler

import (
	"encoding/json"
	ms_error "github.com/maldan/go-ml/server/error"
	ms_response "github.com/maldan/go-ml/server/response"
	ml_string "github.com/maldan/go-ml/util/string"
	"reflect"
	"strconv"

	"io/ioutil"
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

	// Has 2 arg
	if functionType.NumIn() == 3 {
		// Get last arg
		arg := reflect.New(functionType.In(2)).Interface()

		// Get type
		argType := reflect.TypeOf(arg).Elem()

		// Create new value
		argValue := reflect.New(argType)

		b, _ := json.Marshal(params)
		json.Unmarshal(b, argValue.Elem().Addr().Interface())

		// Is struct
		if argType.Kind() == reflect.Struct {
			return virtualCall(method, controller, &Context{}, argValue.Elem().Interface())
		}
	}

	return reflect.ValueOf("")
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

	// Read body
	bodyBytes, _ := ioutil.ReadAll(args.Request.Body)
	if len(bodyBytes) > 0 {
		// Parse json body and
		jsonMap := map[string]any{}
		err := json.Unmarshal(bodyBytes, &jsonMap)
		ms_error.FatalIfError(err)

		// Collect params
		for key, element := range jsonMap {
			params[key] = element
		}
	}

	// Get controller
	path := strings.Split(strings.Replace(args.Request.URL.Path, args.Route, "", 1), "/")
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
		ctrName := ml_string.UnTitle(cc[len(cc)-1])
		if ctrName == controllerName {
			controller = c
			break
		}
	}

	// Find method
	var method any = nil
	controllerType := reflect.TypeOf(controller)
	for i := 0; i < controllerType.NumMethod(); i++ {
		if controllerType.Method(i).Name == methodName {
			method = controllerType.Method(i)
			break
		}
	}

	// Call method
	value := callMethod(method.(reflect.Method), controller, params).Interface()

	switch value.(type) {
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
