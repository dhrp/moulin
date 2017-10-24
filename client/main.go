package main

import (
	"fmt"
	"log"

	"flag"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nerdalize/moulin/certificates"
	pb "github.com/nerdalize/moulin/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

const (
	address     = "localhost:8042"
	defaultName = "world"
)

// GRPCDriver is the main instance
type GRPCDriver struct {
	connection *grpc.ClientConn
	client     pb.APIClient
}

var queueID = flag.String("queue", "batch", "Select a queue")

func newGRPCDriver() *GRPCDriver {

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

	gd := &GRPCDriver{connection: conn, client: apiClient}

	return gd
}

func (g GRPCDriver) getHealth() (status string) {
	// first do status
	r, err := g.client.Healthz(context.Background(), &empty.Empty{})
	if err != nil {
		log.Panic("could not greet")
	}
	log.Printf("Health: %s", r.Status)
	return r.Status
}

// LoadTask loads a task from the queue
// ToDo: add a timeout, for testing, and allow selecting queueID
func (g GRPCDriver) LoadTask() (task *pb.Task) {
	// then load a message
	t, err := g.client.LoadTask(context.Background(), &pb.RequestMessage{QueueID: "clientTestSuite"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Task: %s", t.TaskID)
	return t
}

func main() {

	grpcDriver := newGRPCDriver()
	// _ = grpcDriver
	defer grpcDriver.connection.Close()

}
