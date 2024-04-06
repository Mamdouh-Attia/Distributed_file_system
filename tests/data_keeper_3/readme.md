# TCP for Data Keeper nodes as : Client-Server File Transfer

This is a simple TCP client-server file transfer task written in Go. The sender DK node listens for incoming connections from clients, receives a \<filename> request from the receiver DK node, and sends the corresponding file back to the receiver DK node.

## Features

- **Client (receiver DK node):** Sends a request to the server with the filename of the MP4 file it wants to receive.
- **Server(sender DK node):** Listens for incoming connections, receives the filename request from the client, sends the requested file back to the client.

## Usage

1. **Start the Server on a terminal:**

    ```
    go run server.go
    ```

    The server will start listening for incoming connections on `localhost:8080`.

2. **Run the Client on another terminal:**

    ```
    go run client.go
    ```

    The client will prompt you to enter the name of the MP4 file you want to receive. Enter the filename and press Enter.

3. **File Transfer:**

    - The client sends the filename to the server.
    - The server receives the filename request, opens the corresponding file, and sends its contents back to the client.
    - The client receives the file, writes it to disk as "\<filename>.mp4", and displays a success message.

## Requirements

- Go (Golang) installed on your machine.

## Note

- Make sure the MP4 file you want to transfer is located in the same directory as the server executable, or provide the correct path to the file in the client's input.
