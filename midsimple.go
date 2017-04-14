package midsimple

import "net/http"

// MidSimple simple middleware manager
type MidSimple struct {
	list []func(http.Handler) http.Handler
}

// New creates an instance of MidSimple middleware manager
func New() *MidSimple {
	return new(MidSimple)
}

type handlerFuncWrapper struct {
	hfunc http.HandlerFunc
}

// ServeHTTP runs wrapped handler function
func (wrapper *handlerFuncWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wrapper.hfunc(w, r)
}

// Use adds middleware to the queue
func (m *MidSimple) Use(middleware ...func(http.Handler) http.Handler) *MidSimple {
	m.list = append(m.list, middleware...)
	return m
}

// Wrap wraps the handler with middlewares from queue and returns wrapper version
func (m *MidSimple) Wrap(handler http.Handler) http.Handler {
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
