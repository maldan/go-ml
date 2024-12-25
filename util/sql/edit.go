package ml_sql

import (
	"database/sql"
	"fmt"
	ml_slice "github.com/maldan/go-ml/util/slice"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type TableList struct {
	Name string `json:"name"`
}

type TableField struct {
	Id         int    `json:"cid"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	NotNull    int    `json:"notnull"`
	Default    any    `json:"dflt_value"`
	PrimaryKey int    `json:"pk"`
}

type SelectResponseDynamic struct {
	Result  []map[string]any
	IsFound bool
	Count   int
	Error   error
}

type Editor_UpdateQuery struct {
	DB    *sql.DB
	Table string
	Data  map[string]any
	RowId int
}

type Editor_DeleteQueryRow struct {
	DB    *sql.DB
	Table string
	RowId int
}

type Editor_SelectQuery struct {
	DB        *sql.DB
	Table     string
	Where     string
	WhereArgs []any
	OrderBy   string
	Offset    int
	Limit     int
	JoinTable []string
	JoinOn    []string
}

func Editor_GetTableList(DB *sql.DB) ([]TableList, error) {
	r := Select[TableList](SelectQuery{
		DB:        DB,
		Table:     "sqlite_master",
		Where:     "type = ?",
		WhereArgs: []any{"table"},
	})

	return r.Result, r.Error
}

func Editor_GetTableFieldList(DB *sql.DB, name string) ([]TableField, error) {
	r := Raw[TableField](RawQuery{
		DB:    DB,
		Query: "PRAGMA table_info(`" + name + "`)",
		Args:  []any{},
	})
	if r.Error != nil {
		return nil, r.Error
	}
	return r.Result, r.Error
}

func Editor_RawQuery(args RawQuery) SelectResponseDynamic {
	response := SelectResponseDynamic{}

	// Prepare
	statement, err := args.DB.Prepare(args.Query)
	defer statement.Close()
	if err != nil {
		response.Error = err
		return response
	}

	rows, err := statement.Query()
	defer rows.Close()
	if err != nil {
		response.Error = err
		return response
	}

	// Scan rows
	/*for rows.Next() {
		err2 := rows.Scan(destForScan...)
		if err2 != nil {
			fmt.Printf("%v\n", err2)
		}

		out := map[string]any{}

		response.IsFound = true
		response.Result = append(response.Result, out)
	}*/
	response.Count = len(response.Result)

	return response
}

func Editor_SelectDynamic(args Editor_SelectQuery) SelectResponseDynamic {
	fields := make([]string, 0)
	fieldTypes := make([]string, 0)
	fieldList, err := Editor_GetTableFieldList(args.DB, args.Table)
	if err != nil {
		return SelectResponseDynamic{Error: err}
	}

	// Join mode
	isJoinMode := false
	joinString := ""
	if len(args.JoinTable) > 0 {
		isJoinMode = true
	}
	if isJoinMode {
		// Add tables
		for i := 0; i < len(args.JoinTable); i++ {
			joinString += fmt.Sprintf(" JOIN `%v` ON %v", args.JoinTable[i], args.JoinOn[i])
		}

		// Fill fields
		for i := 0; i < len(fieldList); i++ {
			fields = append(fields, fmt.Sprintf("`%v`.`%v`", args.Table, fieldList[i].Name))
			fieldTypes = append(fieldTypes, strings.ToLower(fieldList[i].Type))
		}

		// File join tables field
		for i := 0; i < len(args.JoinTable); i++ {
			joinFieldList, err2 := Editor_GetTableFieldList(args.DB, args.JoinTable[i])
			if err2 != nil {
				return SelectResponseDynamic{Error: err2}
			}

			// Fill fields
			for j := 0; j < len(joinFieldList); j++ {
				fields = append(fields, fmt.Sprintf("`%v`.`%v`", args.JoinTable[i], joinFieldList[j].Name))
				fieldTypes = append(fieldTypes, strings.ToLower(joinFieldList[j].Type))
			}
		}
	} else {
		// Fill fields
		for i := 0; i < len(fieldList); i++ {
			fields = append(fields, fmt.Sprintf("`%v`", fieldList[i].Name))
			fieldTypes = append(fieldTypes, strings.ToLower(fieldList[i].Type))
		}

		fields = append(fields, "ROWID")
		fieldTypes = append(fieldTypes, "integer")
	}

	query := fmt.Sprintf(
		"SELECT %v FROM `%v`",
		strings.Join(fields, ","),
		args.Table,
	)

	// Join
	if isJoinMode {
		query += joinString
	}

	// Where
	if args.Where != "" {
		query += fmt.Sprintf(" WHERE %v", args.Where)
	}

	// Order
	if args.OrderBy != "" {
		query += fmt.Sprintf(" ORDER BY %v", args.OrderBy)
	}

	// Set limit
	if args.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %v", args.Limit)
	}

	// Offset
	if args.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %v", args.Offset)
	}

	fmt.Printf("SELECT: %v\n", query)

	response := SelectResponseDynamic{}

	destForScan := make([]any, len(fields))
	rawResult := make([]sql.RawBytes, len(fields))
	for i, _ := range destForScan {
		destForScan[i] = &rawResult[i]
	}

	// Prepare
	statement, err := args.DB.Prepare(query)
	defer statement.Close()
	if err != nil {
		response.Error = err
		return response
	}

	// Execute statement
	if args.WhereArgs == nil {
		args.WhereArgs = make([]any, 0)
	}

	rows, err := statement.Query(args.WhereArgs...)
	defer rows.Close()
	if err != nil {
		response.Error = err
		return response
	}

	removeQuotes := func(str string) string {
		return strings.ReplaceAll(str, "`", "")
	}

	// Scan rows
	for rows.Next() {
		err2 := rows.Scan(destForScan...)
		if err2 != nil {
			fmt.Printf("%v\n", err2)
		}

		out := map[string]any{}

		// Copy result
		for i := 0; i < len(fields); i++ {
			if fieldTypes[i] == "text" || strings.Contains(fieldTypes[i], "varchar") {
				out[removeQuotes(fields[i])] = string(rawResult[i])
			} else if fieldTypes[i] == "integer" {
				str := string(rawResult[i])
				n, _ := strconv.Atoi(str)
				out[removeQuotes(fields[i])] = n
			} else if fieldTypes[i] == "real" {
				str := string(rawResult[i])
				n, _ := strconv.ParseFloat(str, 64)
				if math.IsNaN(n) || math.IsInf(n, 1) || math.IsInf(n, -1) {
					n = 0
				}
				out[removeQuotes(fields[i])] = n
			} else if fieldTypes[i] == "datetime" {
				layouts := []string{
					"2006-01-02T15:04:05.999999-07:00",
					"2006-01-02 15:04:05.999999-07:00",
					"2006-01-02 15:04:05.999999 -07:00",
					"2006-01-02 15:04:05 -07:00",
					"2006-01-02T15:04:05Z",
					"2006-01-02 15:04:05Z",
					"2006-01-02 15:04:05",
				}
				out[removeQuotes(fields[i])] = "???"
				for j := 0; j < len(layouts); j++ {
					t, err3 := time.Parse(layouts[j], string(rawResult[i]))
					if err3 == nil {
						out[removeQuotes(fields[i])] = t.Format("2006-01-02 15:04:05.999999 -07:00")
						break
					}
				}

			} else if fieldTypes[i] == "boolean" {
				if strings.ToLower(string(rawResult[i])) == "true" {
					out[removeQuotes(fields[i])] = true
				} else if strings.ToLower(string(rawResult[i])) == "false" {
					out[removeQuotes(fields[i])] = false
				} else {
					out[removeQuotes(fields[i])] = string(rawResult[i]) != "0"
				}
			} else {
				out[removeQuotes(fields[i])] = "???"
			}
		}

		response.IsFound = true
		response.Result = append(response.Result, out)
		// outList = append(outList, out)
	}
	response.Count = len(response.Result)

	return response
}

func Editor_UpdateRow(args Editor_UpdateQuery) error {
	set := make([]string, 0)
	values := make([]any, 0)
	for k, v := range args.Data {
		set = append(set, k+"=?")
		values = append(values, v)
	}

	query := fmt.Sprintf("UPDATE %v SET %v WHERE ROWID=%v", args.Table, strings.Join(set, ", "), args.RowId)
	fmt.Printf("%v\n", query)

	// Prepare
	statement, err := args.DB.Prepare(query)
	defer statement.Close()
	if err != nil {
		return err
	}

	// Execute
	_, err = statement.Exec(values...)
	if err != nil {
		return err
	}

	return nil
}

func Editor_DeleteRow(args Editor_DeleteQueryRow) error {
	query := fmt.Sprintf("DELETE FROM %v WHERE ROWID = %v", args.Table, args.RowId)
	fmt.Printf("SQL EDITOR: %v\n", query)

	// Prepare
	statement, err := args.DB.Prepare(query)
	defer statement.Close()
	if err != nil {
		return err
	}

	// Execute
	_, err = statement.Exec(args.RowId)
	if err != nil {
		return err
	}

	return nil
}

func Editor_Delete(args DeleteQuery) (int64, error) {
	query := fmt.Sprintf("DELETE FROM %v WHERE %v", args.Table, args.Where)

	// Prepare
	statement, err := args.DB.Prepare(query)
	defer statement.Close()
	if err != nil {
		return 0, err
	}

	// Execute statement
	r, err := statement.Exec(args.WhereArgs...)
	affected, _ := r.RowsAffected()

	return affected, err
}

func Editor_Update(args UpdateQuery) error {
	query := fmt.Sprintf("UPDATE `%v` SET %v WHERE %v", args.Table, args.Set, args.Where)
	fmt.Printf("UPDATE: %v\n", query)

	// Prepare
	statement, err := args.DB.Prepare(query)
	defer statement.Close()
	if err != nil {
		return err
	}

	// Change time
	for i := 0; i < len(args.SetArgs); i++ {
		switch args.SetArgs[i].(type) {
		case time.Time:
			fieldValue := args.SetArgs[i].(time.Time).UTC()
			args.SetArgs[i] = fieldValue.Format("2006-01-02 15:04:05.000-07:00")
			break
		default:
			break
		}
	}

	// Execute statement
	all := ml_slice.Combine(args.SetArgs, args.WhereArgs)
	_, err = statement.Exec(all...)
	if err != nil {
		return err
	}

	return nil
}

func Editor_Insert(db *sql.DB, table string, data map[string]any) (int64, error) {
	fields := make([]string, 0)
	values := make([]any, 0)
	valuesQ := make([]string, len(data))

	for k, v := range data {
		fields = append(fields, fmt.Sprintf("`%v`", k))
		values = append(values, v)
	}

	query := fmt.Sprintf("INSERT INTO '%v'(\n", table)

	for i := 0; i < len(fields); i++ {
		if fields[i] == "`id`" && (values[i] == nil || reflect.ValueOf(values[i]).IsZero()) {
			values[i] = nil
		}

		fields[i] = "\t" + fields[i]
		valuesQ[i] = "\t?"
	}

	query += "" + strings.Join(fields, ",\n ") + "\n) \n"
	query += "VALUES(\n" + strings.Join(valuesQ, ",\n") + "\n)"

	// Prepare
	statement, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}

	fmt.Printf("%v\n", query)

	// Execute statement
	r, err := statement.Exec(values...)
	if err != nil {
		return 0, err
	}

	lastId, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}
