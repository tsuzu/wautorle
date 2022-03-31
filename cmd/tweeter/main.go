package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

var client = http.Client{
	Timeout: 10 * time.Second,
}

func tweet(result, imageURL string) error {
	type Payload struct {
		Value1 string `json:"value1"`
		Value2 string `json:"value2"`
	}

	data := Payload{
		Value1: result,
		Value2: imageURL,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", os.Getenv("TWEET_IFTTT_IMG_URL"), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func main() {
	payload, err := os.ReadFile(os.Args[1])

	if err != nil {
		panic(err)
	}

	if err := tweet(string(payload), os.Args[2]); err != nil {
		panic(err)
	}
}
