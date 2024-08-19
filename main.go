package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
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

func main() {
	fmt.Println("Listening on port 8080...")
	// data that will be passed to the template - irl coming from DB
	todos := Todos{
		Todos: []Todo{
			{Id: 1, Date: time.Now(), Done: false, Title: "learn go"},
			{Id: 2, Date: time.Now(), Done: false, Title: "learn htmx"},
			{Id: 3, Date: time.Now(), Done: false, Title: "don't be dumb"},
		},
	}

	home := func(w http.ResponseWriter, r *http.Request) {
		// template object renders the name of the template that we want to render
		tmpl := template.Must(template.ParseFiles("index.html"))
		// pass the data to the template with the Execute method
		tmpl.Execute(w, todos)
	}

	addTodo := func(w http.ResponseWriter, r *http.Request) {
		// getting the value from the form
		title := r.PostFormValue("todo")
		// creating a new Todo Obj
		newTodo := makeNewTodo(title)
		// appending the new todo to the current todos
		todos.Todos = append(todos.Todos, newTodo)
		// creating the template to replace the values
		tmpl := template.Must(template.ParseFiles("index.html"))
		err := tmpl.ExecuteTemplate(w, "todos-list-element", newTodo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	deleteTodos := func(w http.ResponseWriter, r *http.Request) {
		log.Print("hit the delete")
	}

	checkTodo := func(w http.ResponseWriter, r *http.Request) {
		// getting the Id from the request URL
		uri := r.RequestURI
		// getting the last digit from the URI, this wouldn't work with two digit numbers....
		stringId := string(uri[len(uri) - 1])
		// converting that string into a number
		id, err := strconv.ParseInt(stringId, 10, 32)
		if err != nil {
			log.Print("Error altering done status of item")
		}
		// toggling the done status of the item
		todos.Todos[id - 1].Done = !todos.Todos[id - 1].Done
		// fmt.Println(todos.Todos[id - 1].Done)
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/add-todo/", addTodo)
	http.HandleFunc("/delete-todos/", deleteTodos)
	http.HandleFunc("/check-todo/", checkTodo)

	log.Fatal(http.ListenAndServe(":8080", nil))
}