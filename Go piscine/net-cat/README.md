# Net-Cat

A simple TCP chat server and client written in Go, inspired by the Unix `netcat` utility. The application allows multiple users to join a shared chat room, exchange messages, and see when users connect or disconnect.

This project was developed as a team project by:

- (Me) Mohammad Mostafavi
- Silver Golovastenko
- Ghazal Varaminian

## Features

- TCP server with support for multiple concurrent clients
- Shared chat room
- Timestamped messages
- Join and leave notifications
- Configurable host and port
- Input and connection validation

## Requirements

- [Go](https://go.dev/) 1.20 or newer

## Installation

Clone the repository and move into the project directory:

```bash
git clone <repository-url>
cd net-cat
```

Build the application:

```bash
go build -o net-cat
```

## Usage

Start the server on the default port:

```bash
./net-cat
```

Or specify a custom port:

```bash
./net-cat 8080
```

To connect as a client, open another terminal and use NetCat:

```bash
nc localhost 8080
```

You can replace `localhost` with the server's IP address when connecting from another machine.

## Example

```text
$ ./net-cat 8080
Listening on port 8080...

$ nc localhost 8080
Welcome to TCP-Chat!
Please enter your name:
Alice
[2026-11-19 12:30:00][Alice]: Hello everyone!
```

## Project Structure

```text
.
├── README.md
├── go.mod
└── *.go
```

## Running Tests

If the project includes Go tests, run:

```bash
go test ./...
```

## Contributing

This project was created by a team; contributions should be discussed and tested before being submitted. Please keep changes consistent with the existing Go style and avoid modifying unrelated files.

## License

This project is licensed under the terms included in the repository's `LICENSE` file.
