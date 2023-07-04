package response_test

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/infraboard/mcube/http/response"
	"github.com/stretchr/testify/require"
)

type hijackableResponse struct {
	Hijacked bool
}

func newHijackableResponse() *hijackableResponse {
	return &hijackableResponse{}
}

func (h *hijackableResponse) Header() http.Header           { return nil }
func (h *hijackableResponse) Write(buf []byte) (int, error) { return 0, nil }
func (h *hijackableResponse) WriteHeader(code int)          {}
func (h *hijackableResponse) Flush()                        {}
func (h *hijackableResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h.Hijacked = true
	return nil, nil, nil
}

func TestResponseWriterBeforeWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := response.NewResponse(rec)

	should := require.New(t)
	should.Equal(rw.Status(), 0)
	should.Equal(rw.Written(), false)
}

func TestResponseWriterBeforeFuncHasAccessToStatus(t *testing.T) {
	var status int

	rec := httptest.NewRecorder()
	rw := response.NewResponse(rec)

	rw.Before(func(w response.Response) {
		status = w.Status()
	})
	rw.WriteHeader(http.StatusCreated)

	should := require.New(t)
	should.Equal(status, http.StatusCreated)
}

func TestResponseWriterWritingString(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := response.NewResponse(rec)

	rw.Write([]byte("Hello world"))

	should := require.New(t)
	should.Equal(rec.Code, rw.Status())
	should.Equal(rec.Body.String(), "Hello world")
	should.Equal(rw.Status(), http.StatusOK)
	should.Equal(rw.Size(), 11)
	should.Equal(rw.Written(), true)
}

func TestResponseWriterWritingStrings(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := response.NewResponse(rec)

	rw.Write([]byte("Hello world"))
	rw.Write([]byte("foo bar bat baz"))

	should := require.New(t)
	should.Equal(rec.Code, rw.Status())
	should.Equal(rec.Body.String(), "Hello worldfoo bar bat baz")
	should.Equal(rw.Status(), http.StatusOK)
	should.Equal(rw.Size(), 26)
}

func TestResponseWriterWritingHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := response.NewResponse(rec)

	rw.WriteHeader(http.StatusNotFound)

	should := require.New(t)
	should.Equal(rec.Code, rw.Status())
	should.Equal(rec.Body.String(), "")
	should.Equal(rw.Status(), http.StatusNotFound)
	should.Equal(rw.Size(), 0)
}

func TestResponseWriterBefore(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := response.NewResponse(rec)
	result := ""

	rw.Before(func(response.Response) {
		result += "foo"
	})
	rw.Before(func(response.Response) {
		result += "bar"
	})

	rw.WriteHeader(http.StatusNotFound)

	should := require.New(t)
	should.Equal(rec.Code, rw.Status())
	should.Equal(rec.Body.String(), "")
	should.Equal(rw.Status(), http.StatusNotFound)
	should.Equal(rw.Size(), 0)
	should.Equal(result, "barfoo")
}

func TestResponseWriterHijack(t *testing.T) {
	should := require.New(t)

	hijackable := newHijackableResponse()
	rw := response.NewResponse(hijackable)
	hijacker, ok := rw.(http.Hijacker)
	should.Equal(ok, true)

	_, _, err := hijacker.Hijack()
	if err != nil {
		t.Error(err)
	}
	should.Equal(hijackable.Hijacked, true)
}

func TestResponseWriteHijackNotOK(t *testing.T) {
	should := require.New(t)

	hijackable := new(http.ResponseWriter)
	rw := response.NewResponse(*hijackable)
	hijacker, ok := rw.(http.Hijacker)
	should.Equal(ok, true)

	_, _, err := hijacker.Hijack()
	should.EqualError(err, "the ResponseWriter doesn't support the Hijacker interface")
}

func TestResponseWriterFlusher(t *testing.T) {
	should := require.New(t)

	rec := httptest.NewRecorder()
	rw := response.NewResponse(rec)

	_, ok := rw.(http.Flusher)
	should.Equal(ok, true)
}

func TestResponseWriter_Flush_marksWritten(t *testing.T) {
	should := require.New(t)
	rec := httptest.NewRecorder()
	rw := response.NewResponse(rec)

	rw.Flush()
	should.Equal(rw.Status(), http.StatusOK)
	should.Equal(rw.Written(), true)
}
