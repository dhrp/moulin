package process

import (
	"log"

	"github.com/dhrp/moulin/client"
)

func (suite *MainTestSuite) TestWork() {
	log.Println("*** testing Work")
	var result int
	var err error

	grpcDriver := client.NewGRPCDriver()
	defer grpcDriver.Connection.Close()

	result, err = Work(grpcDriver, "q1", "once")

	_, _ = result, err

	// func Work(grpcDriver *client.GRPCDriver, queueID, workType string) (result int, err error) {
}
