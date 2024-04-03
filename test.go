package main

import (
	"fmt"
	"os"
)

func main() {

	migrationsPath := fmt.Sprintf("file:%s", getEnv("MIGRATION_PATH", "E:/cmd/migrate/migrations"))
	fmt.Println(migrationsPath)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
