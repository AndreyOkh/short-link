package res

import (
	"encoding/json"
	"log"
	"net/http"
)

func Json(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("error json encode:", err)
	}
}
