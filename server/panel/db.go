package ms_panel

import (
	"github.com/maldan/go-ml/db/mdb"
	ms_error "github.com/maldan/go-ml/server/error"
	ml_convert "github.com/maldan/go-ml/util/convert"
	ml_slice "github.com/maldan/go-ml/util/slice"
)

type DB struct {
	DataBase *map[string]*mdb.DataTable
}

type ArgsDBSearch struct {
	Page  int    `json:"page"`
	Table string `json:"table"`
	Where string `json:"where"`
}

func (d DB) GetList() []string {
	return ml_slice.GetKeys(*d.DataBase)
}

func (d DB) GetSearch(args ArgsDBSearch) mdb.SearchResult {
	table, ok := (*d.DataBase)[args.Table]
	ms_error.FatalIf(!ok, ms_error.Error{Code: 404})

	whereB, _ := ml_convert.FromBase64(args.Where)
	where := string(whereB)
	if where == "" {
		where = "1 == 1"
	}

	res := table.FindBy(mdb.ArgsFind{
		Offset:          args.Page * 20,
		Limit:           20,
		WhereExpression: where,
	})

	res.Total = table.Count(mdb.ArgsFind{WhereExpression: where})
	res.PerPage = 20
	res.Count = len(res.Result)

	return res
}
