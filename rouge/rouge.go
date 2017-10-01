package rouge

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

// RedClient Basic client class to group Redis functions
type RedClient struct {
	clientpool *pool.Pool
	host       string
	client     *redis.Client
}

func (red *RedClient) init() *RedClient {

	df := func(network, addr string) (*redis.Client, error) {
		client, err := redis.Dial(network, addr)
		if err != nil {
			return nil, err
		}
		// TODO: Review if we need a password
		// set password with CONFIG SET requirepass "nevermind"
		// if err = client.Cmd("AUTH", "nevermind").Err; err != nil {
		// 	client.Close()
		// 	return nil, err
		// }
		return client, nil
	}

	client, err := pool.NewCustom("tcp", red.host, 10, df)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("redis client connected successfully with radix driver")
	red.clientpool = client

	return red
}

// Load loads a message from the queue, and mark it as processing
func (red *RedClient) Load(queueID string, expirationSec int) TaskMessage {

	log.Println("**********")
	log.Println("LOAD START")

	receivedList := queueID + ".received"
	runningSet := queueID + ".running"
	var taskMessage TaskMessage

	// try if there is an expired task to pick up, and if so, return it.
	member, err := red.fetchAndUpdateExpired(runningSet, expirationSec)
	if err != nil && err.Error() != "No expired members retrieved" {
		log.Panic(err)
	}
	if member != "" {
		msg, errIn := red.get(queueID + "." + member)
		if errIn == nil {
			taskMessage.fromString(msg)
			return taskMessage
		}
		if errIn.Error() == "Nothing found at key" {
			log.Println(errIn)
		}
	}

	// try if there is a task on the received queue to pick up this
	// rare condition would only happen if there would be a failure after
	// the next brpoplpush, but before the popQueueAndSaveKeyToSet
	msg, err := red.rpop(receivedList)
	if err == nil {
		taskMessage.fromString(msg)
		return taskMessage
	}

	// block; wait and switch a task from the queue to the received queue
	red.brpoplpush(queueID, receivedList)

	// fetch from received queue, save the message to it's key
	// and add the ID to the running set
	msg, err = red.popQueueAndSaveKeyToSet(queueID, expirationSec)
	if err != nil {
		// retry.
		log.Println("Didn't find an item on the received queue (exception), will try to pop a new one from incoming queue")
		red.Load(queueID, expirationSec)
	}

	taskMessage.fromString(msg)

	log.Println("LOAD END")
	log.Println("**********")

	return taskMessage
}

// Heartbeat updates the status of a message
func (red *RedClient) Heartbeat(queueID string, taskID string, expirationSec int) bool {

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

	log.Println("***************")
	log.Println("COMPLETE START")

	from := fmt.Sprintf("%s.running", queueID)
	to := fmt.Sprintf("%s.complete", queueID)
	member := taskID

	count, _ := red.moveMemberFromSetToSet(from, to, member)
	if count != 1 {
		log.Printf("Didn't complete the right amount of items!: %d", count)
		return false
	}

	log.Println("COMPLETE END")
	log.Println("***************")
	return true
}

//
// func main() {
//
// 	red := RedClient{host: "localhost:6379"}
// 	_ = red.init()
//
// 	red.Load("queue", 300)
//
// 	fmt.Println("DEBUGGING")
//
// 	// _ = red.set(taskId, taskMessage)
//
// 	// fmt.Println(taskMessage.body)
//
// }
