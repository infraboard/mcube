package mongo_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/infraboard/mcube/ioc/apps/oss"
)

func TestUploadFile(t *testing.T) {
	buf := io.NopCloser(strings.NewReader("test log 1222222sdfsdfsdsfdf"))
	defer buf.Close()

	req := oss.NewUploadFileRequest("task_log", "test.log", buf)
	err := impl.UploadFile(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDownload(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	req := oss.NewDownloadFileRequest("task_log", "test.log", buf)
	err := impl.Download(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(buf.String())
}
