package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	// POSTするエンドポイント

	url := "https://jsonplaceholder.typicode.com/posts"

	// 送信するJSONデータ
	data := map[string]string{
		"title": "foo",
		"body": "bar",
		"userId": "1",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println(result)
}