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
		log.Fatalln(err.Error())
		// os.Exit(1)
	}

	grpcServer := newGRPCServer(rougeClient)

	lis, err := net.Listen("tcp", serveAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Starting server on %s\n", serveAddress)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
