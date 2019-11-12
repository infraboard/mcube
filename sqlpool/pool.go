// 集中管理SQL语句, 提供给SQL审计系统使用
// 集中管理SQL的执行, 如果有慢查询可以 提供给metric系统
// 集中管理SQL的执行统计, 每个sql语句执行的次数

package sqlpool

import (
	"database/sql"
)

// Pool sql语句池
type Pool interface {
	AddSQL(key, sqlStr string)
	AddStmtCachedSQL(key, sqlStr string)
	GetSQL(key string) string
	Exec(key string) (sql.Result, error)
	Query(args ...interface{}) (*sql.Rows, error)
	QueryRow(args ...interface{}) *sql.Row
}

// PoolStati sqlpool 统计信息
type PoolStati interface {
	ALLSQL() []string
	CallCount() map[string]int64
}

// PoolProvider pool manager
type PoolProvider interface {
	NewPool(name string) Pool
	GetPoolStati(name string) PoolStati
}
