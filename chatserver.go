package walmart

import (
	pb "Walmart/protobuf"
	"context"
	"errors"
	"sync"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClientConn struct {
	msgStream pb.Chat_SubscribeServer
	name      string
	// Need this to send error to Subscribe method of ChatServer interface
	error chan error
}

type ChatServer struct {
	*pb.UnimplementedChatServer
	// Keeping it simple here.
	//
	// When a user registers value of nextUserId is assgined to the user
	//
	nextUserId uint64
	// mutex for allocating userId
	mutex sync.RWMutex
	// map of userId to ClientConn
	idToClientConn map[uint64]*ClientConn
}

func (cs *ChatServer) getNextUID() (UID uint64, err error) {
	// In the interest of time, I am using an incrementing counter for UID.
	if _, ok := cs.idToClientConn[cs.nextUserId]; ok {
		err = errors.New("Available User ID's are taken")
	} else {
		UID, err = cs.nextUserId, nil
		cs.nextUserId++
	}
	return
}

// Register RPC
//
func (cs *ChatServer) Register(ctx context.Context, user *pb.User) (*pb.UserId, error) {
	glog.V(9).Infof("Register: User - %+v", user)
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	UID, err := cs.getNextUID()
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	glog.V(9).Infof(" User - %+v, Id = %v", user, UID)
	cs.idToClientConn[UID] = &ClientConn{name: user.Name}
	return &pb.UserId{Id: UID}, nil

}

// Subscribe RPC
func (cs *ChatServer) Subscribe(uID *pb.UserId, msgStream pb.Chat_SubscribeServer) error {
	glog.V(9).Infof("Register: UserId - %+v", uID)
	userConn, ok := cs.idToClientConn[uID.Id]
	if !ok {
		return status.Errorf(codes.FailedPrecondition, "Unregistered user - %v", uID.Id)
	}
	cs.mutex.Lock()
	userConn.msgStream = msgStream
	userConn.error = make(chan error)
	cs.mutex.Unlock()
	err := <-userConn.error
	// Client exited either gracefully or non-gracefully
	glog.V(9).Infof("Register: Closing Conn - %+v", cs.idToClientConn[uID.Id])
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	delete(cs.idToClientConn, uID.Id)
	return err
}

// Send RPC
func (cs *ChatServer) Send(ctx context.Context, msg *pb.Message) (*pb.NoResponse, error) {
	glog.V(9).Infof("Send: Message - %+v", msg)
	// Read lock should be sufficient here
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	wait := sync.WaitGroup{}
	for uID, userConn := range cs.idToClientConn {
		if uID == msg.UID.Id || userConn.msgStream == nil {
			// Skip sending message to the sender and to users who haven't subscribed yet
			continue
		}
		wait.Add(1)
		go func(msg *pb.Message, userConn *ClientConn) {
			defer wait.Done()
			glog.V(9).Infof("Sending message: %v to %v", msg.Msg, msg.User)
			err := userConn.msgStream.Send(msg)
			if err != nil {
				glog.Errorf("Error sending message to - %+v, err = %+v", msg.User, err)
				userConn.error <- err
			}
		}(msg, userConn)
	}
	wait.Wait()
	return &pb.NoResponse{}, nil
}

// Quit RPC
func (cs *ChatServer) Quit(ctx context.Context, uID *pb.UserId) (*pb.NoResponse, error) {
	if userConn, ok := cs.idToClientConn[uID.Id]; ok {
		glog.V(9).Infof("Quit: User - %+v", userConn.name)
		userConn.error <- nil
	}
	return &pb.NoResponse{}, nil
}

func registerServer(server *grpc.Server) {
	cs := &ChatServer{
		nextUserId:     1,
		mutex:          sync.RWMutex{},
		idToClientConn: make(map[uint64]*ClientConn)}
	pb.RegisterChatServer(server, cs)
}
