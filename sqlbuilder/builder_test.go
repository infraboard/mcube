package sqlbuilder_test

import (
	"testing"

	"github.com/infraboard/mcube/sqlbuilder"
	"github.com/stretchr/testify/assert"
)

func TestQueryBuild(t *testing.T) {
	should := assert.New(t)
	q := sqlbuilder.NewQuery("SELECT * FROM t")
	q.Where("t.a = ? AND t.c = ? AND t.d LIKE ?", "one", "two", "three")
	q.Desc("t.create_at")
	q.Limit(0, 20)
	stmt, args := q.Build()
	should.Equal(stmt, "SELECT * FROM t WHERE t.a = ? AND t.c = ? AND t.d LIKE ? ORDER BY t.create_at DESC LIMIT ?,? ;")
	should.Equal(args, []interface{}{"one", "two", "three", int64(0), uint(20)})
	should.Equal(q.WhereStmt(), []string{"t.a = ? AND t.c = ? AND t.d LIKE ?"})
	should.Equal(q.WhereArgs(), []interface{}{"one", "two", "three"})
}
