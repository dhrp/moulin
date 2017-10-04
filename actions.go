package main

import (
	// "errors"
	"fmt"
	"log"

	"github.com/nerdalize/moulin/rouge"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
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
