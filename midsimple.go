package midsimple

import "net/http"

type MidSimple struct {
	list []func(http.Handler) http.Handler
}

func New() *MidSimple {
	return new(MidSimple)
}

type handlerFuncWrapper struct {
	hfunc http.HandlerFunc
}

func (wrapper *handlerFuncWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wrapper.hfunc(w, r)
}

func (m *MidSimple) Use(middleware func(http.Handler) http.Handler) {
	m.list = append(m.list, middleware)
}

func (m *MidSimple) Wrap(handler http.Handler) http.Handler {
	res := handler
	for i := len(m.list) - 1; i >= 0; i-- {
		res = m.list[i](res)
	}
	return res
}

func (m *MidSimple) WrapFunc(handlerFunc http.HandlerFunc) http.Handler {
	handler := http.Handler(&handlerFuncWrapper{hfunc: handlerFunc})
	return m.Wrap(handler)
}

func (m *MidSimple) WrapRevert(handler http.Handler) http.Handler {
	res := handler
	for i := range m.list {
		res = m.list[i](res)
	}
	return res
}

func (m *MidSimple) WrapRevertFunc(handlerFunc http.HandlerFunc) http.Handler {
	handler := http.Handler(&handlerFuncWrapper{hfunc: handlerFunc})
	return m.WrapRevert(handler)
}
