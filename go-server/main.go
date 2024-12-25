package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Form Parsing Error: %v", err), http.StatusInternalServerError)
		return
	}
	name := r.FormValue("name")
	address := r.FormValue("address")
	fmt.Fprintf(w, "Name: %s\n", name)
	fmt.Fprintf(w, "Address: %s\n", address)
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	if url := r.URL.Path; url != "/hello" {
		http.Error(w, fmt.Sprintf("Bad Request: %v", url), http.StatusBadRequest)
		return
	}
	if method := r.Method; method != "GET" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	fmt.Fprintf(w, "Hello bro!")
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", handleForm)
	http.HandleFunc("/hello", handleHello)
	fmt.Println("Serving on port 1234")
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatalln("Error when stating server:", err)
	}
}
