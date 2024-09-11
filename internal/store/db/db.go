package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func Connect() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
	}
}