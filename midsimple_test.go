package midsimple

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestMidSimple_New(t *testing.T) {
	ms := New(newTestMiddleware("1"), newTestMiddleware("2"))
	resp := httptest.NewRecorder()
	ms.Wrap(&testHandler{result: "handler"}).ServeHTTP(resp, nil)
	assert.Equal(t, []byte("12handler"), resp.Body.Bytes())
}

func TestMidSimple_UseReturnsMidSimple(t *testing.T) {
	ms := New()
	returnedMs := ms.Use(newTestMiddleware("1"))
	assert.Equal(t, ms, returnedMs)
}

func TestMidSimple_Reset(t *testing.T) {
	ms := New()
	ms.Use(newTestMiddleware("1"))
	ms.Reset().
		Use(newTestMiddleware("2"))
	resp := httptest.NewRecorder()
	ms.Wrap(&testHandler{result: "handler"}).ServeHTTP(resp, nil)
	assert.Equal(t, []byte("2handler"), resp.Body.Bytes())
}

func TestHandlerFuncWrapper(t *testing.T) {
	var success bool
	h := handlerFuncWrapper{hfunc: func(http.ResponseWriter, *http.Request) { success = true }}
	h.ServeHTTP(nil, nil)
	assert.True(t, success)
}

func TestMidSimple_WrapHandlerRegular(t *testing.T) {
	ms := New()
	ms.Use(newTestMiddleware("1"))
	ms.Use(newTestMiddleware("2"), newTestMiddleware("3"), newTestMiddleware("4")).
		Use(newTestMiddleware("5"))
	ms.Use(newTestMiddleware("6"))
	resp := httptest.NewRecorder()
	ms.Wrap(&testHandler{result: "handler"}).ServeHTTP(resp, nil)
	assert.Equal(t, []byte("123456handler"), resp.Body.Bytes())
}

func TestMidSimple_WrapHandlerRevert(t *testing.T) {
	ms := New()
	ms.Use(newTestMiddleware("1"))
	ms.Use(newTestMiddleware("2"), newTestMiddleware("3"), newTestMiddleware("4")).
		Use(newTestMiddleware("5"))
	ms.Use(newTestMiddleware("6"))
	resp := httptest.NewRecorder()
	ms.WrapRevert(&testHandler{result: "handler"}).ServeHTTP(resp, nil)
	assert.Equal(t, []byte("654321handler"), resp.Body.Bytes())
}

func TestMidSimple_WrapHandlerFuncRegular(t *testing.T) {
	ms := New()
	ms.Use(newTestMiddleware("1"))
	ms.Use(newTestMiddleware("2"), newTestMiddleware("3"), newTestMiddleware("4")).
		Use(newTestMiddleware("5"))
	ms.Use(newTestMiddleware("6"))

	resp := httptest.NewRecorder()
	ms.WrapFunc((&testHandler{result: "handler"}).ServeHTTP).ServeHTTP(resp, nil)
	assert.Equal(t, []byte("123456handler"), resp.Body.Bytes())
}

func TestMidSimple_WrapHandlerFuncRevert(t *testing.T) {
	ms := New()
	ms.Use(newTestMiddleware("1"))
	ms.Use(newTestMiddleware("2"), newTestMiddleware("3"), newTestMiddleware("4")).
		Use(newTestMiddleware("5"))
	ms.Use(newTestMiddleware("6"))

	resp := httptest.NewRecorder()
	ms.WrapRevertFunc((&testHandler{result: "handler"}).ServeHTTP).ServeHTTP(resp, nil)
	assert.Equal(t, []byte("654321handler"), resp.Body.Bytes())
}
