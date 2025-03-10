package cors

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

// func SetHeaders(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, cache-control")

// 		h.ServeHTTP(w, r)
// 	})
// }

// func New(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		c := cors.New(cors.Options{
// 			AllowedOrigins: []string{"*"},
// 			AllowedMethods: []string{http.MethodPost, http.MethodDelete, http.MethodGet, http.MethodPut},
// 			AllowCredentials: true,
// 			Debug: true,
// 		})

// 		h = c.Handler(h)
// 	})
// }

func New() func (h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {


		fn := func(w http.ResponseWriter, r *http.Request) {
			c := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{http.MethodPost, http.MethodDelete, http.MethodGet, http.MethodPut},
			AllowCredentials: true,
			Debug: false,
			})

			h = c.Handler(h)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			h.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}