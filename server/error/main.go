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
	Stack       []string `json:"stack,omitempty"`
	EndPoint    string   `json:"endPoint,omitempty"`
}

func Fatal(err Error) {
	if err.Code == 0 {
		err.Code = 500
	}
	if err.Type == "" {
		err.Type = "unknown"
	}

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
