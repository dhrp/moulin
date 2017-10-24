package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/nerdalize/moulin/rouge"
)

type MainTestSuite struct {
	suite.Suite
	grpcDriver *GRPCDriver
	rouge      *rouge.Client
}

// SetupSuite takes care of starting a rouge client
// and creating a server instance
func (suite *MainTestSuite) SetupSuite() {
	fmt.Println("*** SetupSuite()")

	// Test the error handling. We expect the a panic at this time.
	suite.Panics(suite.TestGetHealthz, "The function did not panic even though there is no connection!?!")

	grpcDriver := newGRPCDriver()
	suite.grpcDriver = grpcDriver

	// initialize the rouge client (on localhost)
	// rougeClient := &rouge.Client{Host: "localhost:6379"}
	// rougeClient.Init()

	// initialize the server, with our rougeClient
	// suite.rouge = &server{rouge: rougeClient}
	// suite.rouge = &rouge.Client{Host: "localhost:6379"}
	// _ = suite.rouge.Init()
}

func (suite *MainTestSuite) TestGetHealthz() {
	fmt.Println("*** TestGetHealthz()")
	result := suite.grpcDriver.getHealth()
	suite.Equal("OK", result, "Didn't receive OK health")
}

func (suite *MainTestSuite) TestLoadTask() {
	fmt.Println("*** TestLoadTask()")

	// ToDo: Set a timeout to loading task, and make a case where we add a task first.
	// task := suite.grpcDriver.LoadTask()
	// suite.Equal(len("0vNrL62AGAdIzRZ9pReEnKeMu44"), len(task.TaskID), "TaskID doesn't look valid")
}

func (suite *MainTestSuite) TearDownSuite() {
	log.Println("Tearing down test suite")
	log.Println("closing grpcDriver connection")
	suite.grpcDriver.connection.Close()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
