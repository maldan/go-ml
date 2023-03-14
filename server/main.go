package ms

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	ms_handler "github.com/maldan/go-ml/server/core/handler"
	ms_error "github.com/maldan/go-ml/server/error"
	ms_log "github.com/maldan/go-ml/server/log"
	ms_panel "github.com/maldan/go-ml/server/panel"
	ml_slice "github.com/maldan/go-ml/util/slice"
	"net/http"
	"runtime"
	"strings"
)

//go:embed panel_frontend/dist/*
var panelFs embed.FS

func endOfRequest(args *ms_handler.Args) {
	err := recover()
	if err == nil {
		return
	}

	// Set error output as json
	args.Response.Header().Add("Content-Type", "application/json")

	switch e := err.(type) {
	case ms_error.Error:
		args.Response.WriteHeader(e.Code)
		e.EndPoint = args.Path
		message, _ := json.Marshal(e)
		args.Response.Write(message)
		ms_log.Log("request error", e)
	default:
		stackInfo := make([]string, 0, 10)
		for i := 0; i < 10; i++ {
			_, f, l, ok := runtime.Caller(i + 1)
			if ok {
				// Skip system libs, no points in it
				if strings.Contains(f, "/local/go/src/") || strings.Contains(f, "/Program Files/Go/src/") {
					continue
				}
				stackInfo = append(stackInfo, fmt.Sprintf("%v:%v", f, l))
			}
		}

		args.Response.WriteHeader(500)

		ee := ms_error.Error{
			Type:        "unknown",
			Description: fmt.Sprintf("%v", e),
			Stack:       stackInfo,
			EndPoint:    args.Path,
		}
		message, _ := json.Marshal(ee)
		args.Response.Write(message)
		ms_log.Log("request error", ee)
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
		},
		{
			Path: "/debug/api",
			Handler: ms_handler.API{
				ControllerList: []any{
					ms_panel.Panel{
						HasLogTab: config.Panel.HasLogTab,
					},
					ms_panel.Log{
						Path: config.LogFile,
					},
				},
			},
		},
	})
}

/*func globalPanicHandler() {
	err := recover()
	if err == nil {
		return
	}
	Log("global panic", err)
}*/

func Start(config Config) {
	// defer globalPanicHandler()
	injectDebug(&config)

	// Entry point
	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		// Prepare args
		args := ms_handler.Args{Response: response, Request: request}
		defer endOfRequest(&args)

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
	})

	// Start logger
	if config.LogFile != "" {
		config.Panel.HasLogTab = false
		ms_log.Init(config.LogFile)
	}

	// fmt.Printf("%v\n", "FUCK")
	ms_log.Log("info", fmt.Sprintf("Mega Server Starts at host %v", config.Host))

	if config.TLS.Enabled {
		err := http.ListenAndServeTLS(config.Host, config.TLS.CertFile, config.TLS.KeyFile, nil)
		ms_error.FatalIfError(err)
	} else {
		err := http.ListenAndServe(config.Host, nil)
		ms_error.FatalIfError(err)
	}
}
