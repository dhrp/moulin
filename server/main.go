package main

import (
	"fmt"
	"log"
	"net"

	"github.com/dhrp/moulin/rouge"
)

const (
	port     = 8042
	hostname = ""
)

// join the two constants for convenience
var serveAddress = fmt.Sprintf("%v:%d", hostname, port)

func main() {

	rougeClient, err := rouge.NewRougeClient()
	if err != nil {
		log.Panic(err.Error())
	}

	grpcServer := newGRPCServer(rougeClient)

	lis, err := net.Listen("tcp", serveAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
