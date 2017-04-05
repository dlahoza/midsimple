package midsimple

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testMiddleware struct {
	result string
	next   http.Handler
}

func (m *testMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(m.result))
	m.next.ServeHTTP(w, r)
}

type testHandler struct {
	result string
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(h.result))
}

func newTestMiddleware(result string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &testMiddleware{
			result: result,
			next:   h,
		}
	}
}

func TestHandlerFuncWrapper(t *testing.T) {
	var success bool
	h := handlerFuncWrapper{hfunc: func(http.ResponseWriter, *http.Request) { success = true }}
	h.ServeHTTP(nil, nil)
	assert.True(t, success)
}

func TestMidSimple(t *testing.T) {
	ms := New()
	ms.Use(newTestMiddleware("1"))
	ms.Use(newTestMiddleware("2"))
	ms.Use(newTestMiddleware("3"))
	t.Run("WrapHandlerRegular", func(t *testing.T) {
		resp := httptest.NewRecorder()
		ms.Wrap(&testHandler{result: "handler"}).ServeHTTP(resp, nil)
		assert.Equal(t, []byte("123handler"), resp.Body.Bytes())
	})
	t.Run("WrapHandlerRevert", func(t *testing.T) {
		resp := httptest.NewRecorder()
		ms.WrapRevert(&testHandler{result: "handler"}).ServeHTTP(resp, nil)
		assert.Equal(t, []byte("321handler"), resp.Body.Bytes())
	})
	t.Run("WrapHandlerFuncRegular", func(t *testing.T) {
		resp := httptest.NewRecorder()
		ms.WrapFunc((&testHandler{result: "handler"}).ServeHTTP).ServeHTTP(resp, nil)
		assert.Equal(t, []byte("123handler"), resp.Body.Bytes())
	})
	t.Run("WrapHandlerFuncRevert", func(t *testing.T) {
		resp := httptest.NewRecorder()
		ms.WrapRevertFunc((&testHandler{result: "handler"}).ServeHTTP).ServeHTTP(resp, nil)
		assert.Equal(t, []byte("321handler"), resp.Body.Bytes())
	})
}
