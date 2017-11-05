package client

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nerdalize/moulin/certificates"
	pb "github.com/nerdalize/moulin/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

const (
	address     = "localhost:8042"
	defaultName = "world"
)

// GRPCDriver is the main instance
type GRPCDriver struct {
	Connection *grpc.ClientConn
	client     pb.APIClient
}

var queueID = flag.String("queue", "batch", "Select a queue")

// NewGRPCDriver creates and initializes a new GRPC client and connection
func NewGRPCDriver() *GRPCDriver {

	keyPair, certPool := certificates.GetCert()
	_ = keyPair

	fmt.Println(queueID)

	var opts []grpc.DialOption
	creds := credentials.NewClientTLSFromCert(certPool, "localhost:8042")
	opts = append(opts, grpc.WithTransportCredentials(creds))
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	apiClient := pb.NewAPIClient(conn)

	gd := &GRPCDriver{Connection: conn, client: apiClient}

	return gd
}

func (g GRPCDriver) GetHealth() (status string) {
	// first do status
	r, err := g.client.Healthz(context.Background(), &empty.Empty{})
	if err != nil {
		log.Panic("could not greet")
	}
	log.Printf("Health: %s", r.Status)
	return r.Status
}

// PushTask loads a task from the queue
// ToDo: add a timeout, for testing, and allow selecting queueID
func (g GRPCDriver) PushTask(task *pb.Task) string {
	// then load a message

	md := metadata.Pairs("authorization", "open sesame")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	result, err := g.client.PushTask(ctx, task)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Result: %v", result)
	return result.Status
}

// LoadTask loads a task from the queue
// ToDo: add a timeout, for testing, and allow selecting queueID
func (g GRPCDriver) LoadTask(queueID string) (task *pb.Task) {
	// then load a message
	t, err := g.client.LoadTask(context.Background(), &pb.RequestMessage{QueueID: queueID})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Task: %s", t.TaskID)
	return t
}
