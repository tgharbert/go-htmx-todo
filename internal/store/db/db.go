package db

import (
	"context"
	"fmt"
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


func Connect() *pgx.Conn {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	// defer conn.Close(context.Background())
	return conn
}

func GetTodos(c *pgx.Conn) ([]Todo, error) {
		// get the todos...
		query := "SELECT * FROM todos"
		rows, err := c.Query(context.Background(), query)
		if err != nil {
			return nil, fmt.Errorf("initial query failed: %v", err)
		}
		defer rows.Close()

		var todos []Todo
		for rows.Next(){
			var todo Todo
			err := rows.Scan(&todo.Id, &todo.Date, &todo.Done, &todo.Title)
			if err != nil {
				return nil, fmt.Errorf("error scanning row %v", err)
			}
			todos = append(todos, todo)
		}
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		return todos, nil
}

// MODIFY QUERY TO RETURN ADDED VALUE
func AddTodo(c *pgx.Conn, title string) (error) {
	query := `INSERT INTO todos (date, done, title) VALUES ($1, $2, $3)`
	// var todoId int
	_, err := c.Exec(context.Background(), query, time.Now(), false, title)

	if err != nil {
			return fmt.Errorf("failed to insert todo: %v", err)
	}
	return nil
}