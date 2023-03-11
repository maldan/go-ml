package ms_panel

type Panel struct {
	HasLogTab bool
}

func (p Panel) GetSetting() any {
	return map[string]any{
		"hasLogTab": p.HasLogTab,
	}
}
