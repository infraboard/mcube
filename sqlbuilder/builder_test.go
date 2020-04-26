package sqlbuilder_test

import (
	"testing"

	"github.com/infraboard/mcube/sqlbuilder"
	"github.com/stretchr/testify/assert"
)

func TestQueryBuild(t *testing.T) {
	should := assert.New(t)
	q := sqlbuilder.NewQuery("SELECT * FROM t")
	qstmt, args := q.Where("t.a = ? AND t.c = ? AND t.d LIKE ?", "one", "two", "three").Desc("t.create_at").Limit(0, 20).Build()
	should.Equal(qstmt, "SELECT * FROM t WHERE t.a = ? AND t.c = ? AND t.d LIKE ? ORDER BY t.create_at DESC LIMIT ?,? ;")
	should.Equal(args, []interface{}{"one", "two", "three", int64(0), uint(20)})
	should.Equal(q.WhereStmt(), []string{"t.a = ? AND t.c = ? AND t.d LIKE ?"})
	should.Equal(q.WhereArgs(), []interface{}{"one", "two", "three"})

	c := sqlbuilder.NewQuery("SELECT COUNT(*) FROM t")
	cstmt, args := c.WithWhere(q.WhereStmt(), q.WhereArgs()).Build()
	should.Equal(cstmt, "SELECT COUNT(*) FROM t WHERE t.a = ? AND t.c = ? AND t.d LIKE ? ;")
}
