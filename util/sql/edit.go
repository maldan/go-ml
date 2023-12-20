package ml_sql

import (
	"database/sql"
	"fmt"
	"math"
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
		Query: "PRAGMA table_info(" + name + ")",
		Args:  []any{},
	})
	if r.Error != nil {
		return nil, r.Error
	}
	return r.Result, r.Error
}

func Editor_SelectDynamic(args SelectQuery) SelectResponseDynamic {
	fields := make([]string, 0)
	fieldTypes := make([]string, 0)
	fieldList, err := Editor_GetTableFieldList(args.DB, args.Table)
	if err != nil {
		return SelectResponseDynamic{Error: err}
	}

	for i := 0; i < len(fieldList); i++ {
		fields = append(fields, fieldList[i].Name)
		fieldTypes = append(fieldTypes, strings.ToLower(fieldList[i].Type))
	}
	fields = append(fields, "ROWID")
	fieldTypes = append(fieldTypes, "integer")

	query := fmt.Sprintf(
		"SELECT %v FROM %v",
		strings.Join(fields, ","),
		args.Table,
	)

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

	// Scan rows
	for rows.Next() {
		err2 := rows.Scan(destForScan...)
		if err2 != nil {
			fmt.Printf("%v\n", err2)
		}

		out := map[string]any{}

		//out := *new(T)
		//outType := reflect.TypeOf(&out).Elem()

		// Copy result
		for i := 0; i < len(fields); i++ {
			if fieldTypes[i] == "text" || strings.Contains(fieldTypes[i], "varchar") {
				out[fields[i]] = rawResult[i]
			} else if fieldTypes[i] == "integer" {
				str := string(rawResult[i])
				n, _ := strconv.Atoi(str)
				out[fields[i]] = n
			} else if fieldTypes[i] == "real" {
				str := string(rawResult[i])
				n, _ := strconv.ParseFloat(str, 64)
				if math.IsNaN(n) || math.IsInf(n, 1) || math.IsInf(n, -1) {
					n = 0
				}
				out[fields[i]] = n
			} else if fieldTypes[i] == "datetime" {
				t, err2 := time.Parse("2006-01-02T15:04:05.999999-07:00", string(rawResult[i]))
				if err2 != nil {
					// fmt.Printf("%v - %v\n", err2, string(rawResult[i]))
					t2, err3 := time.Parse("2006-01-02T15:04:05Z", string(rawResult[i]))
					if err3 != nil {
						// fmt.Printf("%v - %v\n", err3, string(rawResult[i]))
					} else {
						t = t2
					}
				}

				out[fields[i]] = t
				//ptr := unsafe.Add(unsafe.Pointer(&out), outType.Field(i).Offset)
				//*(*time.Time)(ptr) = t
			} else {
				out[fields[i]] = "???"
			}
			/*if outType.Field(i).Type.Kind() == reflect.Int8 ||
				outType.Field(i).Type.Kind() == reflect.Int16 ||
				outType.Field(i).Type.Kind() == reflect.Int32 ||
				outType.Field(i).Type.Kind() == reflect.Int ||
				outType.Field(i).Type.Kind() == reflect.Int64 {
				str := string(rawResult[i])
				n, _ := strconv.Atoi(str)
				reflect.ValueOf(&out).Elem().Field(i).SetInt(int64(n))
			}

			if outType.Field(i).Type.Kind() == reflect.Uint8 ||
				outType.Field(i).Type.Kind() == reflect.Uint16 ||
				outType.Field(i).Type.Kind() == reflect.Uint32 ||
				outType.Field(i).Type.Kind() == reflect.Uint ||
				outType.Field(i).Type.Kind() == reflect.Uint64 {
				str := string(rawResult[i])
				n, _ := strconv.Atoi(str)
				reflect.ValueOf(&out).Elem().Field(i).SetUint(uint64(n))
			}

			if outType.Field(i).Type.Kind() == reflect.Float32 ||
				outType.Field(i).Type.Kind() == reflect.Float64 {
				str := string(rawResult[i])
				n, _ := strconv.ParseFloat(str, 64)
				if math.IsNaN(n) || math.IsInf(n, 1) || math.IsInf(n, -1) {
					n = 0
				}
				reflect.ValueOf(&out).Elem().Field(i).SetFloat(n)
			}

			if outType.Field(i).Type.Kind() == reflect.Bool {
				if len(rawResult[i]) > 0 {
					reflect.ValueOf(&out).Elem().Field(i).SetBool(rawResult[i][0] == 49)
				}
			}
			if outType.Field(i).Type.Kind() == reflect.String {
				if len(rawResult[i]) > 0 {
					reflect.ValueOf(&out).Elem().Field(i).SetString(string(rawResult[i]))
				}
			}
			if outType.Field(i).Type.Name() == "Time" {
				t, err2 := time.Parse("2006-01-02T15:04:05.999999-07:00", string(rawResult[i]))
				if err2 != nil {
					// fmt.Printf("%v - %v\n", err2, string(rawResult[i]))
					t2, err3 := time.Parse("2006-01-02T15:04:05Z", string(rawResult[i]))
					if err3 != nil {
						// fmt.Printf("%v - %v\n", err3, string(rawResult[i]))
					} else {
						t = t2
					}
				}

				ptr := unsafe.Add(unsafe.Pointer(&out), outType.Field(i).Offset)
				*(*time.Time)(ptr) = t
			}*/
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
	for k, _ := range args.Data {
		set = append(set, k)
	}

	query := fmt.Sprintf("UPDATE %v SET %v WHERE ROWID=%v", args.Table, strings.Join(set, ","), args.RowId)
	fmt.Printf("%v\n", query)

	// Prepare
	statement, err := args.DB.Prepare(query)
	defer statement.Close()
	if err != nil {
		return err
	}

	// Change time
	/*for i := 0; i < len(args.SetArgs); i++ {
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
	}*/

	return nil
}
