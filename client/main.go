package main

import (
	"log"

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

func main() {
	keyPair, certPool := certificates.GetCert()
	_ = keyPair

	var opts []grpc.DialOption
	creds := credentials.NewClientTLSFromCert(certPool, "localhost:8042")
	opts = append(opts, grpc.WithTransportCredentials(creds))
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	c := pb.NewAPIClient(conn)

	// first do status
	r, err := c.Healthz(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Health: %s", r.Status)

	// then load a message
	t, err := c.LoadTask(context.Background(), &pb.RequestMessage{QueueID: "foobar"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Task: %s", t.TaskID)

}
