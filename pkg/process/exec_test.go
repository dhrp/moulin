package process

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ProcessTestSuite struct {
	suite.Suite
}

//
// func (suite *MainTestSuite) TestExec() {
// 	log.Println("*** testing Exec")
// 	var task *pb.Task
//
// 	type testCommand struct {
// 		command  string
// 		exitCode int
// 	}
//
// 	commands := []testCommand{
// 		testCommand{command: "echo hello world", exitCode: 0},
// 		testCommand{command: "echo hello; false", exitCode: 0}, // this is echo "hello; false"
// 		testCommand{command: "/bin/bash -c \"echo XXX\"", exitCode: 0},
// 		testCommand{command: "/bin/bash -c \"echo hello; echo sending error; false\"", exitCode: 1},
// 		testCommand{command: "", exitCode: 1},
// 		testCommand{command: "doesntexist", exitCode: 1},
// 	}
//
// 	for i := 0; i < len(commands); i++ {
// 		task = &pb.Task{
// 			QueueID: "1010",
// 			Body:    commands[i].command,
// 		}
// 		statusCode, err := Exec(task)
// 		suite.Nil(err)
// 		suite.Equal(commands[i].exitCode, statusCode)
// 	}
// }

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestProcessTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessTestSuite))
}
