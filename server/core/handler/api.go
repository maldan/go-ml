package ms_handler

import (
	"encoding/json"
	"fmt"
	ms_error "github.com/maldan/go-ml/server/error"
	ms_response "github.com/maldan/go-ml/server/response"
	ml_file "github.com/maldan/go-ml/util/io/fs/file"
	ml_slice "github.com/maldan/go-ml/util/slice"
	ml_string "github.com/maldan/go-ml/util/string"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type API struct {
	ControllerList []any
	Middleware     func(args *Args)
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

func mapHeader(context *Context, argValue reflect.Value) {
	argPointer := argValue.Elem().Addr().Interface()
	valueOf := reflect.ValueOf(argPointer).Elem()
	typeOf := reflect.TypeOf(argPointer).Elem()
	allowedMap := []string{"Authorization"}

	// Get fields
	fieldMap := map[string]reflect.Value{}
	for i := 0; i < valueOf.NumField(); i++ {
		f := typeOf.Field(i)
		fieldMap[f.Name] = valueOf.Field(i)
		if f.Tag.Get("json") != "" {
			fieldMap[f.Tag.Get("json")] = valueOf.Field(i)
		}
	}

	for k, v := range context.Request.Header {
		// Map only allowed
		if !ml_slice.Includes(allowedMap, k) {
			continue
		}

		// Get field
		field, ok := fieldMap[k]
		if !ok {
			continue
		}

		// Map
		if field.CanSet() {
			if field.Kind() == reflect.String {
				field.SetString(v[0])
			}
		}
	}
}

func mapGet(context *Context, argValue reflect.Value) {
	argPointer := argValue.Elem().Addr().Interface()
	valueOf := reflect.ValueOf(argPointer).Elem()
	typeOf := reflect.TypeOf(argPointer).Elem()

	// Get fields
	fieldMap := map[string]reflect.Value{}
	for i := 0; i < valueOf.NumField(); i++ {
		f := typeOf.Field(i)
		fieldMap[f.Name] = valueOf.Field(i)
		if f.Tag.Get("json") != "" {
			fieldMap[f.Tag.Get("json")] = valueOf.Field(i)
		}
	}

	for k, v := range context.Request.URL.Query() {
		// Map strings
		field, ok := fieldMap[k]
		if !ok {
			continue
		}

		// Set
		if field.CanSet() {
			if field.Kind() == reflect.String {
				field.SetString(v[0])
			}
			if field.Kind() == reflect.Int {
				i, err := strconv.ParseInt(v[0], 10, 64)
				if err != nil {
					fmt.Printf("MS GET Parser [Warning]: %v\n", err)
					// panic(err)
				}
				field.SetInt(i)
			}
		}
	}
}

func mapJson(context *Context, argValue reflect.Value) {
	// Read all data
	bodyBytes, err := io.ReadAll(context.Request.Body)
	if err != nil {
		panic(err)
	}

	// Close
	err = context.Request.Body.Close()
	if err != nil {
		panic(err)
	}

	// Unpack json to struct
	err = json.Unmarshal(bodyBytes, argValue.Elem().Addr().Interface())
	if err != nil {
		panic(err)
	}
}

func mapFormData(context *Context, argValue reflect.Value) {
	argPointer := argValue.Elem().Addr().Interface()
	valueOf := reflect.ValueOf(argPointer).Elem()
	typeOf := reflect.TypeOf(argPointer).Elem()

	// Parse multipart
	err := context.Request.ParseMultipartForm(0)
	if err != nil {
		panic(err)
	}

	// Get fields
	fieldMap := map[string]reflect.Value{}
	fieldFileMap := map[string]reflect.Value{}

	for i := 0; i < valueOf.NumField(); i++ {
		f := typeOf.Field(i)

		if f.Type.Name() == "File" {
			fieldFileMap[f.Name] = valueOf.Field(i)
			if f.Tag.Get("json") != "" {
				fieldFileMap[f.Tag.Get("json")] = valueOf.Field(i)
			}
		} else {
			fieldMap[f.Name] = valueOf.Field(i)
			if f.Tag.Get("json") != "" {
				fieldMap[f.Tag.Get("json")] = valueOf.Field(i)
			}
		}
	}

	// Parse values
	for k, v := range context.Request.MultipartForm.Value {
		// Map strings
		field, ok := fieldMap[k]
		if !ok {
			continue
		}

		// Set
		if field.CanSet() {
			if field.Kind() == reflect.String {
				field.SetString(v[0])
			}
		}
	}

	// Collect files
	if len(context.Request.MultipartForm.File) > 0 {
		for k, fileHeaders := range context.Request.MultipartForm.File {
			for _, header := range fileHeaders {
				// Map strings
				field, ok := fieldFileMap[k]
				if !ok {
					continue
				}

				if field.CanSet() {
					file := ml_file.NewWithMime(
						reflect.ValueOf(header).Elem().FieldByName("tmpfile").String(),
						header.Header.Get("Content-Type"),
					)
					file.SetAttribute("FileName", header.Filename)
					field.Set(reflect.ValueOf(*file))
				}
			}
		}
	}
}

func callMethod(
	method reflect.Method,
	context *Context,
	controller any,
	params map[string]any,
) reflect.Value {
	functionType := reflect.TypeOf(method.Func.Interface())

	// Has 0 arg
	if functionType.NumIn() == 1 {
		return virtualCall(method, controller)
	}

	// Has 1 arg
	if functionType.NumIn() == 2 {
		// Get struct arg
		arg := reflect.New(functionType.In(1)).Interface()

		// Get type
		argType := reflect.TypeOf(arg).Elem()

		// Create new value
		argValue := reflect.New(argType)

		// Args is context
		if argType.Kind() == reflect.Pointer {
			return virtualCall(method, controller, context)
		}

		// If received data as json
		types := []string{"application/json", "text/plain"}
		if ml_slice.Includes(types, context.Request.Header.Get("Content-Type")) {
			mapJson(context, argValue)
		}

		// If multipart data
		if strings.Contains(context.Request.Header.Get("Content-Type"), "multipart/form-data") {
			mapFormData(context, argValue)
		}

		// Map get params
		mapHeader(context, argValue)

		// Map header
		mapGet(context, argValue)

		// Copy json to that value
		/*b, err := json.Marshal(params)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(b, argValue.Elem().Addr().Interface())
		if err != nil {
			panic(err)
		}*/

		// Is struct
		if argType.Kind() == reflect.Struct {
			return virtualCall(method, controller, argValue.Elem().Interface())
		} else {
			panic("Argument must be struct type")
		}
	}

	// Has 2 arg
	/*if functionType.NumIn() == 3 {
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
			return virtualCall(method, controller, context, argValue.Elem().Interface())
		} else {
			panic("Argument must be struct type")
		}
	}*/

	panic("Method not found")
}

func (a API) Handle(args *Args) {
	// Get authorization
	/*authorization := args.Request.Header.Get("Authorization")
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
	}*/

	// Read body
	/*bodyBytes, _ := io.ReadAll(args.Request.Body)
	err := args.Request.Body.Close()
	ms_error.FatalIfError(err)
	args.Body = bodyBytes*/

	// Call middleware
	if a.Middleware != nil {
		a.Middleware(args)
	}

	// Parse multipart body
	if strings.Contains(args.Request.Header.Get("Content-Type"), "multipart/form-data") {
		// Parse multipart body and collect params
		/*err := args.Request.ParseMultipartForm(0)
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
		}*/
	} else {
		// Read body as json
		/*if len(args.Body) > 0 {
			// Parse json body and
			jsonMap := map[string]any{}
			err2 := json.Unmarshal(args.Body, &jsonMap)
			ms_error.FatalIfError(err2)

			// Collect params
			for key, element := range jsonMap {
				params[key] = element
			}
		}*/
	}

	// Get controller name
	path := strings.Split(strings.Replace(args.Request.URL.Path, args.Route, "", 1), "/")
	if len(path) <= 1 {
		panic("controller not specified")
	}
	controllerName := path[1]

	// Get method name
	methodName := ""
	if len(path) > 2 {
		methodName = path[2]
	}
	if methodName == "" {
		methodName = "Index"
	}
	methodName = ml_string.Title(strings.ToLower(args.Request.Method)) + ml_string.Title(methodName)

	// Find controller struct
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
		ms_error.Fatal(ms_error.Error{
			Code:        404,
			Description: fmt.Sprintf("controller %v not found", controllerName),
		})
	}

	// Find method function in controller struct
	var method any = nil
	controllerType := reflect.TypeOf(controller)
	for i := 0; i < controllerType.NumMethod(); i++ {
		if strings.ToLower(controllerType.Method(i).Name) == strings.ToLower(methodName) {
			method = controllerType.Method(i)
			break
		}
	}
	if method == nil {
		ms_error.Fatal(ms_error.Error{
			Code:        404,
			Description: fmt.Sprintf("method %v not found", methodName),
		})
	}

	// Get client ip
	remoteIp := ""
	if len(args.Request.Header["X-Forwarded-For"]) > 0 {
		remoteIp = args.Request.Header["X-Forwarded-For"][0]
	}

	// Call method
	reflectValue := callMethod(method.(reflect.Method), &Context{
		// AccessToken: authorization,
		Response: args.Response,
		Request:  args.Request,
		RemoteIP: remoteIp,
	}, controller, nil)
	value := reflectValue.Interface()

	switch value.(type) {
	case ms_response.File:
		v := value.(ms_response.File)

		// Copy headers
		for k, v2 := range v.Headers {
			args.Response.AddHeader(k, v2)
		}

		http.ServeFile(args.Response, args.Request, v.Path)
		break
	case ms_response.Custom:
		v := value.(ms_response.Custom)

		// Copy headers
		for k, v2 := range v.Headers {
			args.Response.AddHeader(k, v2)
		}

		if v.Reader != nil {
			_, err2 := io.Copy(args.Response, v.Reader)
			ms_error.FatalIfError(err2)
		} else {
			_, err2 := args.Response.Write(v.Body)
			ms_error.FatalIfError(err2)
		}

		break
	case nil:
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
		data, err2 := json.Marshal(&value)
		ms_error.FatalIfError(err2)

		// Write response
		args.Response.AddHeader("Content-Type", "application/json")
		_, err2 = args.Response.Write(data)
		ms_error.FatalIfError(err2)
		break
	}
}
