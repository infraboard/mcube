package redis_test

import (

	"testing"

	"github.com/stretchr/testify/require"

	"github.com/infraboard/mcube/cache"
	"github.com/infraboard/mcube/cache/redis"
)

var (
	redisConf = `{"address": "127.0.0.1:6379", "db": "", "password": ""}`
)

type adapterSuit struct {
	adapter cache.Cache
	testKey string
	testVal string
}

func (a *adapterSuit) SetUp() {
	adapter := redis.NewCache()
	if err := adapter.Config(redisConf); err != nil {
		panic(err)
	}

	a.adapter = adapter
	a.testKey = "testkey01"
	a.testVal = "testval01"
}

func (a *adapterSuit) TearDown() {
	a.adapter.Close()
}

func TestRedisAdapterSuit(t *testing.T) {
	suit := new(adapterSuit)
	suit.SetUp()
	defer suit.TearDown()

	t.Run("PutOK", testPutOK(suit))
	t.Run("ExistOK", testExistOK(suit))
	t.Run("ExistNotOK", testExistNotOK(suit))
	t.Run("DelOK", testDelOK(suit))
}

func testPutOK(a *adapterSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)

		err := a.adapter.Put(a.testKey, a.testVal)
		should.NoError(err)
	}
}

func testGetOK(a *adapterSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)

		val := ""
		a.adapter.Get(a.testKey, &val)
		should.Equal(a.testVal, val)
	}
}

func testExistOK(a *adapterSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)

		ok := a.adapter.IsExist(a.testKey)
		should.Equal(true, ok)
	}
}

func testExistNotOK(a *adapterSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)

		ok := a.adapter.IsExist("not exist key")
		should.Equal(false, ok)
	}
}

func testDelOK(a *adapterSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)

		err := a.adapter.Delete(a.testKey)
		should.NoError(err)
	}
}