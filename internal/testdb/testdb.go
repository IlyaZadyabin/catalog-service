package testdb

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/IlyaZadyabin/catalog-service/internal/sqlscript"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var mu sync.Mutex

func DSN() string {
	if dsn := os.Getenv("TEST_DATABASE_URL"); dsn != "" {
		return dsn
	}
	return "postgres://postgres:password@localhost:5432/challenge?sslmode=disable"
}

func Open(t *testing.T) *gorm.DB {
	t.Helper()
	mu.Lock()
	t.Cleanup(mu.Unlock)

	dsn := DSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("postgres not available: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Skipf("postgres ping failed: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Skipf("postgres not available: %v", err)
	}

	if err := sqlscript.ExecDir(context.Background(), dsn, findSQLDir(t)); err != nil {
		t.Fatalf("apply sql: %v", err)
	}

	t.Cleanup(func() {
		_ = sqlDB.Close()
	})
	return db
}

func findSQLDir(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		candidate := filepath.Join(dir, "sql")
		if st, err := os.Stat(candidate); err == nil && st.IsDir() {
			return candidate
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("sql directory not found")
	return ""
}
