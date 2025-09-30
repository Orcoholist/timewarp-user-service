// handlers/timewarp.go
package handlers

import (
	"net/http"
)

func TimewarpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Timewarp endpoint"))
	}
}
