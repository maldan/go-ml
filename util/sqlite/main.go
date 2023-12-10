package ml_sqlite

import "database/sql"

type QueryBuilder struct {
	DB    *sql.DB
	Query string
}

func (q *QueryBuilder) Select(fields []string) {
	q.Query += "SELECT "
}

func (q *QueryBuilder) Where() {

}

func (q *QueryBuilder) OrderBy() {

}

func (q *QueryBuilder) Execute() (any, error) {

	return nil, nil
}
