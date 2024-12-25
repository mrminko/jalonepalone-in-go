package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/events", streamEvent)
	http.Handle("/", fileServer)
	http.ListenAndServe(":1234", nil)
}

func streamEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	tokens := []string{"hello", "this", "is", "a", "demo", "for", "event-stream"}
	for _, token := range tokens {
		content := fmt.Sprintf("data: %s\n\n", token)
		w.Write([]byte(content))
		w.(http.Flusher).Flush()
		time.Sleep(time.Millisecond * 1000)
	}
}
