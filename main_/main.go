package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", MainServer)
	http.ListenAndServe(":8000", nil)
}

func MainServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
