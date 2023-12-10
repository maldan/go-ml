package ml_sql

import "database/sql"

type DeleteQuery struct {
	DB        *sql.DB
	Table     string
	Where     string
	WhereArgs []any
}

type SelectQuery struct {
	DB        *sql.DB
	Table     string
	Where     string
	WhereArgs []any
	OrderBy   string
	// OrderDirection bool
	Offset  int
	Limit   int
	FieldAs map[string]string
	// EncryptedFields []string
}

type RawQuery struct {
	DB    *sql.DB
	Query string
	Args  []any
}

type SelectResponse[T any] struct {
	Result  []T
	IsFound bool
	Count   int
	Error   error
}

type UpdateQuery struct {
	DB        *sql.DB
	Table     string
	Where     string
	WhereArgs []any
	Set       string
	SetArgs   []any
}

type UpdateSimpleQuery struct {
	DB        *sql.DB
	Table     string
	Where     string
	WhereArgs []any
	Set       string
	SetArgs   []any
}

type CountQuery struct {
	DB        *sql.DB
	Table     string
	Where     string
	WhereArgs []any
}
