package cors

import (
	"net/http"
)

func SetHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, cache-control")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		h.ServeHTTP(w, r)
	})
}

// func NewCORS() {
//     mux.HandleFunc("/", handler)

//     // Настройка CORS с разрешением на доступ с любого источника
//     c := cors.New(cors.Options{
//         AllowedOrigins:   []string{"*"}, // Разрешает доступ с любого домена
//         AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"}, // Разрешенные HTTP методы
//         AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
//         AllowCredentials: true,
//     })

//     // Применение CORS к маршрутам
//    	return c.Handler(mux)
// }
 
