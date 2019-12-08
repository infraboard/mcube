package ftime_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/infraboard/mcube/types/ftime"
	"github.com/stretchr/testify/assert"
)

type test struct {
	CreateAt ftime.Time `json:"create_at,omitempty"`
}

func TestJSONMashal(t *testing.T) {
	shoud := assert.New(t)

	now := time.Now()
	test := &test{
		CreateAt: ftime.T(now),
	}

	fj, err := json.Marshal(test)
	expected := fmt.Sprintf(`{"create_at":%d}`, now.UnixNano()/1000000)

	if shoud.NoError(err) {
		shoud.Equal(expected, string(fj))
	}
}

func TestJSONUnMashal(t *testing.T) {
	shoud := assert.New(t)

	now := time.Now()
	ts := now.UnixNano() / 1000000
	expected := fmt.Sprintf(`{"create_at":%d}`, ts)

	test := new(test)

	err := json.Unmarshal([]byte(expected), test)
	if shoud.NoError(err) {
		shoud.Equal(ts, test.CreateAt.T().UnixNano()/1000000)
	}

}
