package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// const fileName = "cache.json"
const cacheFileName = "cache/cache.json"

func LoadCache(cache map[string]Product) error {
	fileInfo, err := os.Stat(cacheFileName)
	if err != nil {
		fmt.Println("Error loading the file", cacheFileName)
		return err
	}

	if fileInfo.Size() == 0 {
		cache = make(map[string]Product)
		return nil
	}

	bytes, err := os.ReadFile(cacheFileName)
	if err != nil {
		fmt.Println("Error loading the file", cacheFileName)
		return err
	}

	return json.Unmarshal(bytes, &cache)
}

func PersistCacheWorker(writeCacheCh chan map[string]Product) {
	for {
		select {
		case c := <-writeCacheCh:
			PersistCache(c)
		}
	}
}

func PersistCache(c map[string]Product) {
	bytes, err := json.Marshal(c)
	if err != nil {
		fmt.Println("Error marshaling cache:", err)
		return
	}

	err = os.WriteFile(cacheFileName, bytes, 0644)
	if err != nil {
		fmt.Println("Error writing cache to disk:", err)
	}

	fmt.Println("Cache persisted")
}
