package ms_panel

import ms_response "github.com/maldan/go-ml/server/response"

type Panel struct{}

func (p Panel) GetIndex() ms_response.Custom {
	return ms_response.Custom{
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Body: Html,
	}
}

func (p Panel) GetCss() ms_response.Custom {
	return ms_response.Custom{
		Headers: map[string]string{
			"Content-Type": "text/css",
		},
		Body: Css,
	}
}

func (p Panel) GetJs() ms_response.Custom {
	return ms_response.Custom{
		Headers: map[string]string{
			"Content-Type": "text/javascript",
		},
		Body: Js,
	}
}

func (p Panel) GetSetting() {

}
