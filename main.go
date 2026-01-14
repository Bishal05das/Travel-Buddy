package main

import (
	"fmt"
	"net/http"
)

func main(){
	mux := http.NewServeMux()
	fmt.Println("Listening to server on port 3000")
	err := http.ListenAndServe(":3000",mux)
	if err != nil {
		fmt.Println(fmt.Println("Server failed to start:", err))
	}
}