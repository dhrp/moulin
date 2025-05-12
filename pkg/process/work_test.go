package process

import (
	"log"

	"github.com/dhrp/moulin/client"
	pb "github.com/dhrp/moulin/pkg/protobuf"
)

func (suite *ProcessTestSuite) TestWork() {
	log.Println("*** testing Work")
	var result int
	var err error

	grpcDriver := client.NewGRPCDriver()
	defer grpcDriver.Connection.Close()

	task := new(pb.Task)
	task.QueueID = "q1"
	task.Body = `sh -c "echo this is a test && echo done"`
	grpcDriver.PushTask(task)

	result, err = Work(grpcDriver, "q1", "once")
	suite.Nil(err)

	_, _ = result, err

	// func Work(grpcDriver *client.GRPCDriver, queueID, workType string) (result int, err error) {
}
