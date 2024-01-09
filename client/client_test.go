package client

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	pb "github.com/dhrp/moulin/pkg/protobuf"
	"github.com/dhrp/moulin/pkg/rouge"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	rougeClient, err := rouge.NewRougeClient()
	suite.NoError(err, "no error")

	suite.rouge = rougeClient
}

func (suite *MainTestSuite) TestGetHealthz() {
	fmt.Println("*** TestGetHealthz()")
	result, err := suite.grpcDriver.GetHealth()
	suite.Nil(err)
	suite.Equal(pb.Status_SUCCESS, result.Status, "Didn't receive OK health")
}

func (suite *MainTestSuite) TestLoadExpire() {
	fmt.Println("*** TestLoadExpire()")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := suite.grpcDriver.LoadTask(ctx, "clientTest")
	code := status.Code(err)

	fmt.Println(code)
	suite.Equal(code, codes.DeadlineExceeded)

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
	returnedTask, err := suite.grpcDriver.LoadTask(context.Background(), "clientTest")
	suite.Nil(err)
	suite.Equal(len("0vNrL62AGAdIzRZ9pReEnKeMu4x"), len(returnedTask.TaskID), "TaskID doesn't look valid")

	// ToDo: Set a timeout to loading task, and make a case where we add a task first.
	result = suite.grpcDriver.HeartBeat("clientTest", returnedTask.TaskID, 301)
	suite.Equal(pb.Status_SUCCESS, result.Status)

	result = suite.grpcDriver.HeartBeat("clientTest", "doesnt-exist", 301)
	suite.Equal(pb.Status_FAILURE, result.Status)

	result = suite.grpcDriver.Complete("clientTest", returnedTask.TaskID)
	suite.Equal(pb.Status_SUCCESS, result.Status)
}

// TestTaskConnectFirst is a test to show a problem where, if LoadTask is
// started before a task is on the queue, it will not return the first item
// added to that queue. It will add (and return) subsequent items..
// func TestTaskConnectFirst(t *testing.T) {
//
// 	// run this test ten times
// 	for i := 0; i < 10; i++ {
//
// 		grpcDriver := NewGRPCDriver()
//
// 		// initialize the rouge client (on localhost)
// 		rougeClient := &rouge.Client{Host: "localhost:6379"}
// 		rougeClient.Init()
//
// 		fmt.Println("*** TestLoadTask()")
// 		var inputTask *pb.Task
// 		var result *pb.StatusMessage
// 		queueID := "clientTest2"
//
// 		rougeClient.ClearQueue(queueID)
//
// 		channel := make(chan bool)
// 		go LoadRoutine(grpcDriver, 1, channel)
// 		// go LoadRoutine(grpcDriver, 2, channel)
//
// 		fmt.Print("[main] sleep 100 ms..\n")
// 		time.Sleep(100 * time.Millisecond)
// 		fmt.Print("[main] done sleeping\n")
//
// 		fmt.Print("push task..\n")
// 		inputTask = &pb.Task{
// 			QueueID: queueID,
// 			Body:    "Task #1",
// 		}
// 		result = grpcDriver.PushTask(inputTask)
//
// 		progress, _ := grpcDriver.Progress(queueID)
// 		fmt.Printf("incoming: %d\n", progress.IncomingCount)
// 		fmt.Printf("received: %d\n", progress.ReceivedCount)
// 		fmt.Printf("running:  %d\n", progress.RunningCount)
//
// 		loadResult := <-channel
//
// 		assert.True(t, loadResult, "channel didn't reply true (in time?)")
// 		fmt.Println("received true from loadroutine")
//
// 		if loadResult != true {
// 			return
// 		}
//
// 		assert.Equal(t, pb.Status_SUCCESS, result.Status, "result was not OK")
// 		grpcDriver.Connection.Close()
//
// 	}
// }

// func (suite *MainTestSuite) TearDownSuite() {
// 	log.Println("Tearing down test suite")
// 	log.Println("closing grpcDriver connection")
// 	suite.grpcDriver.Connection.Close()
// }

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
