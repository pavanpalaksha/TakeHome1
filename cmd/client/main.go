package main

import (
	walmart "Walmart"
	"flag"
)

func main() {
	server := flag.String("Chat server", "localhost", "ip address of the chat server")
	port := flag.Int("port", 13000, "Chat server port")
	userName := flag.String("username", "FNU", "username")
	flag.Parse()
	client := walmart.NewChatClient(*server, *port, *userName)
	client.Run()
}
