package main

import (
	"log"
	"net/http"
	routes "todo/internal/handlers"
	db "todo/internal/store/db"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	panic(err)
	// }
	// testVar := os.Getenv("DATABASE_URL");
	db.Connect()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", routes.Homepage)
	http.HandleFunc("/add-todo/", routes.AddTodo)
	http.HandleFunc("/delete-todos/", routes.DeleteTodos)
	http.HandleFunc("/check-todo/", routes.CheckTodo)
	http.HandleFunc("/delete-all/", routes.DeleteAllTodos)

	log.Print("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}