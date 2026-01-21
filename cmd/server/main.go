package main

import (
	"log"
	"net/http"
)

func main() {
	 http.Handle("/", http.FileServer(http.Dir("./public")))
	 
	 log.Println("Server started on :8080")
	 err := http.ListenAndServe(":8080", nil)
	 if err != nil {
		log.Println("Error starting server:", err)
	 }
}