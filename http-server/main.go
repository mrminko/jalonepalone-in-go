package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Todo string

type todoDatabase map[Todo]bool

func (d todoDatabase) Create(w http.ResponseWriter, r *http.Request) {
	newTodo := Todo(r.URL.Query().Get("name"))
	if _, ok := db[newTodo]; !ok {
		d[newTodo] = false
		fmt.Fprintln(w, "created")
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "item already exist")
}

func (d todoDatabase) Update(w http.ResponseWriter, r *http.Request) {
	todo := Todo(r.URL.Query().Get("name"))
	value := r.URL.Query().Get("value")
	var newValue bool
	switch value {
	case "True", "true":
		newValue = true
	case "False", "false":
		newValue = false
	}
	if _, ok := db[todo]; ok {
		d[todo] = newValue
		fmt.Fprintln(w, "updated")
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "item not found")
}

func (d todoDatabase) Delete(w http.ResponseWriter, r *http.Request) {
	todo := Todo(r.URL.Query().Get("name"))
	if _, ok := db[todo]; ok {
		delete(d, todo)
		fmt.Fprintln(w, "deleted")
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "item not found")
}

func (d todoDatabase) Get(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	for k, v := range d {
		buf.WriteString(fmt.Sprintf("%s: %v", k, v))
	}
	fmt.Fprintln(w, buf.String())
}

func (d todoDatabase) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		d.Get(w, r)
	case "POST":
		d.Create(w, r)
	case "DELETE":
		d.Delete(w, r)
	case "PUT":
		d.Update(w, r)
	}
}

var db = todoDatabase{}

func main() {
	port := os.Args[1]
	http.Handle("/todo", db)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf("localhost:%s", port), nil))
}
