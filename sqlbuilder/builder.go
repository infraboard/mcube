package sqlbuilder

import (
	"fmt"
	"strings"
)

// NewQuery new 实例
func NewQuery(querySQL string) *Query {
	return &Query{
		query:     querySQL,
		whereStmt: []string{},
		whereArgs: []interface{}{},
		limitArgs: []interface{}{},
	}
}

// Query 查询mysql数据库
type Query struct {
	query     string
	whereStmt []string
	whereArgs []interface{}
	limitStmt string
	limitArgs []interface{}
	order     string
}

// Where 添加参数
func (q *Query) Where(stmt string, v ...interface{}) *Query {
	q.whereStmt = append(q.whereStmt, stmt)
	q.whereArgs = append(q.whereArgs, v...)
	return q
}

// WithWhere 携带条件
func (q *Query) WithWhere(stmts []string, args []interface{}) *Query {
	q.whereStmt = append(q.whereStmt, stmts...)
	q.whereArgs = append(q.whereArgs, args...)
	return q
}

// Limit 这周Limit
func (q *Query) Limit(offset int64, limit uint) *Query {
	q.limitStmt = "LIMIT ?,? "
	q.limitArgs = append(q.limitArgs, offset, limit)
	return q
}

// Order todo
func (q *Query) Order(d string) *Query {
	q.order = fmt.Sprintf("ORDER BY %s ", d)
	return q
}

// Desc todo
func (q *Query) Desc() *Query {
	q.order = fmt.Sprintf("%s DESC ", strings.TrimSpace(q.order))
	return q
}

func (q *Query) whereBuild() string {
	return "WHERE " + strings.Join(q.whereStmt, " AND ") + " "
}

// WhereArgs where 语句的参数
func (q *Query) WhereArgs() []interface{} {
	return q.whereArgs
}

// WhereStmt where条件列表
func (q *Query) WhereStmt() []string {
	return q.whereStmt
}

// Build 组件SQL
func (q *Query) Build() (stmt string, args []interface{}) {
	stmt = q.query + " " + q.whereBuild() + q.order + q.limitStmt + ";"

	args = append(args, q.whereArgs...)
	args = append(args, q.limitArgs...)
	return
}
