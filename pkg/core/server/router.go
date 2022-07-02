package server

import "net/http"

type Router struct {
	middlewares []func(http.Handler) http.Handler
}

func NewRouter() *Router {
	return &Router{}
}

func (m *Router) Handle(pattern string, handler http.Handler) {
	http.Handle(pattern, m.applyMiddlewares(handler))
}

func (m *Router) AddMiddleware(middleware func(http.Handler) http.Handler) {
	m.middlewares = append(m.middlewares, middleware)
}

func (m *Router) applyMiddlewares(handler http.Handler) http.Handler {
	for _, middleware := range m.middlewares {
		handler = middleware(handler)
	}

	return handler
}
