package main

import (
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	// address     = "localhost:50051"
	// address     = "dev.nlze.nl:50051"
	defaultName = "world"
	
	defaultServer = "localhost"
	defaultPort = "50051"
)

func main() {
	// Set up a connection to the server.
	server := defaultServer
	port := defaultPort

	if len(os.Args) > 1 {
		server = os.Args[1]
	}
	if len(os.Args) > 2 {
		port = os.Args[2]
	}

	remoteServer := server + ":" + port

	conn, err := grpc.Dial(remoteServer, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)


	// Contact the server and print out its response.
	name := defaultName

	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)


	// second function
	r, err = c.SayHelloAgain(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
	        log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
