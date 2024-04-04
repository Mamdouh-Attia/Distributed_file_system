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

	// Open the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Iterate through the files in the directory
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".mp4") {
			// Extract the filename without the directory path
			filename := filepath.Base(file.Name())
			// Append the filename to the mp4Files slice
			mp4Files = append(mp4Files, filename)
		}
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
func OpenFileFromDirectory(dir string, filename string) (*os.File, error) {
	// Change directory to dir
	errDir := os.Chdir(dir)
	if errDir != nil {
		fmt.Println("Error changing directory:", errDir)
		return nil, errDir
	}

	// Open the file
	file, err := os.Open(filename)

	return file, err
}
