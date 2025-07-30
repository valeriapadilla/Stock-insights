package repository

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	selectClause string
	fromClause   string
	whereClause  []string
	orderClause  string
	limitClause  string
	offsetClause string
	args         []interface{}
	argIndex     int
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		args:     make([]interface{}, 0),
		argIndex: 1,
	}
}

func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.selectClause = strings.Join(columns, ", ")
	return qb
}

func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.fromClause = table
	return qb
}

func (qb *QueryBuilder) Where(condition string, value interface{}) *QueryBuilder {
	if value != nil && value != "" {
		placeholder := fmt.Sprintf("$%d", qb.argIndex)
		condition = strings.Replace(condition, "?", placeholder, 1)
		qb.whereClause = append(qb.whereClause, condition)
		qb.args = append(qb.args, value)
		qb.argIndex++
	}
	return qb
}

func (qb *QueryBuilder) OrderBy(column string, direction string) *QueryBuilder {
	qb.orderClause = fmt.Sprintf("ORDER BY %s %s", column, direction)
	return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limitClause = fmt.Sprintf("LIMIT %d", limit)
	return qb
}

func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offsetClause = fmt.Sprintf("OFFSET %d", offset)
	return qb
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
	query := fmt.Sprintf("SELECT %s FROM %s", qb.selectClause, qb.fromClause)

	if len(qb.whereClause) > 0 {
		query += " WHERE " + strings.Join(qb.whereClause, " AND ")
	}

	if qb.orderClause != "" {
		query += " " + qb.orderClause
	}

	if qb.limitClause != "" {
		query += " " + qb.limitClause
	}

	if qb.offsetClause != "" {
		query += " " + qb.offsetClause
	}

	return query, qb.args
}

func (qb *QueryBuilder) CountQuery() (string, []interface{}) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", qb.fromClause)

	if len(qb.whereClause) > 0 {
		query += " WHERE " + strings.Join(qb.whereClause, " AND ")
	}

	return query, qb.args
}
