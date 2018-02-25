package rouge

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/segmentio/ksuid"
)

// Client Basic client class to group Redis functions
type Client struct {
	Host       string
	clientpool *pool.Pool
}

// Init Initializes the Rouge.Client, and saves it to the client struct
func (red *Client) Init() error {

	// df := func(network, addr string) (*redis.Client, error) {
	// 	client, err := redis.Dial(network, addr)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	// TODO: Review if we need a password
	// 	// set password with CONFIG SET requirepass "nevermind"
	// 	// if err = client.Cmd("AUTH", "nevermind").Err; err != nil {
	// 	// 	client.Close()
	// 	// 	return nil, err
	// 	// }
	// 	return client, nil
	// }

	// pool, err := pool.NewCustom("tcp", red.Host, 10, df)
	pool, err := pool.New("tcp", red.Host, 10)
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

//
// func reportDone(ctx context.Context) {
// 	done := <-ctx.Done()
// 	fmt.Println("Context done!!", done)
// }

// SlimLoad loads a message from the queue in a lean way. For debugging
func (red *Client) SlimLoad(ctx context.Context, queueID string, expirationSec int) (TaskMessage, error) {

	msg := red.brpop(ctx, queueID)

	var taskMessage TaskMessage
	// taskMessage := TaskMessage{ID: "myid", Body: "echo mybody"}
	taskMessage.FromString(msg)

	return taskMessage, nil
}

// Load loads a message from the queue, and mark it as processing
func (red *Client) Load(ctx context.Context, queueID string, expirationSec int) (TaskMessage, error) {

	for {

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
				taskMessage.FromString(msg)
				log.Println("LOAD END:  Returning expired member")
				log.Println("**********")
				return taskMessage, nil
			}
			if errIn.Error() == "Nothing found at key" {
				log.Println(errIn)
				return TaskMessage{}, errIn
			} else {
				return TaskMessage{}, errIn
			}
		}

		// try if there is a task on the received queue to pick up this
		// rare condition would only happen if there would be a failure after
		// the next brpoplpush, but before the popQueueAndSaveKeyToSet
		// rpopmsg, err := red.rpop(receivedList)
		// ToDo: Review what we do with using this same function twice, with different variable names
		rcvdmsg, err := red.popQueueAndSaveKeyToSet(queueID, receivedList, runningSet, expirationSec)
		if err == nil {
			taskMessage.FromString(rcvdmsg)
			log.Println("LOAD END:  Returning lost member (member on received queue)")
			log.Println("**********")
			return taskMessage, nil
		} else {
			// let go, this is ok
		}

		waitChan := make(chan error)

		go func(waitChan chan error) {
			// block; wait and switch a task from the queue to the received queue
			_, err := red.brpoplpush(queueID, receivedList)
			log.Println("Received message from the queue..")
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered from disconnected channel error. Going to move item back to incoming")
					_, errin := red.lpoprpush(receivedList, queueID)
					if err != nil {
						log.Panic(errin)
					}
				}
			}()
			waitChan <- err
		}(waitChan)

		select {
		case err := <-waitChan:
			if err != nil {
				return TaskMessage{}, errors.Wrap(err, "failed brpoplpush")
			}
			log.Println("ok, got msg")
		case <-ctx.Done():
			log.Println("cancelling blocking pop")
			close(waitChan)
			return TaskMessage{}, errors.New("the context was cancelled")
		}

		// fetch from received queue, save the message to it's key
		// and add the ID to the running set
		setmsg, err := red.popQueueAndSaveKeyToSet(queueID, receivedList, runningSet, expirationSec)
		if err != nil {
			// retry.
			log.Println("Didn't find an item on the received queue (exception), " +
				"will try to pop a new one from incoming queue")
			continue
		}
		taskMessage.FromString(setmsg)

		if taskMessage.Body == "" {
			log.Printf("Error: msg was %s", setmsg)
			log.Panic("Body is empty!")
		}

		log.Println("LOAD END:  Returning new member from the incoming queue")
		log.Println("**********")

		return taskMessage, nil
	}
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
	to := fmt.Sprintf("%s.completed", queueID)
	member := taskID

	ok, err := red.moveMemberFromSetToSet(from, to, member)
	if err != nil {
		return false, errors.Wrap(err, "couldn't complete the item")
	}

	log.Println("COMPLETE END")
	log.Println("***************")
	return ok, nil
}

// Progress gets the status of the current lists in the queue
func (red *Client) Progress(queueID string) (QueueInfo, error) {
	log.Println("***************")
	log.Println("PROGRESS START")

	var queueInfo QueueInfo
	var err error

	// show length of incoming list
	queueInfo.incomingCount, err = red.getListLength(queueID)
	if err != nil {
		log.Panic(err)
	}

	// show length of received list
	queueInfo.receivedCount, err = red.getListLength(queueID + ".received")
	if err != nil {
		log.Panic(err)
	}

	// show count of non-expired items in running set
	now := time.Now().Unix()
	score := strconv.FormatInt(now, 10)
	queueInfo.runningCount, _ = red.zcount(queueID+".running", score, "inf")

	// show count of expired items in running set
	queueInfo.expiredCount, _ = red.zcount(queueID+".running", "-inf", score)

	// show count of items in completed set
	queueInfo.completedCount, _ = red.zcount(queueID+".completed", "-inf", "inf")

	// show count of items in failed set
	queueInfo.failedCount, _ = red.zcount(queueID+".failed", "-inf", "inf")

	log.Println(queueInfo.ToString())
	log.Println("PROGRESS END")
	log.Println("***************")
	return queueInfo, nil
}

// Peek gets a list of the most or least recent items from a given queue
// and queue phase
func (red *Client) Peek(queueID, phase string, limit int) (int, []TaskMessage, error) {
	log.Println("***************")
	log.Println("PEEK START")
	var rawList, members []string
	var taskList []TaskMessage
	var err error
	var item string
	var count int

	if phase == "incoming" {
		from := 0 - limit // set is inclusive
		to := -1          // this is the last item
		count, _ = red.getListLength(queueID)
		rawList, err = red.lrange(queueID, from, to)
		if err != nil {
			return count, taskList, errors.Wrap(err, "couldn't lrange incoming list")
		}
	} else {
		if phase == "running" {
			timestamp := int64(time.Now().Unix())
			score := strconv.FormatInt(timestamp, 10)
			count, _ = red.zcount(queueID+".running", score, "inf")
			members, err = red.zrangebyscore(queueID+".running", score, "inf", limit)
		} else if phase == "expired" {
			timestamp := int64(time.Now().Unix())
			score := strconv.FormatInt(timestamp, 10)
			count, _ = red.zcount(queueID+".running", "-inf", score)
			members, err = red.zrangebyscore(queueID+".running", "-inf", score, limit)
		} else {
			return 0, nil, errors.New("phase not found")
		}
		if err != nil {
			return count, taskList, errors.Wrap(err, "failed to peek")
		}

		for i := 0; i < len(members); i++ {
			item, err = red.get(queueID + "." + members[i])
			rawList = append(rawList, item)
		}
		if err != nil {
			return count, taskList, errors.Wrap(err, "failed getting one of the items from store")
		}
	}

	for i := 0; i < len(rawList); i++ {
		task := TaskMessage{}
		task.FromString(rawList[i])
		taskList = append(taskList, task)
	}

	log.Println("PEEK END")
	log.Println("***************")
	return count, taskList, nil
}

// AddTask adds a new task to a queue.
func (red *Client) AddTask(queueID string, task TaskMessage) (int, error) {
	log.Println("***************")
	log.Println("ADDTASK START")

	if red.clientpool == nil {
		return -1, errors.New("Connection to Redis not initialized. Did you forget to initialize?")
	}

	taskMessageStr := task.ToString()
	newlength, err := red.lpush(queueID, taskMessageStr)
	if err != nil {
		log.Panic(err)
	}

	log.Println("ADDTASK END")
	log.Println("***************")
	return newlength, nil
}

// ClearQueue deletes everything belonging to a given queue.
func (red *Client) ClearQueue(queueID string) (bool, error) {
	log.Println("***************")
	log.Println("CLEARQUEUE START")

	red.del(queueID)
	red.del(queueID + ".running")
	red.del(queueID + ".expired")
	red.del(queueID + ".completed")
	// ToDo: !! Clear the individual keys

	log.Println("CLEARQUEUE END")
	log.Println("***************")
	return true, nil
}

// AddTaskFromString converts a plain string to a task and puts it on the queue
func (red *Client) AddTaskFromString(queueID string, message string) (int, error) {
	log.Println("***************")
	log.Println("ADDTASKSFROMSTRING START")

	id := ksuid.New().String()
	task := TaskMessage{ID: id, Body: message}
	len, err := red.AddTask(queueID, task)
	if err != nil {
		return -1, errors.Wrap(err, "failed to add tasks from string")
	}
	log.Println("ADDTASKSFROMSTRING END")
	log.Println("***************")
	return len, nil
}

// AddTasksFromFile is a function for loading from a file
func (red *Client) AddTasksFromFile(queueID, filePath string) (queueLength int, count int, err error) {
	log.Println("***************")
	log.Println("ADDTASKSFROMFILE START")
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
	log.Println("ADDTASKSFROMFILE END")
	log.Println("***************")
	return queueLength, count, nil
}
