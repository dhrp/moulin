package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

// TaskMessage is the 'message' format used on our queue
type TaskMessage struct {
	ID   int64    `json:"id"`
	Body string   `json:"body"`
	Envs []string `json:"envs"`
}

func newScore() string {
	// determine the score (when task should expire)
	timestamp := int64(time.Now().Unix())
	expires := timestamp + 300 // 5 min
	return strconv.FormatInt(expires, 10)
}

// Load loads a message from the queue, and mark it as processing
func Load(red RedClient, queue string) bool {

	// Get taskMessage (string) from queue
	rawMsgString := red.brpop(queue)

	// Get the taskID from the taskMessage
	var taskMessage TaskMessage
	err := json.Unmarshal([]byte(rawMsgString), &taskMessage)
	if err != nil {
		log.Println(err)
		// Here we should really log something, or notify; as the json
		// cannot be parsed
		return false
	}

	// taskId := taskMessage.ID
	taskID := strconv.FormatInt(taskMessage.ID, 10)

	// save taskId to sorted set with timeout
	red.set(queue+"."+taskID, rawMsgString)

	/// save the item to the sorted set
	set := queue + ".running"
	score := newScore()
	value := queue + "." + taskID

	// ZADD Q_working_set <now>+300 queue-id.task-id
	red.zadd(set, score, value)

	return true
}

// heartbeat updates the status of a message
func Heartbeat(red RedClient, queueID string, taskID string, score string) bool {

	set := queueID + ".running"
	value := queueID + "." + taskID

	// _, _, _ = set, score, value
	count, _ := red.zaddUpdate(set, score, value)
	if count == 0 {
		log.Println("Heartbeat: no item could be found to update, was item already " +
			"completed?, or perhaps the item was updated < 1 sec ago.")
		return false
	}
	return true
}

// complete marks the item as completed.
func Complete(red RedClient, queueID string, taskID string) bool {

	// count, _ := red.moveValueFromToSet(from, to, value)(int, error)
	return true
}

func main() {

	red := RedClient{host: "localhost:6379"}
	_ = red.init()

	Load(red, "queue")

	fmt.Println("DEBUGGING")

	// _ = red.set(taskId, taskMessage)

	// fmt.Println(taskMessage.body)

}
