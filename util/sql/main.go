package ml_sql

import (
	"database/sql"
	"fmt"
	ml_slice "github.com/maldan/go-ml/util/slice"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func typeToSqlType(t string, size string) string {
	switch t {
	case "bool":
		return "BOOLEAN"
	case "int8":
		return "TINYINT"
	case "uint8":
		return "TINYINT UNSIGNED"
	case "int16":
		return "SMALLINT"
	case "uint16":
		return "SMALLINT UNSIGNED"
	case "int32", "int":
		return "INTEGER"
	case "uint32", "uint":
		return "INTEGER UNSIGNED"
	case "int64":
		return "BIGINT"
	case "uint64":
		return "BIGINT UNSIGNED"
	case "float32", "float64":
		return "REAL"
	case "string":
		if size == "" {
			return "TEXT"
		} else {
			return "VARCHAR(" + size + ")"
		}
	case "Time":
		return "DATETIME"
	default:
		panic("unknown type " + t)
	}
	return ""
}

func getValueFieldNames[T any](v T, useQuotes bool) []string {
	typeOf := reflect.TypeOf(v)
	out := make([]string, 0)

	for i := 0; i < typeOf.NumField(); i++ {
		fieldName := typeOf.Field(i).Name
		if typeOf.Field(i).Tag.Get("json") != "" {
			fieldName = typeOf.Field(i).Tag.Get("json")
			if fieldName == "-" {
				continue
			}
		}
		if useQuotes {
			out = append(out, "`"+fieldName+"`")
		} else {
			out = append(out, fieldName)
		}
	}

	return out
}

func getValues[T any](v T) []any {
	typeOf := reflect.TypeOf(v)
	valueOf := reflect.ValueOf(v)
	out := make([]any, 0)

	for i := 0; i < valueOf.NumField(); i++ {
		// Check name
		if typeOf.Field(i).Tag.Get("json") != "" {
			fieldName := typeOf.Field(i).Tag.Get("json")
			if fieldName == "-" {
				continue
			}
		}

		// Time
		if typeOf.Field(i).Type.Name() == "Time" {
			fieldValue := valueOf.Field(i).Interface().(time.Time).UTC()
			out = append(out, fieldValue.Format("2006-01-02 15:04:05.000-07:00"))
		} else {
			fieldValue := valueOf.Field(i).Interface()
			out = append(out, fieldValue)
		}
	}

	return out
}

func CreateTable[T any](db *sql.DB, name string) error {
	out := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (\n", name)

	typeOf := reflect.TypeOf(*new(T))

	for i := 0; i < typeOf.NumField(); i++ {
		fieldName := typeOf.Field(i).Name
		fieldType := typeToSqlType(typeOf.Field(i).Type.Name(), typeOf.Field(i).Tag.Get("len"))
		opts := ""

		if typeOf.Field(i).Name == "Id" {
			if typeOf.Field(i).Type.Kind() == reflect.String {
				fieldType = "TEXT"
				opts += "PRIMARY KEY "
			} else {
				fieldType = "INTEGER"
				opts += "PRIMARY KEY AUTOINCREMENT "
			}
		} else {
			if typeOf.Field(i).Type.Name() == "string" {
				opts += "DEFAULT \"\" "
			} else {
				opts += "DEFAULT 0 "
			}
		}
		opts += "NOT NULL "

		if typeOf.Field(i).Tag.Get("json") != "" {
			fieldName = typeOf.Field(i).Tag.Get("json")
			if fieldName == "-" {
				continue
			}
		}

		out += fmt.Sprintf("\t\t`%v` %v %v", fieldName, fieldType, opts)

		if i != typeOf.NumField()-1 {
			out += ","
		}
		out += "\n"
	}

	out += ");\n"
	fmt.Printf("%v", out)
	// Prepare
	statement, err := db.Prepare(out)
	if err != nil {
		return err
	}

	// Execute
	_, err = statement.Exec()

	return err
}

func Insert[T any](db *sql.DB, table string, value T) (int64, error) {
	fields := getValueFieldNames(value, true)
	values := getValues(value)
	valuesQ := make([]string, len(values))

	query := fmt.Sprintf("INSERT INTO '%v'(\n", table)

	for i := 0; i < len(fields); i++ {
		if fields[i] == "`id`" && reflect.ValueOf(values[i]).IsZero() {
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

	// fmt.Printf("%v\n", query)
	// fmt.Printf("%v\n", values[0])

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

/*func SelectOne[T any](db *sql.DB, from string, where string, values ...any) (*T, error) {
	out := *new(T)
	outType := reflect.TypeOf(&out).Elem()

	fields := getValueFieldNames(out, false)
	query := fmt.Sprintf("SELECT %v FROM %v WHERE %v LIMIT 1", strings.Join(fields, ","), from, where)

	destForScan := make([]any, len(fields))
	rawResult := make([]sql.RawBytes, len(fields))
	for i, _ := range destForScan {
		destForScan[i] = &rawResult[i]
	}

	// Prepare
	statement, err := db.Prepare(query)
	defer statement.Close()
	if err != nil {
		return nil, err
	}

	// Execute statement
	rows, err := statement.Query(values...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	// Scan rows
	found := false
	for rows.Next() {
		err2 := rows.Scan(destForScan...)
		if err2 != nil {
			return nil, err2
		}
		found = true
	}

	if !found {
		return nil, nil
	}

	// Copy result
	for i := 0; i < len(fields); i++ {
		if outType.Field(i).Type.Kind() == reflect.Int8 ||
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
				fmt.Printf("%v\n", err2)
				t2, err3 := time.Parse("2006-01-02T15:04:05Z", string(rawResult[i]))
				if err3 != nil {
					fmt.Printf("%v\n", err3)
				} else {
					t = t2
				}
			}

			ptr := unsafe.Add(unsafe.Pointer(&out), outType.Field(i).Offset)
			*(*time.Time)(ptr) = t
		}
	}

	return &out, err
}*/

/*func SelectMany[T any](db *sql.DB, from string, where string, values ...any) ([]T, error) {
	fields := getValueFieldNames(*new(T), false)
	query := fmt.Sprintf("SELECT %v FROM %v WHERE %v", strings.Join(fields, ","), from, where)

	destForScan := make([]any, len(fields))
	rawResult := make([]sql.RawBytes, len(fields))
	for i, _ := range destForScan {
		destForScan[i] = &rawResult[i]
	}

	// Prepare
	statement, err := db.Prepare(query)
	defer statement.Close()
	if err != nil {
		return make([]T, 0), err
	}

	// Execute statement
	rows, err := statement.Query(values...)
	defer rows.Close()
	if err != nil {
		return make([]T, 0), err
	}

	// Scan rows
	outList := make([]T, 0)
	for rows.Next() {
		err2 := rows.Scan(destForScan...)
		if err2 != nil {
			fmt.Printf("%v\n", err2)
		}

		out := *new(T)
		outType := reflect.TypeOf(&out).Elem()

		// Copy result
		for i := 0; i < len(fields); i++ {
			if outType.Field(i).Type.Kind() == reflect.Int8 ||
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
					fmt.Printf("%v\n", err2)
					t2, err3 := time.Parse("2006-01-02T15:04:05Z", string(rawResult[i]))
					if err3 != nil {
						fmt.Printf("%v\n", err3)
					} else {
						t = t2
					}
				}

				ptr := unsafe.Add(unsafe.Pointer(&out), outType.Field(i).Offset)
				*(*time.Time)(ptr) = t
			}
		}

		outList = append(outList, out)
	}

	return outList, err
}*/

func Count(db *sql.DB, from string, where string, values ...any) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %v WHERE %v LIMIT 1", from, where)
	count := 0

	// Prepare
	statement, err := db.Prepare(query)
	defer statement.Close()
	if err != nil {
		return 0, err
	}

	// Execute statement
	rows, err := statement.Query(values...)
	defer rows.Close()
	if err != nil {
		return 0, err
	}

	// Scan rows
	found := false
	for rows.Next() {
		err2 := rows.Scan(&count)
		if err2 != nil {
			return 0, err2
		}
		found = true
	}

	if !found {
		return 0, nil
	}

	return count, err
}

func Delete(args DeleteQuery) error {
	query := fmt.Sprintf("DELETE FROM %v WHERE %v", args.Table, args.Where)

	// Prepare
	statement, err := args.DB.Prepare(query)
	defer statement.Close()
	if err != nil {
		return err
	}

	// Execute statement
	_, err = statement.Exec(args.WhereArgs...)
	return err
}

func Update(args UpdateQuery) error {
	query := fmt.Sprintf("UPDATE %v SET %v WHERE %v", args.Table, args.Set, args.Where)
	fmt.Printf("%v\n", query)
	// Prepare
	statement, err := args.DB.Prepare(query)
	defer statement.Close()
	if err != nil {
		return err
	}

	// Execute statement
	all := ml_slice.Combine(args.SetArgs, args.WhereArgs)
	_, err = statement.Exec(all...)
	if err != nil {
		return err
	}

	return nil
}

func Select[T any](args SelectQuery) SelectResponse[T] {
	fields := getValueFieldNames(*new(T), false)
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

	response := SelectResponse[T]{}

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

		out := *new(T)
		outType := reflect.TypeOf(&out).Elem()

		// Copy result
		for i := 0; i < len(fields); i++ {
			if outType.Field(i).Type.Kind() == reflect.Int8 ||
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
			}
		}

		response.IsFound = true
		response.Result = append(response.Result, out)
		// outList = append(outList, out)
	}
	response.Count = len(response.Result)

	return response
}

func AlterTableAddColumn(db *sql.DB, table string, name string, kind string, defaultValue string) error {
	query := fmt.Sprintf("ALTER TABLE %v ADD COLUMN %v %v DEFAULT %v NOT NULL", table, name, kind, defaultValue)

	// Prepare
	statement, err := db.Prepare(query)
	defer statement.Close()
	if err != nil {
		return err
	}

	// Execute statement
	_, err = statement.Exec()
	return err
}
