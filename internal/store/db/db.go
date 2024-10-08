package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Todos struct {
	Todos []Todo
}

type Todo struct {
	Id int
	Date time.Time
	Done bool
	Title string
}


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

	// get the todos...
	query := "SELECT * FROM todos"
	rows, err := conn.Query(context.Background(), query)
	// No longer failing here.
	if err != nil {
		log.Fatalf("Initial query failed: %v\n", err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next(){

		var todo Todo
		err := rows.Scan(&todo.Id, &todo.Date, &todo.Done, &todo.Title)
		if err != nil {
			log.Fatalf("Error scanning row %v\n", err)
		}
		todos = append(todos, todo)
	}

	fmt.Println(todos)
}
