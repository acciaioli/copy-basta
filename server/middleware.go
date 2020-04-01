package server

import (
	"net/http"
	"time"

	"github.com/spin14/copy-basta/logging"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		reqURL := req.URL.String()
		reqMethod := req.Method
		next.ServeHTTP(w, req)
		reqDuration := time.Since(start) / time.Millisecond
		logging.Info(req.Context(), "request", &logging.Data{"url": reqURL, "method": reqMethod, "duration": int(reqDuration)})
	})
}
