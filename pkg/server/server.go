package server

import (
	"fmt"
	"net/http"
)

func StartServer(port string) {
	fmt.Println("Server started on :8080")
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
