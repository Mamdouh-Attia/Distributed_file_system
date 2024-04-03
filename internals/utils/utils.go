// utils file
package utils

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GenerateID() int {
	// Generate a random number
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(1000)

	// Get the current timestamp in nanoseconds
	currentTime := time.Now().UnixNano()

	// Combine the timestamp and the random number to form the ID
	id := int(currentTime)%1000*1000 + randomNumber

	return id
}

// FindMP4Files finds all the .mp4 files in the specified directory and its subdirectories
func FindMP4Files(dir string) ([]string, error) {
	var mp4Files []string

	// Walk through the directory and its subdirectories
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if the file has a .mp4 extension
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".mp4") {
			mp4Files = append(mp4Files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return mp4Files, nil
}



func SaveFile(directory, filename string, data []byte) error {
    err := os.MkdirAll(directory, 0755) // Create the directory if it doesn't exist
    if err != nil {
        return err
    }

    fullPath := filepath.Join(directory, filename)
    file, err := os.Create(fullPath)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.Write(data)
    if err != nil {
        return err
    }

    fmt.Printf("File saved as %s\n", fullPath)
    return nil
}
