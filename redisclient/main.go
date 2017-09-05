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

func (t *TaskMessage) toString() string {
	b, err := json.Marshal(t)
	if err != nil {
		log.Panic("error:", err)
	}
	return string(b)
}

func (t *TaskMessage) fromString(jsonStr string) *TaskMessage {

	err := json.Unmarshal([]byte(jsonStr), &t)
	if err != nil {
		log.Panic(err)
	}
	return t
}

func newScore() string {
	// determine the score (when task should expire)
	timestamp := int64(time.Now().Unix())
	expires := timestamp + 300 // 5 min
	return strconv.FormatInt(expires, 10)
}

// Load loads a message from the queue, and mark it as processing
func Load(red RedClient, queueID string) TaskMessage {

	log.Println("**********")
	log.Println("LOAD START")

	receivedQueue := queueID + ".received"

	// try if there is an expired task to pick up
	var taskMessage TaskMessage

	// try if there is a task on the received queue to pick up
	msg, err := red.rpop(receivedQueue)
	if err == nil {
		taskMessage.fromString(msg)
		return taskMessage
	}

	// switch a task from the queue to the received queue
	red.brpoplpush(queueID, receivedQueue)

	// fetch from received queue, save the message to it's key
	// and add the ID to the running set

	destinationSet := fmt.Sprintf("%s.running", queueID)
	msg, err = red.popQueueAndSaveKeyToSet(receivedQueue, destinationSet, 300)
	if err == nil {
		taskMessage.fromString(msg)
		fmt.Println(taskMessage)
	}

	// return taskMessage
	// // pop it from the received queue and set it on a key
	// msg, err = red.rpop(receivedQueue)
	// if err == nil {
	// 	taskMessage.fromString(msg)
	// }

	// taskId := taskMessage.ID
	// taskID := strconv.FormatInt(taskMessage.ID, 10)
	// key := fmt.Sprintf("%s.%s", queueID, taskID)

	// save taskMessage (string) to a key.
	// red.set(key, rawMsgString)w

	/// save the item to the sorted set
	// set := queueID + ".running"
	// score := newScore()
	// member := taskID

	// ZADD Q_working_set <now>+300 queue-id.task-id
	// red.zadd(set, score, member)

	log.Println("LOAD END")
	log.Println("**********")

	return taskMessage
}

// Heartbeat updates the status of a message
func Heartbeat(red RedClient, queueID string, taskID string, expirationSec int) bool {

	log.Println("***************")
	log.Println("HEARTBEAT START")

	set := queueID + ".running"
	member := taskID

	expiresAt := int64(time.Now().Unix()) + int64(expirationSec)
	score := strconv.FormatInt(expiresAt, 10)

	// _, _, _ = set, score, value
	count, _ := red.zaddUpdate(set, score, member)
	if count == 0 {
		log.Println("Heartbeat: no item could be found to update, was item already " +
			"completed?, or perhaps the item was updated < 1 sec ago.")
		return false
	}

	log.Println("HEARTBEAT END")
	log.Println("***************")

	return true
}

// Complete marks the item as completed.
func Complete(red RedClient, queueID string, taskID string) bool {

	from := fmt.Sprintf("%s.running", queueID)
	to := fmt.Sprintf("%s.complete", queueID)
	member := taskID

	count, _ := red.moveMemberFromSetToSet(from, to, member)
	if count != 1 {
		log.Printf("Didn't complete the right amount of items!: %d", count)
		return false
	}
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
