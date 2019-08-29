package process

import (
	"log"

	"github.com/dhrp/moulin/pkg/client"
	pb "github.com/dhrp/moulin/protobuf"
)

func (suite *MainTestSuite) TestWork() {
	log.Println("*** testing Work")
	var result int
	var err error

	grpcDriver := client.NewGRPCDriver()
	defer grpcDriver.Connection.Close()

	task := new(pb.Task)
	task.QueueID = "q1"
	task.Body = "echo this is a test"
	grpcDriver.PushTask(task)

	result, err = Work(grpcDriver, "q1", "once")

	_, _ = result, err

	// func Work(grpcDriver *client.GRPCDriver, queueID, workType string) (result int, err error) {
}
