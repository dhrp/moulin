package client

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"

	pb "github.com/nerdalize/moulin/protobuf"
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

	grpcDriver := NewGRPCDriver()
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
	result := suite.grpcDriver.GetHealth()
	suite.Equal(pb.Status_SUCCESS, result.Status, "Didn't receive OK health")
}

func (suite *MainTestSuite) TestOneTaskEndToEnd() {
	fmt.Println("*** TestLoadTask()")

	inputTask := &pb.Task{
		QueueID: "clientTest",
		// TaskID: taskID,  // we let the server create a taskID
		Body: "Just Testing!",
	}

	result := suite.grpcDriver.PushTask(inputTask)
	suite.Equal(pb.Status_SUCCESS, result.Status, "result was not OK")

	// ToDo: Set a timeout to loading task, and make a case where we add a task first.
	returnedTask := suite.grpcDriver.LoadTask("clientTest")
	suite.Equal(len("0vNrL62AGAdIzRZ9pReEnKeMu4x"), len(returnedTask.TaskID), "TaskID doesn't look valid")

	// ToDo: Set a timeout to loading task, and make a case where we add a task first.
	result = suite.grpcDriver.HeartBeat("clientTest", returnedTask.TaskID, 301)
	suite.Equal(pb.Status_SUCCESS, result.Status)

	result = suite.grpcDriver.HeartBeat("clientTest", "doesnt-exist", 301)
	suite.Equal(pb.Status_FAILURE, result.Status)

	result = suite.grpcDriver.Complete("clientTest", returnedTask.TaskID)
	suite.Equal(pb.Status_SUCCESS, result.Status)
}

func (suite *MainTestSuite) TearDownSuite() {
	log.Println("Tearing down test suite")
	log.Println("closing grpcDriver connection")
	suite.grpcDriver.Connection.Close()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
