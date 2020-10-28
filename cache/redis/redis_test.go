package redis_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/infraboard/mcube/cache"
	"github.com/infraboard/mcube/cache/redis"
)

type adapterSuit struct {
	adapter cache.Cache
	testKey string
	testVal string
}

func (a *adapterSuit) SetUp() {
	adapter := redis.NewCache(redis.NewDefaultConfig())
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
	t.Run("GetOK", testGetOK(suit))
	t.Run("GetFailed", testGetFailed(suit))
	t.Run("KeysOK", testKeysOK(suit))
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

func testGetFailed(a *adapterSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)

		val := ""
		err := a.adapter.Get("xxx", &val)
		should.Equal("redis: nil", err.Error())
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

func testKeysOK(a *adapterSuit) func(t *testing.T) {
	return func(t *testing.T) {
		should := require.New(t)

		ks, err := a.adapter.ListKey(cache.NewListKeyRequest("testkey*", 3, 1))
		should.NoError(err)
		should.Equal([]string{a.testKey}, ks.Keys)
	}
}
