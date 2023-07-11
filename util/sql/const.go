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
	Offset    int
	Limit     int
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
