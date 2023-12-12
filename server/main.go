package ms

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	ms_handler "github.com/maldan/go-ml/server/core/handler"
	ms_error "github.com/maldan/go-ml/server/error"

	ml_slice "github.com/maldan/go-ml/util/slice"
	"net/http"
	"runtime"
	"strings"
)

/*
//go:embed panel_frontend/dist/*
var panelFs embed.FS
*/

func endOfRequest(args *ms_handler.Args) {
	err := recover()
	if err == nil {
		// ms_log.LogRequest(args)
		return
	}

	// Set error output as json
	args.Response.AddHeader("Content-Type", "application/json")

	switch e := err.(type) {
	case ms_error.Error:
		args.Response.WriteHeader(e.Code)
		e.EndPoint = args.Path
		message, _ := json.Marshal(e)
		args.Response.Write(message)
		// ms_log.Log("request error", e)
		// ms_log.LogRequest(args)
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
		//ms_log.Log("request error", ee)
		// ms_log.LogRequest(args)
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

/*func initDb(config *Config) {
	if config.DataBase.DataBase == nil {
		return
	}

	if *config.DataBase.DataBase == nil {
		*config.DataBase.DataBase = map[string]*mdb.DataTable{}
	}

	for i := 0; i < len(config.DataBase.TableList); i++ {
		table := config.DataBase.TableList[i]

		(*config.DataBase.DataBase)[table.Name] = mdb.New(
			config.DataBase.Path,
			table.Name,
			table.Type,
			&gosn_driver.Container{},
		)
	}
}*/

func injectDebug(config *Config) {
	// Add debug controller
	config.Router = ml_slice.Prepend(config.Router, []ms_handler.RouteHandler{
		/*{
			Path: "/debug/panel",
			Handler: ms_handler.EmbedFS{
				Root: "panel_frontend/dist",
				Fs:   panelFs,
			},
		},*/
		{
			Path: "/debug/api",
			Handler: ms_handler.API{
				ControllerList: []any{
					/*ms_panel.Panel{
						HasLogTab:      config.Debug.UseLogs,
						HasRequestLogs: config.Debug.UseRequestLogs,
					},*/
					/*ms_panel.Log{
						Path: config.LogFile,
					},*/
					/*ms_panel.Request{},
					ms_panel.Router{
						List: config.Router,
					},*/
					//ms_panel.DB{
					// DataBase: config.DataBase.DataBase,
					//},
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
	// initDb(&config)
	injectDebug(&config)

	/*	for i := 0; i < len(config.Router); i++ {
		fmt.Printf("%v\n", config.Router[i].Path)
	}*/

	// Entry point
	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		// Prepare args
		virtualBuffer := bytes.NewBuffer([]byte{})
		virtualStatus := 200
		virtualResponse := ms_handler.VirtualResponseWriter{
			Response:   response,
			Buffer:     virtualBuffer,
			StatusCode: &virtualStatus,
		}
		args := ms_handler.Args{Response: virtualResponse, Request: request}
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

		// buf := bytes.NewBuffer([]byte{})
		// io.MultiWriter(args.Response, buf)

		// Handle
		h.Handle(&args)
	})

	// Start logger
	/*if config.Debug.UseLogs {
		ms_log.InitLogs(config.LogFile)
	}
	if config.Debug.UseRequestLogs {
		ms_log.InitRequestLogs(config.LogFile)
	}*/

	fmt.Printf("Mega Server Starts at host %v\n", config.Host)
	// ms_log.Log("info", fmt.Sprintf("Mega Server Starts at host %v", config.Host))

	if config.TLS.Enabled {
		err := http.ListenAndServeTLS(config.Host, config.TLS.CertFile, config.TLS.KeyFile, nil)
		ms_error.FatalIfError(err)
	} else {
		err := http.ListenAndServe(config.Host, nil)
		ms_error.FatalIfError(err)
	}
}
