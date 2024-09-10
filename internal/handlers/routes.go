package routes

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

var todos = Todos{
	Todos: []Todo{
		{Id: 1, Date: time.Now(), Done: false, Title: "learn go"},
		{Id: 2, Date: time.Now(), Done: false, Title: "learn htmx"},
		{Id: 3, Date: time.Now(), Done: false, Title: "hook this thing up to MySQL"},
	},
}

func changeStatus (t *Todo) {
	t.Done = !t.Done
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	// template object renders the name of the template that we want to render
	// tmpl := template.Must(template.ParseFiles("../../templates/index.html"))
	// // pass the data to the template with the Execute method
	// tmpl.Execute(w, todos)

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, struct{ Todos []Todo }{Todos: todos.Todos})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AddTodo(w http.ResponseWriter, r *http.Request) {
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
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.ExecuteTemplate(w, "todos-list-element", newTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DeleteTodos(w http.ResponseWriter, r *http.Request) {
	// this is how to delete elements from an array in place
	j := 0
	for i := 0; i < len(todos.Todos); i++ {
		if (!todos.Todos[i].Done) {
			todos.Todos[j] = todos.Todos[i]
			j++
		}
	}
	todos.Todos = todos.Todos[:j]
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	// created a new template in the HTML file and now sending that template to replace the list
	err := tmpl.ExecuteTemplate(w, "todos-list", todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DeleteAllTodos(w http.ResponseWriter, r *http.Request) {
	todos.Todos = []Todo{}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.ExecuteTemplate(w, "todos-list", todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CheckTodo(w http.ResponseWriter, r *http.Request) {
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
			// added a func to remind of pointers and memory concern
			changeStatus(&todos.Todos[i])
			// changedTodo = todos.Todos[i]
		}
	}
	// add return of the ToDo then swap the HTMX stuff??
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.ExecuteTemplate(w, "todos-list-element", changedTodo)
}