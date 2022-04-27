package sqlbuilder_test

import (
	"fmt"
	"testing"

	"github.com/infraboard/mcube/sqlbuilder"
	"github.com/stretchr/testify/assert"
)

func TestQueryBuild(t *testing.T) {
	should := assert.New(t)
	q := sqlbuilder.NewQuery("SELECT * FROM t")
	q.LeftJoin("table_join as tj").ON("tj.t_id = t.id")
	qstmt, args := q.Where("t.a = ? AND t.c = ? AND t.d LIKE ?", "one", "two", "three").GroupBy("t.group").Having("MAX(t.salary) > ?", 10).Order("t.create_at").Desc().Limit(0, 20).BuildQuery()
	should.Equal(qstmt, "SELECT * FROM t LEFT JOIN table_join as tj ON tj.t_id = t.id WHERE t.a = ? AND t.c = ? AND t.d LIKE ? GROUP BY t.group HAVING MAX(t.salary) > ? ORDER BY t.create_at DESC LIMIT ?,? ;")
	should.Equal(args, []interface{}{"one", "two", "three", 10, int64(0), uint(20)})
	should.Equal(q.WhereStmt(), []string{"t.a = ? AND t.c = ? AND t.d LIKE ?"})
	should.Equal(q.WhereArgs(), []interface{}{"one", "two", "three"})

	cstmt, args := q.BuildCount()
	should.Equal(cstmt, "SELECT COUNT(*) FROM t LEFT JOIN table_join as tj ON tj.t_id = t.id WHERE t.a = ? AND t.c = ? AND t.d LIKE ? GROUP BY t.group HAVING MAX(t.salary) > ? ;")
	should.Equal(args, []interface{}{"one", "two", "three", 10})
}

func TestSetBuild(t *testing.T) {
	should := assert.New(t)
	q := sqlbuilder.NewQuery("UPDATE t").Set("t.f1 = ?", "set1").Set("t.f2 = ?", "set2")
	q.LeftJoin("table_join as tj").ON("tj.t_id = t.id")
	cstmt, args := q.Build()
	fmt.Println(cstmt, args)
	should.Equal(cstmt, "UPDATE t LEFT JOIN table_join as tj ON tj.t_id = t.id SET t.f1 = ?,t.f2 = ? ;")
	should.Equal(args, []interface{}{"set1", "set2"})
}
