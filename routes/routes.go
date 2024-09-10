package routes

import (
	"html/template"
	"net/http"
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

var todos = Todos{
	Todos: []Todo{
		{Id: 1, Date: time.Now(), Done: false, Title: "learn go"},
		{Id: 2, Date: time.Now(), Done: false, Title: "learn htmx"},
		{Id: 3, Date: time.Now(), Done: false, Title: "hook this thing up to MySQL"},
	},
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	// template object renders the name of the template that we want to render
	tmpl := template.Must(template.ParseFiles("index.html"))
	// pass the data to the template with the Execute method
	tmpl.Execute(w, todos)
}

