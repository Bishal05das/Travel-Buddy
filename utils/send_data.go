package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendData(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		fmt.Println("Error encoding data:",err)
		return
	}
}