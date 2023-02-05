package ms

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	ms_handler "github.com/maldan/go-ml/server/core/handler"
	ms_error "github.com/maldan/go-ml/server/error"
	ml_slice "github.com/maldan/go-ml/slice"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

//go:embed panel_frontend/dist/*
var panelFs embed.FS

func HandleError(args *ms_handler.Args) {
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

		/*for i := 0; i < 10; i++ {
			p, f, l, ok := runtime.Caller(i)
			if ok {
				fmt.Printf("%v %v:%v\n", p, f, l)
			}
		}
		*/

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

func getHandler(url string, routers []ms_handler.RouteHandler) (string, ms_handler.Handler) {
	for i := 0; i < len(routers); i++ {
		if strings.HasPrefix(url, routers[i].Path) {
			return routers[i].Path, routers[i].Handler
		}
	}

	return "", ms_handler.Undefined{}
}

func injectDebug(config *Config) {
	// Add debug controller
	config.Router = ml_slice.Prepend(config.Router, []ms_handler.RouteHandler{
		{
			Path: "/debug/panel",
			Handler: ms_handler.EmbedFS{
				Root: "panel_frontend/dist",
				Fs:   panelFs,
			},
			/*Handler: ms_handler.API{
				ControllerList: []any{ms_panel.Panel{}},
			},*/
		},
	})

	/*ms_panel.Html = PanelHtml
	ms_panel.Css = PanelCss
	ms_panel.Js = PanelJs*/
}

func Start(config Config) {
	injectDebug(&config)

	// Entry point
	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		// Prepare args
		args := ms_handler.Args{Response: response, Request: request}
		defer HandleError(&args)

		// Disable cors for all queries
		DisableCors(response)

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
