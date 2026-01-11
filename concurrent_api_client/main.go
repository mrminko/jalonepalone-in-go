package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/goccy/go-yaml"
)

type Api struct {
	Label    string
	Name     string
	Endpoint string
	Method   string
	Payload  map[string]interface{}
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
		get(api, ch)
	case "POST", "post":
		post(api, ch)

	}
}

func post(api *Api, ch ApiResponseChannel) {
	jsonData, err := json.Marshal(api.Payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	resp, err := http.Post(api.Endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		ch <- &ApiResponse{
			label: api.Label,
			err:   err,
		}
		return
	}
	handleRespBody(api, resp, ch)
}

func get(api *Api, ch ApiResponseChannel) {
	resp, err := http.Get(api.Endpoint)
	if err != nil {
		ch <- &ApiResponse{
			label: api.Label,
			err:   err,
		}
		return
	}
	handleRespBody(api, resp, ch)
}

func handleRespBody(api *Api, resp *http.Response, ch ApiResponseChannel) {
	defer resp.Body.Close()
	var body []byte
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- &ApiResponse{
			label: api.Label,
			err:   err,
		}
		return
	}
	if resp.Header.Get("content-type") != "application/json" {
		buf := &bytes.Buffer{}
		err = json.Indent(buf, body, "", "  ")
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
	} else {
		ch <- &ApiResponse{
			label:  api.Label,
			result: string(body),
			err:    nil,
		}
	}
}
