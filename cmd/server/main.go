// cmd/server/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"sch.dev/my-kasir-gw/internal/product"
	db "sch.dev/my-kasir-gw/internal/storage/sqlite"

	_ "sch.dev/my-kasir-gw/docs"
)

// @title My Kasir API Guweh
// @version 1.0
// @description Ini adalah server API untuk aplikasi My Kasir Guweh.
// @host localhost:8080
// @BasePath /
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
	 
	 mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	 }))

	 mux.Handle("/swagger/", httpSwagger.WrapHandler)
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