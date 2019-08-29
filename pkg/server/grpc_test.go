package server

import (
	"fmt"
	"log"

	"google.golang.org/grpc/metadata"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"

	pb "github.com/dhrp/moulin/pkg/protobuf"
)

func (suite *MainTestSuite) TestHealthz() {
	req := &empty.Empty{}
	resp, _ := suite.server.Healthz(context.Background(), req)
	suite.NotEmpty(resp, "didnt get a response")
}

func (suite *MainTestSuite) TestPushLoadAndCompleteTask() {
	log.Println("*** testing gRPC TestPushAndLoadTask")

	// create and push the task
	var status *pb.StatusMessage
	md := metadata.Pairs("authorization", "open sesame")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	reqOut := &pb.Task{QueueID: "foobar", TaskID: "ASDF"}
	status, _ = suite.server.PushTask(ctx, reqOut)
	suite.Equal(pb.Status_SUCCESS, status.Status, "Didn't get Success status from PushTask")

	// retrieve the task
	reqIn := &pb.RequestMessage{QueueID: "foobar"}
	var task *pb.Task
	var err error
	task, err = suite.server.LoadTask(context.Background(), reqIn)
	suite.Nil(err)
	fmt.Println(task)

	// heartbeat the task
	m := &pb.Task{QueueID: "foobar", TaskID: task.TaskID, ExpirationSec: 501}
	statusMessage, err := suite.server.HeartBeat(context.Background(), m)
	suite.Nil(err)
	suite.Equal(pb.Status_SUCCESS, statusMessage.Status)

	// complete the task
	m = &pb.Task{QueueID: "foobar", TaskID: task.TaskID}
	statusMessage, err = suite.server.Complete(context.Background(), m)
	suite.Nil(err)
	suite.Equal(pb.Status_SUCCESS, statusMessage.Status)
	//
	// complete the same task again task, this fails
	m = &pb.Task{QueueID: "foobar", TaskID: task.TaskID}
	statusMessage, err = suite.server.Complete(context.Background(), m)
	suite.NotNil(err)
	suite.Equal(pb.Status_FAILURE, statusMessage.Status)
}
