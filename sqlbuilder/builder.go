package sqlbuilder

import (
	"fmt"
	"strings"
)

// NewMySQLQuery new 实例
func NewMySQLQuery(querySQL string) *Query {
	return &Query{
		query:     querySQL,
		whereOr:   []string{},
		whereAnd:  []string{},
		argsOr:    []interface{}{},
		argsAnd:   []interface{}{},
		argsLimit: []interface{}{},
	}
}

// Query 查询mysql数据库
type Query struct {
	query     string
	whereOr   []string
	whereAnd  []string
	argsOr    []interface{}
	argsAnd   []interface{}
	limit     string
	argsLimit []interface{}
	desc      string
}

// AddOrArg 添加参数
func (q *Query) AddOrArg(v ...interface{}) {
	q.argsOr = append(q.argsOr, v...)
}

// AddOrWhere 添加Where语句
func (q *Query) AddOrWhere(v ...string) {
	q.whereOr = append(q.whereOr, v...)
}

// AddAndArg 添加参数
func (q *Query) AddAndArg(v ...interface{}) {
	q.argsAnd = append(q.argsAnd, v...)
}

// AddAndWhere 添加Where语句
func (q *Query) AddAndWhere(v ...string) {
	q.whereAnd = append(q.whereAnd, v...)
}

// Limit 这周Limit
func (q *Query) Limit(offset int64, limit uint) {
	q.limit = "LIMIT ?,? "
	q.argsLimit = append(q.argsLimit, offset, limit)
}

// Desc todo
func (q *Query) Desc(d string) {
	q.desc = fmt.Sprintf("ORDER BY %s DESC ", d)
}

func (q *Query) whereBuild() string {
	if len(q.whereAnd) == 0 && len(q.whereOr) == 0 {
		return ""
	}

	where := []string{}
	if len(q.whereAnd) > 0 {
		where = append(where, strings.Join(q.whereAnd, " AND "))
	}
	if len(q.whereOr) > 0 {
		where = append(where, " ( "+strings.Join(q.whereOr, " OR ")+" ) ")
	}

	return "WHERE " + strings.Join(where, " AND ")
}

// Build 组件SQL
func (q *Query) Build() string {
	return q.query + q.whereBuild() + q.desc + q.limit + ";"
}

// Args sql参数
func (q *Query) Args() []interface{} {
	args := []interface{}{}
	args = append(args, q.argsAnd...)
	args = append(args, q.argsOr...)
	args = append(args, q.argsLimit...)
	return args
}
