package httpapi

import (
	"fmt"
	"log"
	"net/http"
)

func withTextResponse(next returnHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code, body := next(r)

		if err, ok := body.(error); ok || body == nil {
			resp := fmt.Sprintf("%d %s", code, http.StatusText(code))
			body = resp
			if err != nil {
				(r.Context().Value(contextKeyLogData)).(*logData).Error = err.Error()
				if Debug {
					resp += ": " + err.Error()
					body = resp
				}
			}
		}

		w.Header().Set(headerContentType, "text/plain")
		w.WriteHeader(code)

		if _, err := w.Write([]byte(body.(string))); err != nil {
			log.Println("Error writing text response:", err)
		}
	})
}
