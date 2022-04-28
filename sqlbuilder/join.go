package sqlbuilder

import "fmt"

type JoinStmt struct {
	join string
	b    *Builder
}

func (j *JoinStmt) ON(cond string) *Builder {
	joinStmt := fmt.Sprintf("%s ON %s", j.join, cond)
	j.b.join = append(j.b.join, joinStmt)
	return j.b
}
