package ms

import (
	"encoding/json"
	"fmt"
	ms_config "github.com/maldan/go-ml/server/config"
	ms_core "github.com/maldan/go-ml/server/core"
	"github.com/maldan/go-ml/server/core/handler"
	ms_error "github.com/maldan/go-ml/server/error"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func HandleError(args *ms_config.HandlerArgs) {
	err := recover()
	if err == nil {
		return
	}

	// Set error output as json
	args.Response.Header().Add("Content-Type", "application/json")

	switch e := err.(type) {
	case ms_error.Error:
		args.Response.WriteHeader(e.Code)
		message, _ := json.Marshal(e)
		args.Response.Write(message)
		/*if args.DebugMode {
			rapi_debug.Log(args.Id).SetError(e)
			rapi_debug.Log(args.Id).SetArgs(args.MethodArgs)
		}*/
	default:
		_, file, line, _ := runtime.Caller(3)

		for i := 0; i < 10; i++ {
			p, f, l, ok := runtime.Caller(i)
			if ok {
				fmt.Printf("%v %v:%v\n", p, f, l)
			}
		}

		args.Response.WriteHeader(500)
		// fmt.Println(string(debug.Stack()))
		ee := ms_error.Error{
			Code:        500,
			Type:        "unknown",
			Description: fmt.Sprintf("%v", e),
			Line:        line,
			File:        file,
			// Stack:       string(debug.Stack()),
			Created: time.Now(),
		}
		message, _ := json.Marshal(ee)
		args.Response.Write(message)
		/*if args.DebugMode {
			rapi_debug.Log(args.Id).SetError(ee)
			rapi_debug.Log(args.Id).SetArgs(args.MethodArgs)
		}*/
	}
}

func getHandler(url string, routers []ms_config.RouteHandler) (string, ms_config.Handler) {
	for i := 0; i < len(routers); i++ {
		if strings.HasPrefix(url, routers[i].Path) {
			return routers[i].Path, routers[i].Handler
		}
	}

	return "", handler.Undefined{}
}

func Start(config ms_config.Config) {
	// Entry point
	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		// Prepare args
		args := ms_config.HandlerArgs{Response: response, Request: request}
		defer HandleError(&args)

		// Disable cors for all queries
		ms_core.DisableCors(response)

		// Fuck options
		if request.Method == "OPTIONS" {
			response.WriteHeader(200)
			return
		}

		// Get handler
		route, h := getHandler(request.URL.Path, config.Router)
		args.Path = request.URL.Path
		args.Route = route

		// Handle
		h.Handle(args)

		// Done
		fmt.Printf("%+v\n", args)
	})

	log.Printf("Mega Server Starts at host %v\n", config.Host)

	err := http.ListenAndServe(config.Host, nil)
	ms_error.FatalIfError(err)
}
