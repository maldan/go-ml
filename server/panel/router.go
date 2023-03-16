package ms_panel

import (
	ms_handler "github.com/maldan/go-ml/server/core/handler"
	ms_error "github.com/maldan/go-ml/server/error"
	ml_slice "github.com/maldan/go-ml/util/slice"
	ml_string "github.com/maldan/go-ml/util/string"
	"reflect"
	"strings"
)

type Router struct {
	List []ms_handler.RouteHandler
}

/*type MethodInput struct {
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Kind      string         `json:"kind"`
	FieldList []*MethodInput `json:"fieldList"`
}

type MethodInfo struct {
	Uid string `json:"uid"`

	FullPath string `json:"fullPath"`
	Url      string `json:"url"`

	Controller  string `json:"controller"`
	HttpMethod  string `json:"httpMethod"`
	Name        string `json:"name"`
	InputMethod string `json:"inputMethod"`

	Input *MethodInput `json:"input"`
}*/

type RouterInfo struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

type ArgsRouterMethodList struct {
	Path       string
	Controller string
}

/*type MethodArgument struct {
	Type string `json:"type"`
}*/

type MethodInfo struct {
	Url        string   `json:"url"`
	HttpMethod string   `json:"httpMethod"`
	Args       []string `json:"args"`
	Return     []string `json:"return"`
}

func (r Router) GetList() []RouterInfo {
	out := make([]RouterInfo, 0)

	for i := 0; i < len(r.List); i++ {
		out = append(out, RouterInfo{
			Path: r.List[i].Path,
			Type: reflect.TypeOf(r.List[i].Handler).Name(),
		})
	}

	return out
}

func (r Router) GetControllerList(args ArgsRouterMethodList) []string {
	rr, ok := ml_slice.Find(r.List, func(x *ms_handler.RouteHandler) bool {
		return x.Path == args.Path
	})
	ms_error.FatalIf(!ok, ms_error.Error{Code: 404})
	xx := rr.Handler.(ms_handler.API)

	out := make([]string, 0)
	for i := 0; i < len(xx.ControllerList); i++ {
		typeOf := reflect.TypeOf(xx.ControllerList[i])
		controllerName := ml_string.UnTitle(typeOf.Name())
		out = append(out, controllerName)
	}
	return out
}

func (r Router) GetMethodList(args ArgsRouterMethodList) []MethodInfo {
	rr, ok := ml_slice.Find(r.List, func(x *ms_handler.RouteHandler) bool {
		return x.Path == args.Path
	})
	ms_error.FatalIf(!ok, ms_error.Error{Code: 404})
	xx := rr.Handler.(ms_handler.API)

	out := make([]MethodInfo, 0)
	for i := 0; i < len(xx.ControllerList); i++ {
		typeOf := reflect.TypeOf(xx.ControllerList[i])
		controllerName := ml_string.UnTitle(typeOf.Name())

		if controllerName != args.Controller {
			continue
		}

		ml := typeOf.NumMethod()
		for j := 0; j < ml; j++ {
			// Name
			method := typeOf.Method(j)
			methodName := method.Name
			httpMethod := ""
			if strings.Contains(methodName, "Get") {
				methodName = strings.Replace(methodName, "Get", "", 1)
				httpMethod = "GET"
			}
			if strings.Contains(methodName, "Delete") {
				methodName = strings.Replace(methodName, "Delete", "", 1)
				httpMethod = "DELETE"
			}
			if strings.Contains(methodName, "Post") {
				methodName = strings.Replace(methodName, "Post", "", 1)
				httpMethod = "POST"
			}
			if strings.Contains(methodName, "Patch") {
				methodName = strings.Replace(methodName, "Patch", "", 1)
				httpMethod = "PATCH"
			}
			if strings.Contains(methodName, "Put") {
				methodName = strings.Replace(methodName, "Put", "", 1)
				httpMethod = "PUT"
			}

			// Fill args
			methodArgs := make([]string, 0)
			if method.Type.NumIn() > 1 {
				for k := 0; k < method.Type.NumIn()-1; k++ {
					methodArgs = append(methodArgs, method.Type.In(1+k).String())
				}
			}
			methodReturn := make([]string, 0)
			if method.Type.NumIn() > 1 {
				for k := 0; k < method.Type.NumOut(); k++ {
					methodReturn = append(methodReturn, method.Type.Out(k).String())
				}
			}

			// Append
			out = append(out, MethodInfo{
				HttpMethod: httpMethod,
				Args:       methodArgs,
				Return:     methodReturn,
				Url:        args.Path + "/" + controllerName + "/" + ml_string.UnTitle(methodName),
			})
		}
	}
	return out
}
