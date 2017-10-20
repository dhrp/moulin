package main

// Here we define all GRPC handlers

import (
	"golang.org/x/net/context"

	"github.com/pkg/errors"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/nerdalize/moulin/protobuf"
	"github.com/nerdalize/moulin/rouge"
)

func (s *server) Healthz(ctx context.Context, in *empty.Empty) (*pb.StatusMessage, error) {
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
	queueID := in.QueueID

	taskMessage := rouge.TaskMessage{ID: "nonsense", Body: "empty"}
	result, err := s.rouge.AddTask(queueID, taskMessage)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to push task")
	}

	return &pb.StatusMessage{
		Status: "OK",
		Detail: result,
	}, nil
}
