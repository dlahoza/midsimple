package midsimple

import (
	"net/http"
	"sync"
)

// MidSimple simple middleware manager
type MidSimple struct {
	list []func(http.Handler) http.Handler
	sync.Mutex
}

// New creates an instance of MidSimple middleware manager and adds middlewares if any
func New(middleware ...func(http.Handler) http.Handler) *MidSimple {
	ms := new(MidSimple)
	if len(middleware) > 0 {
		ms.Use(middleware...)
	}
	return ms
}

type handlerFuncWrapper struct {
	hfunc http.HandlerFunc
}

// ServeHTTP runs wrapped handler function
func (wrapper *handlerFuncWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wrapper.hfunc(w, r)
}

// Use adds middlewares to the queue
func (m *MidSimple) Use(middleware ...func(http.Handler) http.Handler) *MidSimple {
	m.Lock()
	defer m.Unlock()
	m.list = append(m.list, middleware...)
	return m
}

// Reset clears list of middlewares and allows use reuse it
func (m *MidSimple) Reset() *MidSimple {
	m.Lock()
	defer m.Unlock()
	m.list = m.list[:0]
	return m
}

// Wrap wraps the handler with middlewares from queue and returns wrapper version
func (m *MidSimple) Wrap(handler http.Handler) http.Handler {
	m.Lock()
	defer m.Unlock()
	res := handler
	for i := len(m.list) - 1; i >= 0; i-- {
		res = m.list[i](res)
	}
	return res
}

// WrapFunc wraps HandlerFunc with middlewares from queue and returns wrapper version
func (m *MidSimple) WrapFunc(handlerFunc http.HandlerFunc) http.Handler {

	handler := http.Handler(&handlerFuncWrapper{hfunc: handlerFunc})
	return m.Wrap(handler)
}

// WrapRevert wraps the handler with middlewares from queue but in reverted sequence and returns wrapper version
func (m *MidSimple) WrapRevert(handler http.Handler) http.Handler {
	m.Lock()
	defer m.Unlock()
	res := handler
	for i := range m.list {
		res = m.list[i](res)
	}
	return res
}

// WrapRevertFunc wraps the HandlerFunc with middlewares from queue but in reverted sequence and returns wrapper version
func (m *MidSimple) WrapRevertFunc(handlerFunc http.HandlerFunc) http.Handler {
	handler := http.Handler(&handlerFuncWrapper{hfunc: handlerFunc})
	return m.WrapRevert(handler)
}
