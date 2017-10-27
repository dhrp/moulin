package main

// Here we define all GRPC handlers

import (
	"fmt"

	"google.golang.org/grpc/metadata"

	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/nerdalize/moulin/protobuf"
	"github.com/nerdalize/moulin/rouge"
)

func (s *server) Healthz(ctx context.Context, in *empty.Empty) (*pb.StatusMessage, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println("read incoming context successfull:")
		fmt.Println(md)
		fmt.Println(md["authentication"])
	} else {
		fmt.Println("failed to load context (meta)")
	}

	result, _ := s.rouge.Info()

	return &pb.StatusMessage{
		Status: "OK",
		Detail: result,
	}, nil
}

// LoadTask returns a task from the redis queue
func (s *server) LoadTask(ctx context.Context, in *pb.RequestMessage) (*pb.Task, error) {

	queueID := in.QueueID
	taskMessage := s.rouge.Load(queueID, 300)

	return &pb.Task{
		TaskID: taskMessage.ID,
		Body:   taskMessage.Body,
	}, nil
}

func (s *server) PushTask(ctx context.Context, in *pb.Task) (*pb.StatusMessage, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		var authorization string = md["authorization"][0]
		fmt.Println("authorization: " + authorization)
	} else {
		fmt.Println("failed to load context (meta)")
	}

	queueID := in.QueueID
	taskID := ksuid.New().String()

	taskMessage := rouge.TaskMessage{ID: taskID, Body: in.Body}
	queueLength, err := s.rouge.AddTask(queueID, taskMessage)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to push task")
	}
	result := fmt.Sprintf("queue %v now %d items long", queueID, queueLength)

	return &pb.StatusMessage{
		Status: "OK",
		Detail: result,
	}, nil
}
