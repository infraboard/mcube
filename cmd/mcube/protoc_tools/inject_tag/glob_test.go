package inject_tag

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestGlob(t *testing.T) {
	files, _ := filepath.Glob("../../../../pb/*/*.pb.go")
	fmt.Println(files)
}
