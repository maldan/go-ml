package ms_error

const ErrorTypeUnknown = "unknown"
const ErrorTypeRequired = "required"

type ErrorDebugInfo struct {
	File string `json:"file"`
}

type Error struct {
	Code        int      `json:"-"`
	Type        string   `json:"type"`
	Field       string   `json:"field,omitempty"`
	Description string   `json:"description"`
	Debug       []string `json:"debug,omitempty"`
	EndPoint    string   `json:"endPoint,omitempty"`
}

func Fatal(err Error) {
	/*_, file, line, _ := runtime.Caller(1)
	ff := strings.Split(file, "/")*/

	if err.Code == 0 {
		err.Code = 500
	}
	if err.Type == "" {
		err.Type = "unknown"
	}
	/*err.File = strings.Join(ff[len(ff)-2:], "/")
	err.Line = line
	err.Created = time.Now()*/

	panic(err)
}

func FatalIfError(err error) {
	if err != nil {
		Fatal(Error{
			Description: err.Error(),
		})
	}
}

func FatalIf(ok bool, err Error) {
	if ok {
		Fatal(err)
	}
}
