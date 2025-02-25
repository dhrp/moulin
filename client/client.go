package client

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/dhrp/moulin/pkg/protobuf"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-retry"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	defaultName = "world"
)

// GRPCDriver is the main instance
type GRPCDriver struct {
	Connection *grpc.ClientConn
	client     pb.APIClient
}

// MoulinConnectionConfig defines a configuration for timeouts
type MoulinConnectionConfig struct {
	LoadTaskTimeOut          time.Duration // LoadTaskTimeOut is how long we wait for a task to be loaded (when we are doing `until finished`)
	HeartBeatInterval        time.Duration // HeartBeatInterval is the time between heartbeats
	ServerUnavailableTimeOut time.Duration // ServerUnavailableTimeOut is how long we accept the server to be down before we consider it dead and quit
	TaskExpirationTime       time.Duration // taskExpirationTime is how much time is allowed before the task is considered expired
}

// ClientConfig is the (default) configuration for the client (timeouts etc)
var ClientConfig = MoulinConnectionConfig{
	LoadTaskTimeOut:          30 * time.Second,
	HeartBeatInterval:        2 * time.Minute,
	ServerUnavailableTimeOut: 1 * time.Hour,
	TaskExpirationTime:       5 * time.Minute,
}

// NewGRPCDriver creates and initializes a new GRPC client and connection
func NewGRPCDriver() *GRPCDriver {

	address := os.Getenv("MOULIN_SERVER")
	if address == "" {
		address = "localhost:8042"
	}
	fmt.Printf("connecting to moulinServer on %s\n", address)

	var backoffConfig = backoff.Config{
		MaxDelay: 10 * time.Minute,
	}

	connectParams := grpc.ConnectParams{
		Backoff:           backoffConfig,
		MinConnectTimeout: ClientConfig.ServerUnavailableTimeOut,
	}

	conn, err := grpc.Dial(address, grpc.WithConnectParams(connectParams), grpc.WithInsecure())
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
func (g GRPCDriver) PushTask(task *pb.Task) (*pb.StatusMessage, error) {
	// then load a message

	md := metadata.Pairs("authorization", "open sesame")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	result, err := g.client.PushTask(ctx, task)
	return result, err
}

// LoadTask loads a task from the queue
func (g GRPCDriver) LoadTask(ctx context.Context, queueID string) (task *pb.Task, err error) {
	// then load a message
	t, err := g.client.LoadTask(ctx, &pb.RequestMessage{QueueID: queueID})
	return t, err
}

// HeartBeat updates the expiry of an item on the running set
// ToDo: add a timeout, for testing
func (g GRPCDriver) HeartBeat(queueID, taskID string) (*pb.StatusMessage, error) {
	// then load a message
	task := &pb.Task{
		QueueID:       queueID,
		TaskID:        taskID,
		ExpirationSec: int32(ClientConfig.TaskExpirationTime.Seconds()),
	}
	r, err := g.client.HeartBeat(context.Background(), task)
	if err != nil {
		log.Printf("could not complete heartbeat: %v", err)
	}
	return r, err
}

// Complete moves the task from the running set to the completed set
func (g GRPCDriver) Complete(queueID, taskID string) *pb.StatusMessage {

	task := &pb.Task{
		QueueID: queueID,
		TaskID:  taskID,
	}

	// For the "complete" action we implement an incremental retry function. This is because
	// the server could be down, and if it is, we want to retry the action, as it
	// could prevent a lot of work from being lost.
	ctx := context.Background()
	b := retry.NewFibonacci(1 * time.Second)
	b = retry.WithMaxDuration(20*time.Second, b)

	r, err := retry.DoValue(ctx, b, func(ctx context.Context) (*pb.StatusMessage, error) {
		res, err := g.client.Complete(ctx, task)
		if err != nil {
			log.Printf("Completing the task failed; is the server down? Retrying...")
			return nil, retry.RetryableError(fmt.Errorf("bad response: %v", err))
		}
		return res, nil
	})

	if err != nil {
		// Unpack the gRPC error
		st, _ := status.FromError(err)
		log.Fatalf("Error: Could not complete task: %v", st.Message())
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
		log.Fatalf("could not fail task: %v", err)
	}
	return r
}

// Progress gets the status for a queue
func (g GRPCDriver) Progress(queueID string) (progress *pb.QueueProgress, err error) {
	// then load a message
	progress, err = g.client.Progress(context.Background(), &pb.RequestMessage{QueueID: queueID})
	if err != nil {
		return progress, status.Errorf(codes.Unknown, "could not get progress")
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
		return taskList, status.Errorf(codes.Unknown, err.Error())
	}
	return taskList, nil
}

// ListQueues returns a list of Progress structs
func (g GRPCDriver) ListQueues() (queues map[string]*pb.QueueProgress, err error) {
	queueMap, err := g.client.ListQueues(context.Background(), &empty.Empty{})
	if err != nil {
		st, _ := status.FromError(err)
		return nil, errors.New(st.Message())
	}
	return queueMap.Queues, nil
}

// DeleteQueue deletes a queue
func (g GRPCDriver) DeleteQueue(queueID string) (*pb.StatusMessage, error) {
	res, err := g.client.DeleteQueue(context.Background(), &pb.RequestMessage{QueueID: queueID})
	return res, err
}
