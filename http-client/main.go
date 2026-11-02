package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)
const DEFAULT_CLIENT_FILE = "clientfile"

var CACHE = map[interface{}]interface{}{}

type HTTPMethod string
const (
	GET HTTPMethod = "GET"
	POST HTTPMethod = "POST"
	PUT HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
)
type request struct {
	env    *os.File
	method HTTPMethod
	URL string
}
type clientFile struct {
	requests []*request
}

func main() {
	files := os.Args[1:]
	for _, f := range files {
		parseFile(f)
	}
}

func parseFile(filename string) *clientFile {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}
	defer file.Close()
	cf := &clientFile{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()
		if t[:3] == "###" {
			continue //comment line
		}
		if t == "\n" {
			cf.requests = append(cf.requests, buildRequest(scanner))
		}
	}
	return cf
}

//called at the start of every new request config
func buildRequest(scanner * bufio.Scanner) *request {
	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), " ")
		if !validateMethod(tokens[0]) {
			panic("Invalid method: " + tokens[0])
		}

	}
}

func validateMethod(method string) bool {
	if method == "GET" || method == "POST" || method == "PUT" || method == "DELETE" {
		return true
	}
	return false
}