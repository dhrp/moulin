package main

import (
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"

	pb "github.com/nerdalize/moulin/protobuf"
)

func (suite *MainTestSuite) TestHealthz() {
	req := &empty.Empty{}
	resp, _ := suite.server.Healthz(context.Background(), req)
	suite.NotEmpty(resp, "didnt get a response")
}

func (suite *MainTestSuite) TestPushTask() {
	log.Println("*** testing gRPC pushTask")

	// task := &rouge.TaskMessage{ID: "ASDF", Body: "empty"}

	var status *pb.StatusMessage
	ctx := context.Background()

	req := &pb.Task{QueueID: "foobar", TaskID: "ASDF"}
	status, _ = suite.server.PushTask(ctx, req)

	suite.Equal("OK", status.Status, "Didn't get OK status from PushTask")
}

func (suite *MainTestSuite) TestLoadTask() {
	log.Println("*** testing gRPC LoadTask")

	req := &pb.RequestMessage{QueueID: "foobar"}
	_, _ = suite.server.LoadTask(context.Background(), req)

}
