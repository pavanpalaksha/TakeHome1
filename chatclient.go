package walmart

import (
	pb "Walmart/protobuf"
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type ChatClient interface {
	Run()
}

type Client struct {
	clientConn    pb.ChatClient
	username      string
	uID           uint64
	msgStreamDone chan struct{}
}

func (cc *Client) init() error {
	glog.V(9).Info("Init")
	// Register
	userID, err := cc.clientConn.Register(context.Background(), &pb.User{Name: cc.username})
	if err != nil {
		glog.Errorf("Error during register - %v", err.Error())
		return err
	}
	cc.uID = userID.Id

	// Subscribe
	msgStream, err := cc.clientConn.Subscribe(context.Background(), userID)
	if err != nil {
		glog.Errorf("Error during Subscribe - %v", err.Error())
		return err
	}
	go func(msgs pb.Chat_SubscribeClient) {
		for {
			msg, err := msgs.Recv()
			if err != nil {
				glog.V(9).Infof("Graceful exit from server - %v", err.Error())
				break
			}
			fmt.Printf("\n%s  %s: %s\n", msg.Timestamp, msg.User.Name, msg.Msg)
			cc.printPrompt()
		}
		cc.msgStreamDone <- struct{}{}
	}(msgStream)
	return nil
}

func (cc *Client) parseCommand() (string, string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())
	cmds := strings.SplitN(input, " ", 2)
	glog.V(9).Info("parsedCmd, cmds - ", cmds)
	if len(cmds) > 1 {
		return cmds[0], cmds[1]
	}
	return cmds[0], ""
}

func (cc *Client) sendMsg(msg string) error {
	currTime := time.Now().Format("2006-01-02 15:04:05")
	pbMsg := &pb.Message{
		UID:       &pb.UserId{Id: cc.uID},
		User:      &pb.User{Name: cc.username},
		Msg:       msg,
		Timestamp: currTime}
	_, err := cc.clientConn.Send(context.Background(), pbMsg)
	if err != nil {
		glog.Errorf("Error during Send - %v", err.Error())
		return err
	}
	return nil
}

func (cc *Client) printPrompt() {
	fmt.Printf("%s # ", cc.username)
}

func (cc *Client) Run() {
	startMsg := `
  Chat client connected to server...

  Type HELP to display help message
  `
	helpMsg := `
  Valid commands:
  HELP: print help message
  NICK <NEW_NAME>: change your username
  SEND <MESSAGE>: send message
  QUIT: disconnect from the server
  `
	err := cc.init()
	if err != nil {
		return
	}
	fmt.Println(startMsg)
Loop:
	for true {
		cc.printPrompt()
		cmd, args := cc.parseCommand()
		switch strings.ToUpper(cmd) {
		case "HELP":
			fmt.Println(helpMsg)
		case "SEND":
			err := cc.sendMsg(args)
			if err != nil {
				break Loop
			}
		case "QUIT":
			_, err := cc.clientConn.Quit(context.Background(), &pb.UserId{Id: cc.uID})
			if err != nil {
				glog.Errorf("Error during Quit - %v", err.Error())
			}
			break Loop
		case "NICK":
			cc.username = args
		default:
			fmt.Printf("%s - Unknown command!!\n", cmd)
			fmt.Println(helpMsg)
		}
	}
	<-cc.msgStreamDone
	glog.V(9).Infof("Exiting client")
}

func NewChatClient(server string, portNum int, username string) ChatClient {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", server, portNum), grpc.WithInsecure())
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}
	return &Client{
		clientConn:    pb.NewChatClient(conn),
		username:      username,
		msgStreamDone: make(chan struct{})}
}
