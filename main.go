package main

import (
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
	log.Print("Listening on port 8080...")
	// data that will be passed to the template - irl coming from DB
	todos := Todos{
		Todos: []Todo{
			{Id: 1, Date: time.Now(), Done: false, Title: "learn go"},
			{Id: 2, Date: time.Now(), Done: false, Title: "learn htmx"},
			{Id: 3, Date: time.Now(), Done: false, Title: "hook this thing up to MySQL"},
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
		// checking to see if there is a value sent
		if title == "" {
			return
		}
		// creating a new Todo Obj
		newTodo := makeNewTodo(title)
		// appending the new todo to the current todos
		todos.Todos = append(todos.Todos, newTodo)
		// creating the template to replace the value
		tmpl := template.Must(template.ParseFiles("index.html"))
		err := tmpl.ExecuteTemplate(w, "todos-list-element", newTodo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	deleteTodos := func(w http.ResponseWriter, r *http.Request) {
		// this is how to delete elements from an array in place
		j := 0
		for i := 0; i < len(todos.Todos); i++ {
			if (!todos.Todos[i].Done) {
				todos.Todos[j] = todos.Todos[i]
				j++
			}
		}
		todos.Todos = todos.Todos[:j]
		tmpl := template.Must(template.ParseFiles("index.html"))
		// created a new template in the HTML file and now sending that template to replace the list
		err := tmpl.ExecuteTemplate(w, "todos-list", todos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	checkTodo := func(w http.ResponseWriter, r *http.Request) {
		// getting the Id from the request URL
		newId := r.URL.Path[len("/check-todo/"):]
		// converting that string into a number
		id, err := strconv.ParseInt(newId, 10, 8)
		if err != nil {
			log.Print("Error altering done status of item")
		}
		var changedTodo Todo
		for i := 0; i < len(todos.Todos); i++ {
			if todos.Todos[i].Id == int(id) {
				todos.Todos[i].Done = !todos.Todos[i].Done
				changedTodo = todos.Todos[i]
			}
		}
		// add return of the ToDo then swap the HTMX stuff??
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "todos-list-element", changedTodo)
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/add-todo/", addTodo)
	http.HandleFunc("/delete-todos/", deleteTodos)
	http.HandleFunc("/check-todo/", checkTodo)

	log.Fatal(http.ListenAndServe(":8080", nil))
}