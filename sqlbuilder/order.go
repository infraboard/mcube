package sqlbuilder

import "fmt"

type OrderStmt struct {
	order string
	b     *Builder
}

func (j *OrderStmt) Desc() *Builder {
	orderStmt := fmt.Sprintf("%s DESC", j.order)
	j.b.order = append(j.b.order, orderStmt)
	return j.b
}

func (j *OrderStmt) Asc() *Builder {
	orderStmt := fmt.Sprintf("%s ASC", j.order)
	j.b.order = append(j.b.order, orderStmt)
	return j.b
}
