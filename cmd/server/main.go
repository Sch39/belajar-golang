// cmd/server/main.go
package main

import (
	"log"
	"net/http"

	"sch.dev/my-kasir-gw/internal/product"
	db "sch.dev/my-kasir-gw/internal/storage/sqlite"
)

func main() {
	databasePath := "data.db"
	database, err := db.NewSQLiteDB(databasePath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return
	}

	productRepo := product.NewRepository(database)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	mux := http.NewServeMux()

	productHandler.RegisterRoutes(mux)

	 mux.Handle("/", http.FileServer(http.Dir("./public")))
	 
	 server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	 }

	 log.Println("Server started on :8080")
	 err = server.ListenAndServe()
	 if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %s", err)
	 }
}