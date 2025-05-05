package cors

import (
	"net/http"

	"github.com/rs/cors"
)

func NewCORS(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:5173/"},
		AllowedMethods:   []string{http.MethodPost, http.MethodDelete, http.MethodGet, http.MethodPut, http.MethodPatch},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization", "cache-control"},
		AllowCredentials: true,
		MaxAge: 120,			
	})

	h = c.Handler(h)
		
	return h
}