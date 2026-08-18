package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/IlyaZadyabin/catalog-service/internal/sqlscript"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	dir := os.Getenv("POSTGRES_SQL_DIR")
	if err := sqlscript.ExecDir(context.Background(), dsn, dir); err != nil {
		log.Fatalf("seed failed: %v", err)
	}
	log.Println("Seed completed")
}
