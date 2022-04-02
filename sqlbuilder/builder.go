package sqlbuilder

import (
	"fmt"
	"strings"
)

// NewQuery new 实例
func NewQuery(querySQL string, args ...interface{}) *Builder {
	return NewBuilder(querySQL, args...)
}

// NewBuilder new 实例
func NewBuilder(querySQL string, args ...interface{}) *Builder {
	return &Builder{
		base:      fmt.Sprintf(querySQL, args...),
		whereStmt: []string{},
		whereArgs: []interface{}{},
		limitArgs: []interface{}{},
	}
}

// Query 查询mysql数据库
type Builder struct {
	base       string
	join       []string
	whereStmt  []string
	whereArgs  []interface{}
	limitStmt  string
	limitArgs  []interface{}
	order      string
	groupBy    string
	havingStmt []string
	havingArgs []interface{}
}

// Where 添加参数
func (q *Builder) Where(stmt string, v ...interface{}) *Builder {
	q.whereStmt = append(q.whereStmt, stmt)
	q.whereArgs = append(q.whereArgs, v...)
	return q
}

// WithWhere 携带条件
func (q *Builder) WithWhere(stmts []string, args []interface{}) *Builder {
	q.whereStmt = append(q.whereStmt, stmts...)
	q.whereArgs = append(q.whereArgs, args...)
	return q
}

// Having 添加参数
func (q *Builder) Having(stmt string, v ...interface{}) *Builder {
	q.havingStmt = append(q.havingStmt, stmt)
	q.havingArgs = append(q.havingArgs, v...)
	return q
}

// WithHaving 携带条件
func (q *Builder) WithHaving(stmts []string, args []interface{}) *Builder {
	q.havingStmt = append(q.havingStmt, stmts...)
	q.havingArgs = append(q.havingArgs, args...)
	return q
}

// Limit Limit
func (q *Builder) Limit(offset int64, limit uint) *Builder {
	q.limitStmt = "LIMIT ?,? "
	q.limitArgs = append(q.limitArgs, offset, limit)
	return q
}

// Order todo
func (q *Builder) Order(d string) *Builder {
	q.order = fmt.Sprintf("ORDER BY %s ", d)
	return q
}

// LeftJoin("xxxx").ON("xxx")
func (q *Builder) LeftJoin(j string) *JoinStmt {
	return q.joinStmt(fmt.Sprintf("LEFT JOIN %s", j))
}

// RIGHT("xxxx").ON("xxx")
func (q *Builder) RightJoin(j string) *JoinStmt {
	return q.joinStmt(fmt.Sprintf("RIGHT JOIN %s", j))
}

// Desc todo
func (q *Builder) Desc() *Builder {
	q.order = fmt.Sprintf("%s DESC ", strings.TrimSpace(q.order))
	return q
}

// GroupBy todo
func (q *Builder) GroupBy(d string) *Builder {
	q.groupBy = fmt.Sprintf("GROUP BY %s ", strings.TrimSpace(d))
	return q
}

func (q *Builder) whereBuild() string {
	if len(q.whereStmt) == 0 {
		return ""
	}

	return "WHERE " + strings.Join(q.whereStmt, " AND ") + " "
}

// WhereArgs where 语句的参数
func (q *Builder) WhereArgs() []interface{} {
	return q.whereArgs
}

// WhereStmt where条件列表
func (q *Builder) WhereStmt() []string {
	return q.whereStmt
}

func (q *Builder) havingBuild() string {
	if len(q.havingStmt) == 0 {
		return ""
	}

	return "HAVING " + strings.Join(q.havingStmt, " AND ") + " "
}

// HavingArgs where 语句的参数
func (q *Builder) HavingArgs() []interface{} {
	return q.havingArgs
}

// HavingStmt where条件列表
func (q *Builder) HavingStmt() []string {
	return q.havingStmt
}

// Build 组件SQL
func (q *Builder) BuildQuery() (stmt string, args []interface{}) {
	stmt = q.base + " " + q.joinBuild() + q.whereBuild() + q.groupBy + q.havingBuild() + q.order + q.limitStmt + ";"

	args = append(args, q.whereArgs...)
	args = append(args, q.havingArgs...)
	args = append(args, q.limitArgs...)
	return
}

// 提供 base 替换
// 只提供条件语句 where, group, having
func (q *Builder) BuildFromNewBase(base string) (stmt string, args []interface{}) {
	stmt = base + " " + q.joinBuild() + q.whereBuild() + q.groupBy + q.havingBuild() + ";"
	args = append(args, q.whereArgs...)
	args = append(args, q.havingArgs...)
	return stmt, args
}

func (q *Builder) BuildCount() (stmt string, args []interface{}) {
	uppterQuery := strings.ToUpper(q.base)
	start := strings.Index(uppterQuery, "SELECT")
	end := strings.Index(uppterQuery, "FROM")

	countQuery := fmt.Sprintf("%s %s %s", q.base[:start+6], "COUNT(*)", q.base[end:])
	return q.BuildFromNewBase(countQuery)
}

func (q *Builder) joinBuild() string {
	return strings.Join(q.join, " ") + " "
}

func (q *Builder) joinStmt(joinSQL string) *JoinStmt {
	return &JoinStmt{
		join: joinSQL,
		b:    q,
	}
}

type JoinStmt struct {
	join string
	b    *Builder
}

func (j *JoinStmt) ON(cond string) *Builder {
	joinStmt := fmt.Sprintf("%s ON %s", j.join, cond)
	j.b.join = append(j.b.join, joinStmt)
	return j.b
}
