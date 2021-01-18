## Simple Telnet like chat application
A simple Chat application where Chat clients can communicate with each other by connecting to server and sending message to all connected users.

## Components 
### Server
* Should accept new client request
* Should distinguish between `Client-A` and `Client-B`
* Upon receiving message from `Client-A`, it should send message to all connected clients except `Client-A`

### Client
* Can send message to all the users
* Can change display name
* Can gracefully Quit

## Building Application
1. make deps
2. go generate ./...
3. make

## Running Application
### Server
* To run server on a given port
```
./server -port=14000
```
* To run server on a given port with logs to stdout
```
./server -port=14000 -stderrthreshold=INFO -v=9
```

### Client
* To connect client to server
```
./client -server=localhost -port=14000 -username=user1
```
* To connect client to server with logs to stdout
```
./client -server=localhost -port=14000 -username=user1 -stderrthreshold=INFO -v=9
```

## Help
* `./client --help` and `./server --help` should display valid arguments and their usage
* At any time in Client application you can get help about commands by typing `HELP`
```
./client -username=user1 -server=localhost -port=14000
I0117 22:05:34.729454   25763 chatclient.go:28] Init

  Chat client connected to server...

  Type HELP to display help message

user1 # help

  Valid commands:
  HELP: print help message
  NICK <NEW_NAME>: change your username
  SEND <MESSAGE>: send message
  QUIT: disconnect from the server

```
