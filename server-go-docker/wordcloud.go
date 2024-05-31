package main

import (
	"fmt"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"

	"github.com/psykhi/wordclouds"
)

const fontsPath = "fonts/roboto/Roboto-Black.ttf"

func CreateWordCloud(product Product, mostImportantWords map[string]int) string {
	fmt.Println("Creating word cloud for", product.Title)

	scaleFactor := 5
	scaledWordCounts := make(map[string]int)
	for word, count := range mostImportantWords {
		scaledWordCounts[word] = count * scaleFactor
	}

	var DefaultColors = []color.Color{
		color.RGBA{27, 27, 27, 255},    // Dark gray
		color.RGBA{72, 72, 75, 255},    // Gray
		color.RGBA{89, 58, 238, 255},   // Blue
		color.RGBA{101, 205, 250, 255}, // Light blue
		color.RGBA{112, 214, 191, 255}, // Green
	}

	rand.Seed(time.Now().UnixNano())
	randomColors := make([]color.Color, 0, len(scaledWordCounts))
	for i := 0; i < len(scaledWordCounts); i++ {
		randomColor := DefaultColors[rand.Intn(len(DefaultColors))]
		randomColors = append(randomColors, randomColor)
	}

	w := wordclouds.NewWordcloud(
		scaledWordCounts,
		wordclouds.FontFile(fontsPath), // Ensure this file exists
		wordclouds.Height(2048),
		wordclouds.Width(2048),
		wordclouds.Colors(randomColors),
		wordclouds.BackgroundColor(color.RGBA{255, 255, 255, 255}), // White background
	)

	img := w.Draw()
	fileName := fmt.Sprintf("output/%s.png", product.Title)

	outputFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return ""
	}
	defer outputFile.Close()

	if err := png.Encode(outputFile, img); err != nil {
		fmt.Println("Error encoding PNG:", err)
	}

	return fileName
}
