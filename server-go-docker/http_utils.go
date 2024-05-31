package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func MyPost(url string, contentType string, data interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %w", err)
	}

	resp, err := http.Post(url, contentType, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error sending POST request: %w", err)
	}

	return resp, nil
}

func DecodeJSONBody(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("error decoding JSON response: %w", err)
	}
	return nil
}
