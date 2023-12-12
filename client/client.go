package client

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/dhrp/moulin/pkg/protobuf"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

const (
	defaultName = "world"
)

// GRPCDriver is the main instance
type GRPCDriver struct {
	Connection *grpc.ClientConn
	client     pb.APIClient
}

// NewGRPCDriver creates and initializes a new GRPC client and connection
func NewGRPCDriver() *GRPCDriver {

	address := os.Getenv("MOULIN_SERVER")
	if address == "" {
		address = "localhost:8042"
		fmt.Println("setting moulinServer to localhost:8042")
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	sigchan := make(chan os.Signal, 2)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigchan
		log.Println("closing connection")
		conn.Close()
		os.Exit(1)
	}()

	apiClient := pb.NewAPIClient(conn)
	return &GRPCDriver{Connection: conn, client: apiClient}
}

// GetHealth just checks if everything, including Redis is healthy
func (g GRPCDriver) GetHealth() (status pb.StatusMessage, err error) {
	// first do status
	r, err := g.client.Healthz(context.Background(), &empty.Empty{})
	if err != nil {
		return pb.StatusMessage{}, errors.Wrap(err, "could not get healthz")
	}
	log.Printf("Health: %s", r.Status)
	return *r, nil
}

// PushTask pushes a task onto the queue
func (g GRPCDriver) PushTask(task *pb.Task) *pb.StatusMessage {
	// then load a message

	md := metadata.Pairs("authorization", "open sesame")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	result, err := g.client.PushTask(ctx, task)
	if err != nil {
		log.Fatalf("could not push task: %v", err)
	}
	log.Printf("Result: %v", result)
	return result
}

// LoadTask loads a task from the queue
func (g GRPCDriver) LoadTask(queueID string) (task *pb.Task, err error) {
	// then load a message
	t, err := g.client.LoadTask(context.Background(), &pb.RequestMessage{QueueID: queueID})
	if err != nil {
		log.Fatalf("could not load task: %v", err)
	}
	return t, nil
}

// HeartBeat updates the expiry of an item on the running set
// ToDo: add a timeout, for testing
func (g GRPCDriver) HeartBeat(queueID, taskID string, expirationSec int32) *pb.StatusMessage {
	// then load a message
	task := &pb.Task{
		QueueID:       queueID,
		TaskID:        taskID,
		ExpirationSec: expirationSec,
	}

	r, err := g.client.HeartBeat(context.Background(), task)
	if err != nil {
		log.Fatalf("could not complete heartbeat: %v", err)
	}
	return r
}

// Complete moves the task from the running set to the completed set
func (g GRPCDriver) Complete(queueID, taskID string) *pb.StatusMessage {

	task := &pb.Task{
		QueueID: queueID,
		TaskID:  taskID,
	}

	r, err := g.client.Complete(context.Background(), task)
	if err != nil {
		log.Fatalf("could not complete task: %v", err)
	}
	return r
}

// Fail marks the task as failed by pushing it to the failed set
func (g GRPCDriver) Fail(queueID, taskID string) *pb.StatusMessage {

	task := &pb.Task{
		QueueID: queueID,
		TaskID:  taskID,
	}

	r, err := g.client.Fail(context.Background(), task)
	if err != nil {
		log.Fatalf("could not complete task: %v", err)
	}
	return r
}

// Progress gets the status for a queue
func (g GRPCDriver) Progress(queueID string) (progress *pb.QueueProgress, err error) {
	// then load a message
	progress, err = g.client.Progress(context.Background(), &pb.RequestMessage{QueueID: queueID})
	if err != nil {
		return progress, grpc.Errorf(codes.Unknown, "could not get progress")
	}
	return progress, nil
}

// Peek get the n (limit) 'next' messages for a given queue/phase
func (g GRPCDriver) Peek(queueID, phase string, limit int32) (taskList *pb.TaskList, err error) {
	// peek into queue phase
	requestMessage := &pb.RequestMessage{
		QueueID: queueID,
		Phase:   phase,
		Limit:   limit}
	taskList, err = g.client.Peek(context.Background(), requestMessage)
	if err != nil {
		log.Printf(err.Error())
		return taskList, grpc.Errorf(codes.Unknown, err.Error())
	}
	log.Printf("Incoming queue length: %d", taskList.TotalItems)
	return taskList, nil
}

// ListQueues returns a list of Progress structs
func (g GRPCDriver) ListQueues() (queues []*pb.QueueProgress, err error) {
	// then load a message
	// queues, err = g.client.ListQueues(context.Background(), &empty.Empty{})
	// if err != nil {
	// 	return queues, grpc.Errorf(codes.Unknown, "could not get progress")
	// }
	// return nil and an not implemented error
	return nil, grpc.Errorf(codes.Unimplemented, "not implemented")
}
