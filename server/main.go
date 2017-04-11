package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/reflection"
	// "time"
	"github.com/go-redis/redis"
	"fmt"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{
	r *redis.Client
}

// Open a new connection to Redis
func NewRedisClient() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     "dev.nlze.nl:6379",
		Password: "nevermind", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return client
}

// Pop item from the queue, block untill one is available
func getFromQueue(client *redis.Client) string {

	fmt.Println(client)
					   
	val, err := client.BRPopLPush("grpc:1", "grpc:2", 0).Result()

	if err == redis.Nil {
		fmt.Println("key does not exists")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("queue item", val)
	}

	return val
}


// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello, you've connected successfully " + in.Name}, nil
}
func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {

	argument := getFromQueue(s.r)

	fmt.Println(argument)

    return &pb.HelloReply{Message: "Hello " + in.Name + " you've been assigned " + argument}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	s := grpc.NewServer()
	r := NewRedisClient()

	serverinstance := &server{}
	serverinstance.r = r

	pb.RegisterGreeterServer(s, serverinstance)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
