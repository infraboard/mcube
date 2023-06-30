package hash

import (
	"fmt"
	"hash/fnv"
	"strings"
)

// GenHashID hash id
func FnvHash(contents ...string) string {
	h := fnv.New64a()
	h.Write([]byte(strings.Join(contents, "")))
	return fmt.Sprintf("%x", h.Sum64())
}
