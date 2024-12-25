package ms_api

import (
	"database/sql"
	ms_error "github.com/maldan/go-ml/server/error"
	ml_slice "github.com/maldan/go-ml/util/slice"
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
	DbName    string         `json:"dbName"`
	TableName string         `json:"tableName"`
	RowId     int            `json:"rowId"`
	Data      map[string]any `json:"data"`
}

type ArgsSqliteInsertRow struct {
	DbName    string         `json:"dbName"`
	TableName string         `json:"tableName"`
	Data      map[string]any `json:"data"`
}

type ArgsSqliteDeleteRow struct {
	DbName    string `json:"dbName"`
	TableName string `json:"tableName"`
	RowId     int    `json:"rowId"`
}

func (s Sqlite) GetDbList(args ArgsSqliteDbName) []string {
	out := make([]string, 0)
	for k, _ := range s.DataBaseMap {
		out = append(out, k)
	}
	ml_slice.SortAZ(out)
	return out
}

func (s Sqlite) GetTableList(args ArgsSqliteDbName) any {
	list, err := ml_sql.Editor_GetTableList(s.DataBaseMap[args.DbName])
	ms_error.FatalIfError(err)

	out := make([]string, 0)
	for _, v := range list {
		out = append(out, v.Name)
	}
	ml_slice.SortAZ(out)
	return out
}

func (s Sqlite) GetTableFieldList(args ArgsSqliteDbName) any {
	list, err := ml_sql.Editor_GetTableFieldList(s.DataBaseMap[args.DbName], args.TableName)
	ms_error.FatalIfError(err)
	return list
}

func (s Sqlite) PostRawQuery(args ArgsSqliteRowList) map[string]any {

	r := ml_sql.Editor_SelectDynamic(ml_sql.Editor_SelectQuery{
		DB:     s.DataBaseMap[args.DbName],
		Table:  args.TableName,
		Offset: args.Offset,
		Limit:  args.Limit,
	})
	total, _ := ml_sql.Count(ml_sql.CountQuery{
		DB:    s.DataBaseMap[args.DbName],
		Table: args.TableName,
	})

	return map[string]any{
		"count":  r.Count,
		"total":  total,
		"result": r.Result,
	}
}

func (s Sqlite) PostRowList(args ArgsSqliteRowList) map[string]any {
	r := ml_sql.Editor_SelectDynamic(ml_sql.Editor_SelectQuery{
		DB:     s.DataBaseMap[args.DbName],
		Table:  args.TableName,
		Offset: args.Offset,
		Limit:  args.Limit,
	})
	total, _ := ml_sql.Count(ml_sql.CountQuery{
		DB:    s.DataBaseMap[args.DbName],
		Table: args.TableName,
	})

	return map[string]any{
		"count":  r.Count,
		"total":  total,
		"result": r.Result,
	}
}

func (s Sqlite) PostRow(args ArgsSqliteInsertRow) {
	_, err := ml_sql.Editor_Insert(s.DataBaseMap[args.DbName], args.TableName, args.Data)
	ms_error.FatalIfError(err)
}

func (s Sqlite) PatchRow(args ArgsSqliteUpdateRow) {
	err := ml_sql.Editor_UpdateRow(ml_sql.Editor_UpdateQuery{
		DB:    s.DataBaseMap[args.DbName],
		Table: args.TableName,
		RowId: args.RowId,
		Data:  args.Data,
	})
	ms_error.FatalIfError(err)
}

func (s Sqlite) DeleteRow(args ArgsSqliteDeleteRow) {
	err := ml_sql.Editor_DeleteRow(ml_sql.Editor_DeleteQueryRow{
		DB:    s.DataBaseMap[args.DbName],
		Table: args.TableName,
		RowId: args.RowId,
	})
	ms_error.FatalIfError(err)
}
