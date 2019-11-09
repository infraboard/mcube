package request

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/infraboard/mcube/exception"
)

// CheckBody 读取Body当中的数据
func CheckBody(r *http.Request) ([]byte, error) {
	// 检测请求大小
	if r.ContentLength == 0 {
		return nil, exception.NewBadRequest("request body is empty")
	}
	if r.ContentLength > 20971520 {
		return nil, exception.NewBadRequest(
			"the body exceeding the maximum limit, max size 20M")
	}

	// 读取body数据
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		return nil, exception.NewBadRequest(
			fmt.Sprintf("read request body error, %s", err))
	}

	return body, nil
}

// GetObjFromReq todo
func GetObjFromReq(r *http.Request, v interface{}) error {
	body, err := CheckBody(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, v); err != nil {
		return err
	}

	return nil
}
