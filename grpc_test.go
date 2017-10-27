package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	pb "github.com/nerdalize/moulin/protobuf"
)

func (suite *MainTestSuite) TestHealthz() {
	req := &empty.Empty{}
	resp, _ := suite.server.Healthz(context.Background(), req)
	suite.NotEmpty(resp, "didnt get a response")
}

func (suite *MainTestSuite) TestPushAndLoadTask() {
	log.Println("*** testing gRPC TestPushAndLoadTask")

	// create and push the task
	var status *pb.StatusMessage
	md := metadata.Pairs("authorization", "open sesame")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	reqOut := &pb.Task{QueueID: "foobar", TaskID: "ASDF"}
	status, _ = suite.server.PushTask(ctx, reqOut)
	suite.Equal("OK", status.Status, "Didn't get OK status from PushTask")

	// retrieve the task
	reqIn := &pb.RequestMessage{QueueID: "foobar"}
	var task *pb.Task
	var err error
	task, err = suite.server.LoadTask(context.Background(), reqIn)
	suite.Nil(err)
	fmt.Println(task)

}

//
// func (suite *MainTestSuite) TestLoadTask() {
// 	log.Println("*** testing gRPC LoadTask")
//
// 	// req := &pb.RequestMessage{QueueID: "foobar"}
// 	// _, _ = suite.server.LoadTask(context.Background(), req)
//
// }
