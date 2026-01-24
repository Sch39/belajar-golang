// cmd/server/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"

	"sch.dev/my-kasir-gw/internal/category"
	"sch.dev/my-kasir-gw/internal/product"
	"sch.dev/my-kasir-gw/internal/storage/database"
	"sch.dev/my-kasir-gw/internal/storage/repository/sqlite"

	_ "sch.dev/my-kasir-gw/docs"
)

// @title My Kasir API Guweh
// @version 1.0
// @description Ini adalah server API untuk aplikasi My Kasir Guweh.
// @host localhost:8080
// @BasePath /
func main() {
	// Database
	databasePath := "data.db"
	conn, err := database.NewSqliteConn(databasePath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return
	}

	// Migration
	if err := database.MigrateSqlite(conn); err != nil {
		log.Fatal("Database migration failed:", err)
		return
	}

	// Product
	productRepo := sqlite.NewProductRepository(conn)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	// Category
	categoryRepo := sqlite.NewCategoryRepository(conn)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)

	// Mux
	mux := http.NewServeMux()

	// Register route
	productHandler.RegisterRoutes(mux)
	categoryHandler.RegisterRoutes(mux)

	// Front End
	mux.Handle("/", http.FileServer(http.Dir("./public")))

	// Health
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":    "UP",
			"version":   "1.0.0",
			"timestamp": time.Now().Format(time.RFC3339),
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
