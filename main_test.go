package main

import (
	"log"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/dhrp/moulin/rouge"
)

type MainTestSuite struct {
	suite.Suite
	server *server
}

// SetupSuite takes care of starting a rouge client
// and creating a server instance
func (suite *MainTestSuite) SetupSuite() {
	// initialize the rouge client (on localhost)
	rougeClient := &rouge.Client{Host: "localhost:6379"}
	rougeClient.Init()

	// initialize the server, with our rougeClient
	suite.server = &server{rouge: rougeClient}
}

// func (suite *MainTestSuite) TestCreateGlobalServer() {
// 	globalServer := createGlobalServer()
// 	log.Println(globalServer)
//
// 	suite.Equal("localhost:8042", globalServer.Addr, "host address doesn't match what we set")
// 	suite.Equal(1, len(globalServer.TLSConfig.Certificates), "Server should have a TLS certificate set")
// }

func (suite *MainTestSuite) TearDownSuite() {
	log.Println("closing suite, This would be a good place to close and clean up things")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
