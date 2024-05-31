package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Product struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type PythonResponse struct {
	WordCounts []WordCount `json:"wordCounts"`
	Message    string      `json:"message"`
}

type WordCount struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

func ScrapeProductDescription(decodedUrl string, c chan Product) {
	if decodedUrl == "" {
		c <- Product{"", ""}
		return
	}

	fmt.Println("decoded URL: " + decodedUrl)
	baseUrl := "http://localhost:8082/scrape"

	params := url.Values{}
	params.Add("url", decodedUrl)

	finalUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	resp, err := http.Get(finalUrl)
	if err != nil {
		fmt.Println("Error making request:", err)
		c <- Product{"", ""}
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-OK status code", resp.StatusCode)
		c <- Product{"", ""}
		return
	}

	var product Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		c <- Product{"", ""}
		return
	}

	c <- product
}

func FilterText(description string) map[string]int {
	productDescription := map[string]string{
		"productDescription": description,
	}

	resp, err := MyPost("http://127.0.0.1:8081/product", "application/json", productDescription)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-OK response status", resp.Status)
		return nil
	}

	var wordCount PythonResponse
	if err := DecodeJSONBody(resp, &wordCount); err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return nil
	}

	wordCounts := make(map[string]int)
	for _, wc := range wordCount.WordCounts {
		wordCounts[wc.Word] = wc.Count
	}

	return wordCounts
}
