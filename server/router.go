package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
}

type RouterOption func(*Router) error

func RegisterMiddleware(middleware func(http.Handler) http.Handler) RouterOption {
	return func(router *Router) error {
		router.Use(middleware)
		return nil
	}
}

func RegisterHandler(name, path string, handler HttpHandler, methods ...string) RouterOption {
	return func(router *Router) error {
		if methods == nil {
			return errors.New("trying to register handler with no methods")
		}
		router.HandleFunc(fmt.Sprintf("/api/%s", path), JsonHandlerWrapper(handler)).Name(name).Methods(methods...)
		return nil
	}
}

func NewRouter(opts ...RouterOption) (*Router, error) {
	router := &Router{Router: mux.NewRouter()}
	router.Use(LoggerMiddleware)

	defaultOptions := []RouterOption{
		RegisterMiddleware(LoggerMiddleware),
		RegisterHandler("HealthCheck", "healthcheck", HealthCheckHandler(router.Router), http.MethodGet),
		func(router *Router) error {
			router.Router.NotFoundHandler = http.HandlerFunc(JsonHandlerWrapper(NotFoundHandler(router.Router)))
			return nil
		},
	}

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(router); err != nil {
			return nil, err
		}
	}

	return router, nil
}
