package main

// Here we define all GRPC handlers

import (
	"fmt"
	"log"

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
		Status: pb.Status_SUCCESS,
		Detail: result,
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
		Status: pb.Status_SUCCESS,
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

// LoadTask returns a task from the redis queue
func (s *server) HeartBeat(ctx context.Context, in *pb.Task) (*pb.StatusMessage, error) {

	queueID := in.QueueID
	taskID := in.TaskID

	var expirationSec int32
	if in.ExpirationSec != 0 {
		expirationSec = in.ExpirationSec
	} else {
		expirationSec = 300 // default expiration is 5 min
	}

	expires, err := s.rouge.Heartbeat(queueID, taskID, expirationSec)
	if err != nil {
		// ToDo return error to client, it could be because item was just updated
		log.Println(err)
		return &pb.StatusMessage{
			Status: pb.Status_FAILURE,
			Detail: fmt.Sprintf("Heartbeat failed, task not found or updated in the same second"),
		}, nil
	}

	return &pb.StatusMessage{
		Status: pb.Status_SUCCESS,
		Detail: fmt.Sprintf("Heartbeat successfull, task will now expire at: %v", expires),
	}, nil
}

// Complete moves a task from the running to the complete set
func (s *server) Complete(ctx context.Context, in *pb.Task) (*pb.StatusMessage, error) {

	queueID := in.QueueID
	taskID := in.TaskID

	_, err := s.rouge.Complete(queueID, taskID)
	if err != nil {
		return &pb.StatusMessage{Status: pb.Status_FAILURE, Detail: err.Error()}, err
	}

	return &pb.StatusMessage{
		Status: pb.Status_SUCCESS,
		Detail: "sucessfully marked item as complete",
	}, nil
}
