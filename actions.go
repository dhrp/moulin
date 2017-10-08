package main

import (
	// "errors"
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"

	pb "github.com/nerdalize/moulin/helloworld"
	"github.com/nerdalize/moulin/rouge"
)

func createTaskFromJSON(json []byte) (*rouge.TaskMessage, error) {

	ID, _ := ksuid.NewRandom()
	IDstr := ID.String()

	task := &rouge.TaskMessage{ID: IDstr}
	_, err := task.FromString(json)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create task from string")
	}
	if task.Body == "" {
		return nil, errors.New("Error: JSON needs to have at least 'body' field")
	}
	result := fmt.Sprintf("Task with body: %v, set ID %v \n", task.Body, task.ID)
	log.Println(result)

	return task, nil
}

// SayHello implements helloworld.GreeterServer
func SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "You've connected"}, nil
}
