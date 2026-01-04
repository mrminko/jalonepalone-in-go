package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := os.Args[1]
	r, err := http.Get(addr)
	if err != nil {
		log.Fatalln(err)
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Body.Close()
	fmt.Printf("%s\n", b)
}
