package main

import (
	"fmt"

	"github.com/mrminko/jalonepalone-in-go/port-scanner/port"
)

func main() {
	open := port.ScanPort("tcp", "localhost", 8000)
	fmt.Printf("Port open: %t\n", open)
}
