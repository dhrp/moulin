package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/nerdalize/moulin/helloworld"
	"google.golang.org/grpc/reflection"
	// "time"
	"github.com/go-redis/redis"
	"fmt"
	"os"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{
	r *redis.Client
	hostname string
}

// Open a new connection to Redis
func NewRedisClient() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     "dev.nlze.nl:6379",
		Password: "nevermind", // no password set
		DB:       0,  // use default DB
		// PoolSize:    300,
	})

	pong, err := client.Ping().Result()
	if err == nil {
		var _ = pong
		fmt.Println("redis client connected successfully")		
	} else {
		fmt.Println(err)		
	}
	return client
}

// Pop item from the queue, block untill one is available
func getFromQueue(client *redis.Client) string {

	fmt.Println("getting item from queue")
					   
	val, err := client.BRPopLPush("grpc:1", "grpc:2", 0).Result()

	if err == redis.Nil {
		fmt.Println("key does not exists")
	} else if err != nil {
		panic(err)
	}

	return val
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	// return &pb.HelloReply{Message: "Hello, you've connected successfully " + in.Name}, nil
	return &pb.HelloReply{Message: "You've connected to " + s.hostname}, nil
}

// Here we wait for an item on Redis, and then return an item.
func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("client connected")
	argument := getFromQueue(s.r)
	fmt.Println(argument + " retrieved from redis, sending to client")
    return &pb.HelloReply{Message: "You've been assigned " + argument}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	hostname, err := os.Hostname()
	fmt.Println("server started on " + hostname)

	s := grpc.NewServer()
	r := NewRedisClient()

	serverinstance := &server{}
	serverinstance.r = r
	serverinstance.hostname = hostname

	pb.RegisterGreeterServer(s, serverinstance)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
