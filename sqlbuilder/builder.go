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
		initArgs:  []interface{}{},
		whereStmt: []string{},
		whereArgs: []interface{}{},
		limitArgs: []interface{}{},
	}
}

// Query 查询mysql数据库
type Builder struct {
	base       string
	initArgs   []interface{}
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
func (b *Builder) Where(stmt string, v ...interface{}) *Builder {
	b.whereStmt = append(b.whereStmt, stmt)
	b.whereArgs = append(b.whereArgs, v...)
	return b
}

// WithWhere 携带条件
func (b *Builder) WithWhere(stmts []string, args []interface{}) *Builder {
	b.whereStmt = append(b.whereStmt, stmts...)
	b.whereArgs = append(b.whereArgs, args...)
	return b
}

// Having 添加参数
func (b *Builder) Having(stmt string, v ...interface{}) *Builder {
	b.havingStmt = append(b.havingStmt, stmt)
	b.havingArgs = append(b.havingArgs, v...)
	return b
}

// WithHaving 携带条件
func (b *Builder) WithHaving(stmts []string, args []interface{}) *Builder {
	b.havingStmt = append(b.havingStmt, stmts...)
	b.havingArgs = append(b.havingArgs, args...)
	return b
}

// Limit Limit
func (b *Builder) Limit(offset int64, limit uint) *Builder {
	b.limitStmt = "LIMIT ?,? "
	b.limitArgs = append(b.limitArgs, offset, limit)
	return b
}

// Order todo
func (b *Builder) Order(d string) *Builder {
	b.order = fmt.Sprintf("ORDER BY %s ", d)
	return b
}

// LeftJoin("xxxx").ON("xxx")
func (b *Builder) LeftJoin(j string) *JoinStmt {
	return b.joinStmt(fmt.Sprintf("LEFT JOIN %s", j))
}

// RIGHT("xxxx").ON("xxx")
func (b *Builder) RightJoin(j string) *JoinStmt {
	return b.joinStmt(fmt.Sprintf("RIGHT JOIN %s", j))
}

// Desc todo
func (b *Builder) Desc() *Builder {
	b.order = fmt.Sprintf("%s DESC ", strings.TrimSpace(b.order))
	return b
}

// GroupBy todo
func (b *Builder) GroupBy(d string) *Builder {
	b.groupBy = fmt.Sprintf("GROUP BY %s ", strings.TrimSpace(d))
	return b
}

func (b *Builder) whereBuild() string {
	if len(b.whereStmt) == 0 {
		return ""
	}

	return "WHERE " + strings.Join(b.whereStmt, " AND ") + " "
}

// WhereArgs where 语句的参数
func (b *Builder) WhereArgs() []interface{} {
	return b.whereArgs
}

// WhereStmt where条件列表
func (b *Builder) WhereStmt() []string {
	return b.whereStmt
}

func (b *Builder) havingBuild() string {
	if len(b.havingStmt) == 0 {
		return ""
	}

	return "HAVING " + strings.Join(b.havingStmt, " AND ") + " "
}

// HavingArgs where 语句的参数
func (b *Builder) HavingArgs() []interface{} {
	return b.havingArgs
}

// HavingStmt where条件列表
func (b *Builder) HavingStmt() []string {
	return b.havingStmt
}

// Build 组件SQL
func (b *Builder) Build() (stmt string, args []interface{}) {
	stmt = b.base + " " + b.joinBuild() + b.whereBuild() + b.groupBy + b.havingBuild() + b.order + b.limitStmt + ";"

	args = append(args, b.whereArgs...)
	args = append(args, b.havingArgs...)
	args = append(args, b.limitArgs...)
	return
}

// AddArgs 添加初始化Args
func (b *Builder) AddArgs(args ...interface{}) *Builder {
	b.initArgs = append(b.initArgs, args...)
	return b
}

// DEPRECEATED Build 组件SQL
func (b *Builder) BuildQuery() (stmt string, args []interface{}) {
	return b.Build()
}

// 提供 base 替换
// 只提供条件语句 where, group, having
func (b *Builder) BuildFromNewBase(base string) (stmt string, args []interface{}) {
	stmt = base + " " + b.joinBuild() + b.whereBuild() + b.groupBy + b.havingBuild() + ";"
	args = append(args, b.whereArgs...)
	args = append(args, b.havingArgs...)
	return stmt, args
}

func (b *Builder) BuildCount() (stmt string, args []interface{}) {
	uppterQuery := strings.ToUpper(b.base)
	start := strings.Index(uppterQuery, "SELECT")
	end := strings.Index(uppterQuery, "FROM")

	countQuery := fmt.Sprintf("%s %s %s", b.base[:start+6], "COUNT(*)", b.base[end:])
	return b.BuildFromNewBase(countQuery)
}

func (b *Builder) joinBuild() string {
	return strings.Join(b.join, " ") + " "
}

func (b *Builder) joinStmt(joinSQL string) *JoinStmt {
	return &JoinStmt{
		join: joinSQL,
		b:    b,
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
