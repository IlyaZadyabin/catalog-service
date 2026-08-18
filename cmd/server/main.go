package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IlyaZadyabin/catalog-service/app/catalog"
	"github.com/IlyaZadyabin/catalog-service/app/categories"
	"github.com/IlyaZadyabin/catalog-service/app/database"
	"github.com/IlyaZadyabin/catalog-service/models"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, closeDB := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	defer func() {
		if err := closeDB(); err != nil {
			log.Printf("closing database: %v", err)
		}
	}()

	prodRepo := models.NewProductsRepository(db)
	catRepo := models.NewCategoryRepository(db)

	catalogSvc := catalog.NewService(prodRepo)
	categoriesSvc := categories.NewService(catRepo)

	catalogHandler := catalog.NewCatalogHandler(catalogSvc)
	productHandler := catalog.NewProductHandler(catalogSvc)
	categoriesHandler := categories.NewCategoriesHandler(categoriesSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", catalogHandler.HandleGet)
	mux.HandleFunc("GET /catalog/{code}", productHandler.HandleGetByCode)
	mux.HandleFunc("GET /categories", categoriesHandler.HandleList)
	mux.HandleFunc("POST /categories", categoriesHandler.HandleCreate)

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")),
		Handler: mux,
	}

	go func() {
		log.Printf("Starting server on http://%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}
		log.Println("Server stopped gracefully")
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %s", err)
	}
}
