package recovery_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infraboard/mcube/v2/http/middleware/recovery"
	"github.com/infraboard/mcube/v2/http/router/httprouter"
	"github.com/stretchr/testify/require"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	panic("recovery test")
}

func Test_Recovery(t *testing.T) {
	should := require.New(t)

	router := httprouter.New()

	rm := recovery.New()
	router.Use(rm)
	router.Handle("GET", "/", indexHandler)

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	should.NoError(err)

	router.ServeHTTP(recorder, req)
	should.Equal(recorder.Code, http.StatusInternalServerError)
}

func init() {
}
