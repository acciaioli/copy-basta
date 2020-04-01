package server

import (
	"net/http"
	"time"

	cb_logging "func/copybasta/logging"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		reqURL := req.URL.String()
		reqMethod := req.Method
		next.ServeHTTP(w, req)
		reqDuration := time.Since(start) / time.Millisecond
		cb_logging.Info(req.Context(), "request", &cb_logging.Data{"url": reqURL, "method": reqMethod, "duration": int(reqDuration)})
	})
}
