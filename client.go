package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func constructRequestBody(args []string) []byte {

	body, _ := json.Marshal(map[string]string{
		"name": args[1],
	})

	return body
}

func main() {

	command := os.Args[1:]

	client := &http.Client{
		Transport: nil,
		Jar:       nil,
		Timeout:   0,
	}

	method := ""

	if command[0] == "add" {
		method = "POST"
	} else if command[0] == "get" {
		method = "GET"
	}

	reqBody := bytes.NewReader(constructRequestBody(command))

	req, err := http.NewRequest(method, "http://localhost:6000", reqBody)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	}

	fmt.Println(string(respBody))
}
