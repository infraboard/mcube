package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/infraboard/mcube/v2/exception"
)

// Failed use to response error messge
func Failed(w http.ResponseWriter, err error, opts ...Option) {
	var (
		errCode  int
		httpCode int
		ns       string
		reason   string
		data     any
		meta     any
	)

	switch t := err.(type) {
	case *exception.ApiException:
		errCode = t.ErrorCode()
		reason = t.GetReason()
		data = t.GetData()
		meta = t.Meta
		httpCode = t.GetHttpCode()
		ns = t.GetNamespace()
	default:
		errCode = exception.CODE_UNKNOWN
	}

	if httpCode == 0 {
		httpCode = http.StatusInternalServerError
	}

	resp := Data{
		Code:      &errCode,
		Namespace: ns,
		Reason:    reason,
		Message:   err.Error(),
		Data:      data,
		Meta:      meta,
	}

	for _, opt := range opts {
		opt.Apply(&resp)
	}

	// set response heanders
	w.Header().Set("Content-Type", "application/json")

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
}

// Success use to response success data
func Success(w http.ResponseWriter, data any, opts ...Option) {
	c := 0
	resp := Data{
		Code:    &c,
		Message: "",
		Data:    data,
	}

	for _, opt := range opts {
		opt.Apply(&resp)
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
}
