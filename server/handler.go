package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	cb_logging "func/copybasta/logging"
)

type HttpHandler func(*http.Request) (int, interface{})

func JsonHandlerWrapper(handler HttpHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		statusCode, responseData := handler(req)
		if err := WriteJson(w, statusCode, responseData); err != nil {
			cb_logging.Error(req.Context(), "writeJson", err, &cb_logging.Data{"statusCode": statusCode, "responseData": responseData})
		}
	}
}

func HealthCheckHandler(router *mux.Router) HttpHandler {
	return func(req *http.Request) (int, interface{}) {
		routes, err := listRoutes(router)
		if err != nil {
			cb_logging.Error(req.Context(), "listRoutes", err, nil)
			return http.StatusInternalServerError, nil
		}
		return http.StatusOK, map[string]interface{}{
			"status":          "ok",
			"availableRoutes": routes,
		}
	}
}

func NotFoundHandler(router *mux.Router) HttpHandler {
	return func(req *http.Request) (int, interface{}) {
		routes, err := listRoutes(router)
		if err != nil {
			cb_logging.Error(req.Context(), "listRoutes", err, nil)
			return http.StatusInternalServerError, nil
		}
		return http.StatusNotFound, map[string]interface{}{"requestedPath": req.URL.Path, "availableRoutes": routes}
	}
}

func listRoutes(r *mux.Router) (map[string]string, error) {
	routeList := make(map[string]string)
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		methods, err := route.GetMethods()
		if err != nil {
			return err
		}

		routeList[route.GetName()] = fmt.Sprintf("%s [%s]", pathTemplate, strings.Join(methods, ","))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return routeList, nil
}
