package cli_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infraboard/mcube/tools/cli"
)

const testFullResult = `├───dir1
│	└───file.txt (empty)
├───dir2
│	├───dir1
│	│	└───file.txt (empty)
│	└───file.txt (1b)
└───file.txt (empty)
`

func TestTreeFull(t *testing.T) {
	should := assert.New(t)
	out := new(bytes.Buffer)

	err := cli.Tree(out, "testdata", true)

	if should.NoError(err) {
		should.Equal(testFullResult, out.String())
	}
}

const testDirResult = `├───dir1
└───dir2
	└───dir1
`

func TestTreeDir(t *testing.T) {
	should := assert.New(t)
	out := new(bytes.Buffer)
	err := cli.Tree(out, "testdata", false)

	if should.NoError(err) {
		should.Equal(testDirResult, out.String())
	}
}
