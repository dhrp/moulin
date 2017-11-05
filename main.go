package main

// Here we define all GRPC handlers

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"

	"github.com/nerdalize/moulin/certificates"
	// "github.com/nerdalize/moulin/kafkaproducer"
	pb "github.com/nerdalize/moulin/protobuf"
	"github.com/nerdalize/moulin/rouge"
)

const (
	port     = 8042
	hostname = "localhost"
)

// join the two constants for convenience
var serveAddress string = fmt.Sprintf("%v:%d", hostname, port)

type server struct {
	rouge *rouge.Client
}

func (s *server) makeGRPCServer(certPool *x509.CertPool) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewClientTLSFromCert(certPool, fmt.Sprintf("%v:%d", hostname, port)))}

	//setup grpc server
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAPIServer(grpcServer, s)
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	return grpcServer
}

// getRestMux initializes a new multiplexer, and registers each endpoint
// - in this case only the EchoService
func getRestMux(certPool *x509.CertPool, opts ...runtime.ServeMuxOption) (*runtime.ServeMux, error) {

	// Because we run our REST endpoint on the same port as the GRPC the address is the same.
	upstreamGRPCServerAddress := serveAddress

	// get context, this allows control of the connection
	ctx := context.Background()

	// These credentials are for the upstream connection to the GRPC server
	dcreds := credentials.NewTLS(&tls.Config{
		ServerName: upstreamGRPCServerAddress,
		RootCAs:    certPool,
	})
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	// Which multiplexer to register on.
	gwmux := runtime.NewServeMux()
	err := pb.RegisterAPIHandlerFromEndpoint(ctx, gwmux, upstreamGRPCServerAddress, dopts)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
		return nil, err
	}

	return gwmux, nil
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

// createGlobalServer initializes rouge, the GRPC server and the HTTP mux,
// it returns a fully loaded http server
func createGlobalServer() *http.Server {

	rougeClient := &rouge.Client{Host: "localhost:6379"}
	rougeClient.Init()
	// kfk := &kafkaproducer.KFK{Broker: "localhost:9092"}
	// kfk.Init()

	server := &server{rouge: rougeClient}

	keyPair, certPool := certificates.GetCert()

	grpcServer := server.makeGRPCServer(certPool)
	restMux, err := getRestMux(certPool)
	if err != nil {
		log.Panic(err)
	}

	// register root Http multiplexer (mux)
	mux := http.NewServeMux()

	// we can add any non-grpc endpoints here.
	mux.HandleFunc("/foobar/", simpleHTTPHello)
	mux.HandleFunc("/v1/task_list/batch/", server.createTaskListBatch)
	// router.POST("/v1/task_list/batch/", createTaskListBatch)

	// register the gateway mux onto the root path.
	mux.Handle("/", restMux)

	// the grpcHandlerFunc takes an grpc server and a http muxer and will
	// route the request to the right place at runtime.
	mergeHandler := grpcHandlerFunc(grpcServer, mux)

	// configure TLS for our server. TLS is REQUIRED to make this setup work.
	// check https://golang.org/src/net/http/server.go?s=69823:69872#L2666
	srv := &http.Server{
		Addr:    serveAddress,
		Handler: mergeHandler,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*keyPair},
			NextProtos:   []string{"h2"},
		},
	}
	return srv
}

func main() {

	var srv *http.Server = createGlobalServer()

	// start listening on the socket
	// Note that if you listen on localhost:<port> you'll not be able to accept
	// connections over the network. Change it to ":port"  if you want it.
	conn, err := net.Listen("tcp", serveAddress)
	if err != nil {
		panic(err)
	}

	// start the server
	fmt.Printf("starting GRPC and REST on: %v\n", serveAddress)
	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
