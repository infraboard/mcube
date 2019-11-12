package sqlpool

import "database/sql"

type poolSet map[string]Pool

// NewDefaultPool 初始化一个默认sqlpool
func NewDefaultPool() Pool {
	return &defaultPool{}
}

type defaultPool struct {
}

func (p *defaultPool) AddStmtCachedSQL(key, sqlStr string) {
	return
}

func (p *defaultPool) AddSQL(key, sqlStr string) {

}

func (p *defaultPool) GetSQL(key string) string {
	return ""
}

func (p *defaultPool) Exec(key string) (sql.Result, error) {
	return nil, nil
}

func (p *defaultPool) Query(args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (p *defaultPool) QueryRow(args ...interface{}) *sql.Row {
	return nil
}
