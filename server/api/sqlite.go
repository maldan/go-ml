package ms_api

import (
	"database/sql"
	ms_error "github.com/maldan/go-ml/server/error"
	ml_sql "github.com/maldan/go-ml/util/sql"
)

type Sqlite struct {
	DataBaseMap map[string]*sql.DB
}

type ArgsSqliteDbName struct {
	DbName    string `json:"dbName"`
	TableName string `json:"tableName"`
}

type ArgsSqliteRowList struct {
	DbName    string `json:"dbName"`
	TableName string `json:"tableName"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
}

type ArgsSqliteUpdateRow struct {
	DbName    string `json:"dbName"`
	TableName string `json:"tableName"`
	RowId     int    `json:"rowId"`
}

func (s Sqlite) GetTableList(args ArgsSqliteDbName) any {
	list, err := ml_sql.Editor_GetTableList(s.DataBaseMap[args.DbName])
	ms_error.FatalIfError(err)
	return list
}

func (s Sqlite) GetTableFieldList(args ArgsSqliteDbName) any {
	list, err := ml_sql.Editor_GetTableFieldList(s.DataBaseMap[args.DbName], args.TableName)
	ms_error.FatalIfError(err)
	return list
}

func (s Sqlite) GetRowList(args ArgsSqliteRowList) any {
	r := ml_sql.Editor_SelectDynamic(ml_sql.SelectQuery{
		DB:     s.DataBaseMap[args.DbName],
		Table:  args.TableName,
		Offset: args.Offset,
		Limit:  args.Limit,
	})

	return map[string]any{
		"count":  r.Count,
		"result": r.Result,
	}
}

func (s Sqlite) UpdateRow(args ArgsSqliteUpdateRow) {

}
