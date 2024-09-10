package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	routes "todo/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {

	// db.Connect()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	testVar := os.Getenv("TEST");
	fmt.Println(testVar)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))


	log.Print("Listening on port 8080...")
	http.HandleFunc("/", routes.Homepage)
	http.HandleFunc("/add-todo/", routes.AddTodo)
	http.HandleFunc("/delete-todos/", routes.DeleteTodos)
	http.HandleFunc("/check-todo/", routes.CheckTodo)
	http.HandleFunc("/delete-all/", routes.DeleteAllTodos)

	log.Fatal(http.ListenAndServe(":8080", nil))
}