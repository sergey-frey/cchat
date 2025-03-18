package cors

import (
	"net/http"
)

func NewCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		h.ServeHTTP(w, r)
	})
}

// func NewCORS(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		c := cors.New(cors.Options{
// 			AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:8040"},
// 			AllowedMethods:   []string{http.MethodPost, http.MethodDelete, http.MethodGet, http.MethodPut},
// 			AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization", "Cache-Control"},
// 			AllowCredentials: true,
// 			Debug:            false,			
// 		})

// 		h = c.Handler(h)

// 		h.ServeHTTP(w, r)

// 	})
// }