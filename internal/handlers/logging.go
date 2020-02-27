package handlers

import (
	"fantomrocks-api/internal/services"
	"net/http"
)

// Create new logging middleware HTTP handler.
func LoggingHandler(log services.Logger, h http.Handler) http.Handler {
	// make new handler using closure
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// put out the log
		log.Debugf("Serving %s %s (%s from [%s]) %s", r.Proto, r.Method, r.UserAgent(), r.RemoteAddr, r.URL)

		// pass the request down the chain
		h.ServeHTTP(w, r)
	})
}
