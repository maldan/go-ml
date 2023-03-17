package ms

import (
	"github.com/maldan/go-ml/db/mdb"
	ms_handler "github.com/maldan/go-ml/server/core/handler"
)

type PanelConfig struct {
	Login    string
	Password string

	//HasRequestTab bool
	HasMethodTab  bool
	HasTestTab    bool
	HasDbTab      bool
	HasControlTab bool
	//HasLogTab     bool
}

type SecureConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

type TableConfig struct {
	Name string
	Type any
}

type DebugConfig struct {
	UseLogs        bool
	UseRequestLogs bool
}

type DataBaseConfig struct {
	Path      string
	DataBase  *map[string]*mdb.DataTable
	TableList []TableConfig
}

type Config struct {
	Host     string
	Router   []ms_handler.RouteHandler
	TLS      SecureConfig
	Debug    DebugConfig
	Panel    PanelConfig
	LogFile  string
	DataBase DataBaseConfig
	// TableList []TableConfig
}
