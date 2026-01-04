package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/goccy/go-yaml"
	"io"
	"net/http"
	"os"
	"sync"
)

type Api struct {
	Label    string
	Name     string
	Endpoint string
	Method   string
}

type Config struct {
	Description string
	Apis        []Api
}

type ApiResponse struct {
	label  string
	result string
	err    error
}

type ApiResponseChannel chan *ApiResponse

func main() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	config := &Config{}
	if err = yaml.Unmarshal(data, config); err != nil {
		panic(err)
	}
	run(config.Apis)
}

func writeResultToBuffer(buf *bytes.Buffer, ch ApiResponseChannel, done chan<- struct{}) {
loop:
	for {
		select {
		case s, ok := <-ch:
			if !ok {
				break loop
			}
			if s.err != nil {
				fmt.Println("error: ", s.err)
				continue
			}
			buf.WriteString(fmt.Sprintf("%s: %s\n", s.label, s.result))
		}
	}
	done <- struct{}{}
}

func run(apis []Api) {
	buffer := &bytes.Buffer{}
	wg := &sync.WaitGroup{}
	ch := make(ApiResponseChannel, len(apis))
	done := make(chan struct{})
	go writeResultToBuffer(buffer, ch, done)
	for _, api := range apis {
		wg.Add(1)
		go fetchApi(&api, ch, wg)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	<-done
	fmt.Println("result: ", buffer.String())
	os.Stdout.Write(buffer.Bytes())
}

func fetchApi(api *Api, ch ApiResponseChannel, wg *sync.WaitGroup) {
	defer wg.Done()
	switch api.Method {
	case "GET", "get":
		resp, err := http.Get(api.Endpoint)
		if err != nil {
			ch <- &ApiResponse{
				label: api.Label,
				err:   err,
			}
			return
		}
		defer resp.Body.Close()
		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			ch <- &ApiResponse{
				label: api.Label,
				err:   err,
			}
			return
		}
		buf := &bytes.Buffer{}
		err = json.Indent(buf, raw, "", "  ")
		if err != nil {
			ch <- &ApiResponse{
				label: api.Label,
				err:   err,
			}
			return
		}
		ch <- &ApiResponse{
			label:  api.Label,
			result: buf.String(),
			err:    nil,
		}
	}
}
