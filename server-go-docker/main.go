package main

import (
	"fmt"
	"os"

	"net/http"
	"net/url"

	"sync"

	"github.com/gin-gonic/gin"

	"os/exec"
	"runtime"
)

var (
	urls           []string
	cacheData      = make(map[string]Product)
	lock           sync.RWMutex
	writeCacheCh   = make(chan map[string]Product, 100)
	productFetchCh = make(chan Product, 100)
)

func main() {
	if err := LoadCache(cacheData); err != nil {
		fmt.Println("Error loading cache:", err)
		os.Exit(1)
	}

	go PersistCacheWorker(writeCacheCh)

	router := gin.Default()
	router.POST("/", wordCloudRequest)
	router.Run("localhost:8080")
}

func openImage(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	case "darwin":
		cmd = exec.Command("open", filePath)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", filePath)
	default:
		return fmt.Errorf("unsupported platform")
	}

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	err = cmd.Process.Release()
	if err != nil {
		return fmt.Errorf("failed to release command: %w", err)
	}

	return nil
}

func wordCloudRequest(c *gin.Context) {
	newUrl := c.Query("productUrl")

	if newUrl == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter"})
		return
	}

	decodedUrl, err := url.QueryUnescape(newUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL parameter format"})
		return
	}

	urls = append(urls, decodedUrl)

	lock.RLock()
	product, found := cacheData[decodedUrl]
	lock.RUnlock()

	if !found {
		go ScrapeProductDescription(decodedUrl, productFetchCh)
		product = <-productFetchCh

		lock.Lock()
		cacheData[decodedUrl] = product
		lock.Unlock()
	} else {
		fmt.Println("Using cache value")
	}

	writeCacheCh <- cacheData

	wordCounts := FilterText(product.Description)

	if wordCounts == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process text"})
		return
	}

	if len(wordCounts) == 0 {
		fmt.Printf("Product %s does not have product description\n", product.Title)
		return
	}

	imagePath := CreateWordCloud(product, wordCounts)
	err = openImage(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening image: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Image opened successfully")

	c.JSON(http.StatusCreated, decodedUrl)
}
