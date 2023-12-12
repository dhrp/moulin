package server

// Here we define all GRPC handlers

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	pb "github.com/dhrp/moulin/pkg/protobuf"
	"github.com/dhrp/moulin/pkg/rouge"
)

type server struct {
	rouge *rouge.Client
	pb.UnimplementedAPIServer
}

// NewGRPCServer returns a new GRPC server
func NewGRPCServer(rougeClient *rouge.Client) *grpc.Server {

	// Initialize the API Server
	s := &server{rouge: rougeClient}

	// Setup grpc server
	grpcServer := grpc.NewServer()

	// Register reflection service on gRPC server.
	pb.RegisterAPIServer(grpcServer, s)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	return grpcServer
}

func (s *server) Healthz(ctx context.Context, in *empty.Empty) (*pb.StatusMessage, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println("read incoming context successful:")
		fmt.Println(md)
		fmt.Println(md["authentication"])
	} else {
		fmt.Println("failed to load context (meta)")
	}

	result, _ := s.rouge.Info()

	return &pb.StatusMessage{
		Status: pb.Status_SUCCESS,
		Detail: result,
	}, nil
}

func (s *server) PushTask(ctx context.Context, in *pb.Task) (*pb.StatusMessage, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		//ToDo: on http we get context, but authorization may not be set
		var authorization string
		// var ok bool
		if _, ok = md["authorization"]; ok {
			authorization = md["authorization"][0]
		}

		fmt.Println("authorization: " + authorization)
	} else {
		fmt.Println("failed to load context (meta)")
	}

	queueID := in.QueueID
	taskID := ksuid.New().String()

	taskMessage := rouge.TaskMessage{ID: taskID, Body: in.Body}
	queueLength, err := s.rouge.AddTask(queueID, taskMessage)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to push task")
	}
	result := fmt.Sprintf("queue %v now %d items long", queueID, queueLength)

	return &pb.StatusMessage{
		Status: pb.Status_SUCCESS,
		Detail: result,
	}, nil
}

// LoadTask returns a task from the redis queue
func (s *server) LoadTask(ctx context.Context, in *pb.RequestMessage) (*pb.Task, error) {

	queueID := in.QueueID
	taskMessage, err := s.rouge.Load(ctx, queueID, 300)
	if err != nil {
		err = errors.Wrap(err, "[grpc.go] error in loading message")
		log.Println(err)
		return &pb.Task{}, err
	}
	return &pb.Task{TaskID: taskMessage.ID, Body: taskMessage.Body}, nil

}

// LoadTask returns a task from the redis queue
func (s *server) HeartBeat(ctx context.Context, in *pb.Task) (*pb.StatusMessage, error) {

	queueID := in.QueueID
	taskID := in.TaskID

	var expirationSec int32
	if in.ExpirationSec != 0 {
		expirationSec = in.ExpirationSec
	} else {
		expirationSec = 300 // default expiration is 5 min
	}

	expires, err := s.rouge.Heartbeat(queueID, taskID, expirationSec)
	if err != nil {
		// ToDo return error to client, it could be because item was just updated
		log.Println(err)
		return &pb.StatusMessage{
			Status: pb.Status_FAILURE,
			Detail: fmt.Sprintf("Heartbeat failed, task not found or updated in the same second"),
		}, nil
	}

	return &pb.StatusMessage{
		Status: pb.Status_SUCCESS,
		Detail: fmt.Sprintf("Heartbeat successfull, task will now expire at: %v", expires),
	}, nil
}

// Complete moves a task from the running to the complete set
func (s *server) Complete(ctx context.Context, in *pb.Task) (*pb.StatusMessage, error) {

	queueID := in.QueueID
	taskID := in.TaskID

	_, err := s.rouge.Complete(queueID, taskID)
	if err != nil {
		return &pb.StatusMessage{Status: pb.Status_FAILURE, Detail: err.Error()}, err
	}

	return &pb.StatusMessage{
		Status: pb.Status_SUCCESS,
		Detail: "sucessfully marked item as complete",
	}, nil
}

// Fail moves a task from the running queue to the failed set
func (s *server) Fail(ctx context.Context, in *pb.Task) (*pb.StatusMessage, error) {

	queueID := in.QueueID
	taskID := in.TaskID

	_, err := s.rouge.Fail(queueID, taskID)
	if err != nil {
		return &pb.StatusMessage{Status: pb.Status_FAILURE, Detail: err.Error()}, err
	}

	return &pb.StatusMessage{
		Status: pb.Status_SUCCESS,
		Detail: "sucessfully marked item as failed",
	}, nil
}

// Progress returns a status struct about the requested queue
func (s *server) Progress(ctx context.Context, in *pb.RequestMessage) (*pb.QueueProgress, error) {
	queueInfo, err := s.rouge.Progress(in.QueueID)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, "could not get progress")
	}
	return queueInfo.ToBuff(), nil
}

func (s *server) ListQueues(ctx context.Context, in *pb.RequestMessage) (*pb.QueueList, error) {
	// return nil, and not implemented error
	return nil, grpc.Errorf(codes.Unimplemented, "method ListQueues not implemented")

	// queueList := &pb.QueueList{}
	// queues, err := s.rouge.ListQueues()
	// if err != nil {
	// 	return nil, grpc.Errorf(codes.Unknown, "could not get progress")
	// }
	// queueList.Queues = queues
	// return queueList, nil
}

// Peek returns a count and messageList
func (s *server) Peek(ctx context.Context, in *pb.RequestMessage) (*pb.TaskList, error) {
	// var task pb.Task
	taskList := &pb.TaskList{}

	// get TaskMessageList
	count, tml, err := s.rouge.Peek(in.QueueID, in.Phase, int(in.Limit))
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, "could not get progress")
	}

	taskList.TotalItems = int32(count)
	for i := 0; i < len(tml); i++ {
		task := &pb.Task{
			TaskID: tml[i].ID,
			Body:   tml[i].Body,
			Envs:   tml[i].Envs,
		}
		taskList.Tasks = append(taskList.Tasks, task)
	}
	return taskList, nil
}
