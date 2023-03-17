package ms_panel

type Panel struct {
	HasLogTab      bool
	HasRequestLogs bool
}

func (p Panel) GetSetting() any {
	return map[string]any{
		"hasLogTab":     p.HasLogTab,
		"hasRequestTab": p.HasRequestLogs,
	}
}
