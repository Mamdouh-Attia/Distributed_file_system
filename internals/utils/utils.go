// utils file
package utils

import (
	pb_d "Distributed_file_system/internals/pb/data_node"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
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

func ConvertStrIntoInt(str string) (int, error) {
	var num int
	_, err := fmt.Sscan(str, &num) // Handle potential parsing errors
	if err != nil {
		log.Printf("Failed to parse port number: %v", err)
		return -1, err
	}

	return num, nil
}

// FindMP4Files finds all the .mp4 files in the specified directory and its subdirectories
func FindMP4Files() ([]string, error) {
	var mp4Files []string

	// Open the current directory
	files, err := os.ReadDir("./")
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



func SaveFile(filename string, data []byte) error {

    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.Write(data)
    if err != nil {
        return err
    }

    fmt.Printf("File saved as %s\n", filename)
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


// TCP Listener connection
func ReceiveTCP(ip string, port string) (net.Conn, error) {
	
	//create a listener
	listener, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		fmt.Printf("error creating listener: %v\n", err)
		return nil, err
	}
	defer listener.Close()

	
	fmt.Print("Listening on " + ip + ":" + port + "\n")
	//accept the connection
	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("error accepting connection: %v\n", err)
		return nil, err
	}

	return conn, nil
}


func SendTCP(ip string, port string) (net.Conn, error) {
	conn, errConn := net.Dial("tcp", ip+":"+port)

	fmt.Print("Connected to " + ip + ":" + port + "\n")

	if errConn != nil {
		fmt.Printf("Failed to connect to server: %v", errConn)
		return nil, errConn
	}
	
	return conn, nil
}


// Serialize the request sent
func Serialize(request any, conn net.Conn) error {
		serializedRequest, errSerialize := json.Marshal(request)

	if errSerialize != nil {
		fmt.Printf("Failed to serialize the request: %v", errSerialize)
		return errSerialize
	}

	conn.Write(serializedRequest)


	return nil
}

// Deserialize the response received 
func Deserialize(data []byte, upload bool) error {
	
	file := &pb_d.UploadFileRequest{}
	if !upload {
		file = &pb_d.UploadFileRequest{}
	} 
	err := json.Unmarshal(data, file)
	if err != nil {
		fmt.Println("Error decoding file:", err)
		return err
	}
	
	print("FileName: ", file.FileName)

	// store the file content in the local file
	err = ioutil.WriteFile(file.FileName, file.FileContent, 0644)
	if err != nil {
		fmt.Printf("error saving the file: %v\n", err)
		return err
	}

	return nil
}