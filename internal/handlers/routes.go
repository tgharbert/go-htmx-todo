package routes

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
	"todo/internal/store/db"
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

var id = 3
func makeNewTodo (title string) Todo {
	id++
	return Todo {
		Id: id,
		Date: time.Now(),
		Done: false,
		Title: title,
	}
}

// var todos = Todos{
// 	Todos: []Todo{
// 		{Id: 1, Date: time.Now(), Done: false, Title: "learn go"},
// 		{Id: 2, Date: time.Now(), Done: false, Title: "learn htmx"},
// 		{Id: 3, Date: time.Now(), Done: false, Title: "hook this thing up to MySQL"},
// 	},
// }

// func changeStatus (t *Todo) {
// 	t.Done = !t.Done
// }

type PageData struct {
	Todos []db.Todo
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	// template object renders the name of the template that we want to render
	tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/footer.html", "templates/todolist.html"))
	conn := db.Connect()
	defer conn.Close(context.Background())
	todos, err := db.GetTodos(conn)
	if err != nil {
		log.Fatalf("error getting todos: %v", err)
	}
	// fmt.Println("todos: ", todos)
	data := PageData{
		Todos: todos,
	}
	// pass the data to the template with the Execute method
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AddTodo(w http.ResponseWriter, r *http.Request) {
	// getting the value from the form
	title := r.PostFormValue("todo")
	// checking to see if there is a value sent
	if title == "" {
		fmt.Println("test")
		return
	}
	conn := db.Connect()
	defer conn.Close(context.Background())
	todo, err := db.AddTodo(conn, title)
	if err != nil {
		log.Fatalf("error adding todo: %v", err)
	}

	// checking for bad value, returning a todo no matter what
	if todo.Title == "" {
		return
	}
	// creating the template to replace the value
	tmpl := template.Must(template.ParseFiles("templates/todolist.html"))
	err = tmpl.ExecuteTemplate(w, "todos-list-element", todo)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func DeleteTodos(w http.ResponseWriter, r *http.Request) {
// 	// this is how to delete elements from an array in place
// 	// j := 0
// 	// for i := 0; i < len(todos.Todos); i++ {
// 	// 	if (!todos.Todos[i].Done) {
// 	// 		todos.Todos[j] = todos.Todos[i]
// 	// 		j++
// 	// 	}
// 	// }
// 	// todos.Todos = todos.Todos[:j]

// 	// will delete value from table, then return the entire todos list

// 	db.DeleteTodo()
// 	tmpl := template.Must(template.ParseFiles("templates/todolist.html"))
// 	// created a new template in the HTML file and now sending that template to replace the list
// 	err := tmpl.ExecuteTemplate(w, "todos-list", todos)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// func DeleteAllTodos(w http.ResponseWriter, r *http.Request) {
// 	todos.Todos = []Todo{}
// 	tmpl := template.Must(template.ParseFiles("templates/todolist.html"))
// 	err := tmpl.ExecuteTemplate(w, "todos-list", todos)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

func CheckTodo(w http.ResponseWriter, r *http.Request) {
	// getting the Id from the request URL
	newId := r.URL.Path[len("/check-todo/"):]

	// converting that string into a number
	id, err := strconv.ParseInt(newId, 10, 8)
	if err != nil {
		log.Print("Error altering done status of item")
	}

	conn := db.Connect()
	defer conn.Close(context.Background())

	// var todo Todo
	todo, err := db.ChangeTodo(conn, id, true)
	if err != nil {
		log.Fatalf("error checking todo: %v", err)
	}

	if todo.Title == "" {
		return
	}
	fmt.Println(todo)
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.ExecuteTemplate(w, "todos-list-element", todo)
}