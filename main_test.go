package main

import (
	"context"
	"log"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/nerdalize/moulin/protobuf"
	"github.com/nerdalize/moulin/rouge"
	"github.com/stretchr/testify/suite"
)

type MainTestSuite struct {
	suite.Suite
	server *server
}

func (suite *MainTestSuite) SetupSuite() {

	rougeClient := &rouge.RougeClient{Host: "localhost:6379"}
	rougeClient.Init()
	suite.server = &server{rouge: rougeClient}
}

func (suite *MainTestSuite) TestHealthz() {
	req := &empty.Empty{}
	resp, _ := suite.server.Healthz(context.Background(), req)
	suite.NotEmpty(resp, "didnt get a response")
}

func (suite *MainTestSuite) TestPushTask() {
	log.Println("*** testing gRPC pushTask")

	// task := &rouge.TaskMessage{ID: "ASDF", Body: "empty"}

	req := &pb.Task{QueueID: "foobar", TaskID: "ASDF"}
	resp, _ := suite.server.PushTask(context.Background(), req)

	log.Println(resp)
	// suite.Equal(msg, "", "what was pushed is not what was popped")
}

func (suite *MainTestSuite) TestLoadTask() {
	log.Println("*** testing gRPC LoadTask")

	req := &pb.RequestMessage{QueueID: "foobar"}
	_, _ = suite.server.LoadTask(context.Background(), req)

}

func (suite *MainTestSuite) TearDownSuite() {
	log.Println("closing suite, cleaning up Redis")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
