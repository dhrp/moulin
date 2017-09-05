package main

import (
	"log"
	"net"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	pb "github.com/nerdalize/moulin/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"fmt"
	"os"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	r        *pool.Pool
	hostname string
}

// Open a new connection to Redis
func NewRedisClient() *pool.Pool {
	df := func(network, addr string) (*redis.Client, error) {
		client, err := redis.Dial(network, addr)
		if err != nil {
			return nil, err
		}
		if err = client.Cmd("AUTH", "nevermind").Err; err != nil {
			client.Close()
			return nil, err
		}
		return client, nil
	}

	client, err := pool.NewCustom("tcp", "localhost:6379", 10, df)
	if err == nil {
		fmt.Println("redis client connected successfully with radix driver")
	} else {
		fmt.Println(err)
	}

	fmt.Println(client)
	return client
}

// Pop item from the queue, block untill one is available
func getFromQueue(client *pool.Pool) string {

	fmt.Println("getting item from queue")
	val, err := client.Cmd("BRPOPLPUSH", "grpc:1", "grpc:2", 0).Str()
	fmt.Println(val)

	if err != nil {
		panic(err)
	}

	fmt.Println("got " + val)
	return val
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
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

	serverinstance := &server{r: r, hostname: hostname}
	// serverinstance := &server{}
	// serverinstance.r = r
	// serverinstance.hostname = hostname

	pb.RegisterGreeterServer(s, serverinstance)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
