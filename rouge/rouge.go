package rouge

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/segmentio/ksuid"
)

// Client Basic client class to group Redis functions
type Client struct {
	Host       string
	clientpool *pool.Pool
}

// Init Initializes the Rouge.Client, and saves it to the client struct
func (red *Client) Init() error {

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

	pool, err := pool.NewCustom("tcp", red.Host, 10, df)
	if err != nil {
		return errors.Wrap(err, "init new connection to redis failed")
	}

	log.Println("redis client connected successfully with radix driver")
	red.clientpool = pool

	return nil
}

// Info just returns some information about the redis server. We also use it as
// a health check
func (red *Client) Info() (string, error) {

	if red.clientpool == nil {
		log.Fatal("Connection to Redis not initialized. Did you forget to initialize?")
	}

	resp := red.clientpool.Cmd("INFO", "server")
	if err := resp.Err; err != nil {
		log.Panic(err)
	}
	return resp.Str()
}

// Load loads a message from the queue, and mark it as processing
func (red *Client) Load(queueID string, expirationSec int) TaskMessage {

	if red.clientpool == nil {
		log.Fatal("Connection to Redis not initialized. Did you forget to initialize?")
	}

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
			taskMessage.FromString([]byte(msg))
			log.Println("LOAD END:  Returning expired member")
			log.Println("**********")
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
		taskMessage.FromString([]byte(msg))
		log.Println("LOAD END:  Returning lost member (member on received queue)")
		log.Println("**********")
		return taskMessage
	}

	// block; wait and switch a task from the queue to the received queue
	debugmsg := red.brpoplpush(queueID, receivedList)
	log.Println(debugmsg)

	// fetch from received queue, save the message to it's key
	// and add the ID to the running set
	msg, err = red.popQueueAndSaveKeyToSet(queueID, expirationSec)
	if err != nil {
		// retry.
		log.Println("Didn't find an item on the received queue (exception), will try to pop a new one from incoming queue")
		red.Load(queueID, expirationSec)
	}

	taskMessage.FromString([]byte(msg))

	log.Println("LOAD END:  Returning new member from the incoming queue")
	log.Println("**********")

	return taskMessage
}

// Heartbeat updates the status of a message
func (red *Client) Heartbeat(queueID string, taskID string, expirationSec int32) (int, error) {

	if red.clientpool == nil {
		log.Fatal("Connection to Redis not initialized. Did you forget to initialize?")
	}

	log.Println("***************")
	log.Println("HEARTBEAT START")

	set := queueID + ".running"
	member := taskID

	expiresAt := time.Now().Unix() + int64(expirationSec)
	score := strconv.FormatInt(expiresAt, 10)

	// _, _, _ = set, score, value
	count, _ := red.zaddUpdate(set, score, member)
	if count == 0 {
		errMsg := "Heartbeat: no item could be found to update, was item already " +
			"completed?, or perhaps the item was updated < 1 sec ago."
		log.Println(errMsg)
		return 0, errors.New(errMsg)
	}

	log.Println("HEARTBEAT END")
	log.Println("***************")

	return int(expiresAt), nil
}

// Complete marks the item as completed.
func (red *Client) Complete(queueID string, taskID string) (bool, error) {

	if red.clientpool == nil {
		log.Fatal("Connection to Redis not initialized. Did you forget to initialize?")
	}

	log.Println("***************")
	log.Println("COMPLETE START")

	from := fmt.Sprintf("%s.running", queueID)
	to := fmt.Sprintf("%s.complete", queueID)
	member := taskID

	ok, err := red.moveMemberFromSetToSet(from, to, member)
	if err != nil {
		return false, errors.Wrap(err, "couldn't complete any items")
	}

	log.Println("COMPLETE END")
	log.Println("***************")
	return ok, nil
}

func (red *Client) GetProgress(queueID string) (string, error) {
	// show length of incoming list
	// show length of received list

	// show count of items in data store, and the total size consumed

	// show count of non-expired items in running set
	// show count of expired items in running set
	// show count of items in completed set
	// show count of items in failed set

	// show a list of all items (with content) of now working on and failed
	//  * which worker is working on it
	//  * command and arguments

	// get size of all keys for this queue combined

	answer := "foobar"
	return answer, nil
}

// AddTask adds a new task to a queue.
func (red *Client) AddTask(queueID string, task TaskMessage) (int, error) {

	if red.clientpool == nil {
		return -1, errors.New("Connection to Redis not initialized. Did you forget to initialize?")
	}

	taskMessageStr := task.ToString()
	newlength, err := red.lpush(queueID, taskMessageStr)
	if err != nil {
		log.Panic(err)
	}

	return newlength, nil
}

// AddTaskFromString converts a plain string to a task and puts it on the queue
func (red *Client) AddTaskFromString(queueID string, message string) (int, error) {

	id := ksuid.New().String()
	task := TaskMessage{ID: id, Body: message}
	len, err := red.AddTask(queueID, task)
	if err != nil {
		return -1, errors.Wrap(err, "failed to add tasks from string")
	}
	return len, nil
}

// AddTasksFromFile is a function for loading from a file
func (red *Client) AddTasksFromFile(queueID, filePath string) (queueLength int, count int, err error) {
	if red == nil {
		log.Panic("rouge not initialized in AddTasksFromFile")
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	// var line string
	scanner := bufio.NewScanner(file)

	i := 0
	var values [][]byte

	for scanner.Scan() {
		i++
		values = append(values, scanner.Bytes())
		if i == 1000 {
			// finalize
		}
	}
	count = i

	// check if scanner has any errors
	// ToDo: What if?
	if err = scanner.Err(); err != nil {
		log.Panic(err)
	}

	// send n messages
	for _, value := range values {

		queueLength, err = red.AddTaskFromString(queueID, string(value))
		if err != nil {
			return queueLength, count, errors.Wrap(err, "failed AddTasksFromFile")
		}
	}

	log.Printf("sucessfully added %d items, queue now %d items long", count, queueLength)
	return queueLength, count, nil
}
