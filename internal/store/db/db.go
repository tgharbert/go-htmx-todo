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
	return conn
}

func GetTodos(c *pgx.Conn) ([]Todo, error) {
		// get the todos...
		query := "SELECT * FROM todos ORDER BY id"
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

func AddTodo(c *pgx.Conn, title string) (Todo, error) {
	query := `INSERT INTO todos (date, done, title) VALUES ($1, $2, $3)`
	var todo Todo
	_, err := c.Exec(context.Background(), query, time.Now(), false, title)
	if err != nil {
		return todo, fmt.Errorf("failed to insert todo: %v", err)
	}

	lastQuery := `SELECT * FROM todos WHERE id = (SELECT MAX(id) FROM todos)`
		rows, err := c.Query(context.Background(), lastQuery)
		if err != nil {
			return todo, fmt.Errorf("initial query failed: %v", err)
		}
		defer rows.Close()
		for rows.Next(){
			err := rows.Scan(&todo.Id, &todo.Date, &todo.Done, &todo.Title)
			if err != nil {
				return todo, fmt.Errorf("error scanning row %v", err)
			}
		}
		if rows.Err() != nil {
			return todo, rows.Err()
		}
	return todo, nil
}

func ChangeTodo(c *pgx.Conn, id int64) (Todo, error) {
	var todo Todo
	query := `UPDATE todos SET done = NOT done WHERE id=$1 RETURNING *`
	rows := c.QueryRow(context.Background(), query, id)
	err := rows.Scan(&todo.Id, &todo.Date, &todo.Done, &todo.Title)
	if err != nil {
		if err == pgx.ErrNoRows {
			return todo, fmt.Errorf("no todo found with id: %d", id)
		}
		return todo, fmt.Errorf("updating todo failed: %v", err)
	}
	return todo, nil
}

func DeleteTodos(c *pgx.Conn) ([]Todo, error) {
	var todos []Todo
	query := `DELETE FROM todos WHERE done = true`
	_, err := c.Exec(context.Background(), query)
	if err != nil {
		return todos, fmt.Errorf("unable to delete todos: %v", err)
	}
	allQuery := "SELECT * FROM todos ORDER BY id"
	rows, err := c.Query(context.Background(), allQuery)
	if err != nil {
		return nil, fmt.Errorf("initial query failed: %v", err)
	}
	defer rows.Close()
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

func DeleteAllTodos(c *pgx.Conn) ([]Todo, error) {
	var todos []Todo
	query := "TRUNCATE TABLE todos"
	_, err := c.Exec(context.Background(), query)
	if err != nil {
		return todos, fmt.Errorf("error deleting all todos from db: %v", err)
	}
	return todos, nil
}