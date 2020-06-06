package enum_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/infraboard/mcube/cmd/gencheck/enum"
	"github.com/stretchr/testify/assert"
)

type TestData struct {
	ID     string      `json:"id,omitempty"`
	Status enum.Status `json:"status"`
}

func TestJ(t *testing.T) {
	should := assert.New(t)
	data := TestData{
		ID:     "xxxx",
		Status: enum.Running,
	}

	jd, err := json.Marshal(data)
	should.NoError(err)
	fmt.Println(string(jd))

	td := &TestData{}
	json.Unmarshal(jd, td)
	fmt.Println(td)
	fmt.Println(td.Status.Is(enum.Running))
}
