package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/infraboard/mcube/exception"
)

// Failed use to response error messge
func Failed(w http.ResponseWriter, err error) {
	var (
		errCode  int
		httpCode int
		ns       string
		reason   string
	)

	switch t := err.(type) {
	case exception.APIException:
		errCode = t.ErrorCode()
		reason = t.Reason()
		ns = t.Namespace()
	default:
		errCode = exception.UnKnownException
	}

	// 统一使用业务code, http code固定200
	httpCode = http.StatusOK

	resp := Data{
		Code:      &errCode,
		Namespace: ns,
		Reason:    reason,
		Message:   err.Error(),
	}

	// set response heanders
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Err-Code", fmt.Sprintf("%d", errCode))

	// if marshal json error, use string to response
	respByt, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMSG := fmt.Sprintf(`{"status":"error", "message": "encoding to json error, %s"}`, err)
		w.Write([]byte(errMSG))
		return
	}

	w.WriteHeader(httpCode)
	w.Write(respByt)
	return
}

// Success use to response success data
func Success(w http.ResponseWriter, data interface{}) {
	c := 0
	resp := Data{
		Code:    &c,
		Message: "",
		Data:    data,
	}

	// set response heanders
	w.Header().Set("Content-Type", "application/json")

	respBytes, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMSG := fmt.Sprintf(`{"status":"error", "message": "encoding to json error, %s"}`, err)
		w.Write([]byte(errMSG))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
	return
}
