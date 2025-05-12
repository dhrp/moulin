package rouge

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/suite"
)

type RedClientTestSuite struct {
	suite.Suite
	red           *Client
	sampleMsgBody string
}

func (suite *RedClientTestSuite) SetupSuite() {
	suite.sampleMsgBody = "http://www.peskens.nl"

	client, err := NewRougeClient()
	suite.Nil(err)

	suite.red = client
}

// GenerateMessage generate a json message for piping through the system.
// most notably we add the /current/ timestamp.
// func (suite *RedClientTestSuite) GenerateMessage() string {
func GenerateMessage(body string) TaskMessage {

	randomID := ksuid.New().String()
	randomBytes := ksuid.New()
	_ = randomBytes

	// timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	taskMessage := TaskMessage{
		ID:   randomID,
		Body: body,
		Envs: []string{"FOO=BAR", "CHUNK=1"},
	}
	return taskMessage
}

func (suite *RedClientTestSuite) TestBasicConnection() {
	info, err := suite.red.Info()
	suite.Nil(err, "error, is redis running?")
	suite.NotEmpty(info, "expected redis info")
}

func (suite *RedClientTestSuite) TestPopQueueAndSaveKeyToSet() {
	// prepare one item on the queue
	taskMessage := GenerateMessage(suite.sampleMsgBody)
	taskMessageStr := taskMessage.ToString()
	suite.red.lpush("test.queue.received", taskMessageStr)

	queueID := "test.queue"
	expirationSec := 333

	receivedList := queueID + ".received"
	runningSet := queueID + ".running"

	// positive case; item exists
	msg, err := suite.red.popQueueAndSaveKeyToSet(queueID, receivedList, runningSet, expirationSec)
	suite.Equal(taskMessage.ToString(), msg, "message from popQueueAndSaveKeyToSet didn't match expectation")
	suite.Nil(err, "should have thrown an error when no items in queue")

	// negative case; no item on received queue
	_, err = suite.red.popQueueAndSaveKeyToSet("empty", receivedList, runningSet, expirationSec)
	suite.NotNil(err, "should have thrown an error when no items in queue")
}

func (suite *RedClientTestSuite) TestRPOP() {
	log.Println("*** testing TestRPOP")

	taskMessage := GenerateMessage(suite.sampleMsgBody)
	suite.red.lpush("test.queue.received", taskMessage.ToString())

	msg, _ := suite.red.rpop("test.queue.received")
	suite.Equal(taskMessage.ToString(), msg, "what was pushed is not what was popped")

	msg, _ = suite.red.rpop("test.queue.doesntexist")
	suite.Equal(msg, "", "what was pushed is not what was popped")
}

func (suite *RedClientTestSuite) TestLPOPRPUSH() {

	log.Println("*** testing LPOPRPUSH")

	taskMessage := GenerateMessage(suite.sampleMsgBody)
	suite.red.lpush("test.queue.received", taskMessage.ToString())

	count, err := suite.red.lpoprpush("test.queue.received", "test.queue")
	suite.Equal(count, 1)
	suite.Nil(err)

	msg, _ := suite.red.rpop("test.queue")
	suite.Equal(msg, taskMessage.ToString(), "what was pushed is not what was popped")
}

func (suite *RedClientTestSuite) TestPushAndPopQueue() {

	taskMessage := GenerateMessage(suite.sampleMsgBody)
	taskMessageStr := taskMessage.ToString()

	suite.NotEqual("{}", taskMessageStr, "Err: json unmarshalled empty!")

	// set an item
	var listLength, _ = suite.red.lpush("test.queue", taskMessageStr)
	suite.Equal(listLength, 1, "Expected exactly 1 item in queue")

	// check if the same item is retrieved

	ctx := context.TODO()
	resp := suite.red.brpop(ctx, "test.queue")
	suite.Equal(resp, taskMessageStr)

	// set an item
	listLength, _ = suite.red.lpush("test.queue", taskMessageStr)
	suite.Equal(listLength, 1, "Expected exactly 1 item in queue")

	// check if the same item is retrieved
	resp2, err := suite.red.brpoplpush("test.queue", "test.queue.received")
	suite.Equal(resp2, taskMessageStr)
	suite.Nil(err)
}

func (suite *RedClientTestSuite) TestSortedSet() {

	log.Println("### Testing Sorted Set Start")

	queueID := "test.queue"
	taskID := "123123123123"

	set := queueID + ".running"
	score := newScore()

	// set the original count
	count, _ := suite.red.zadd(set, score, taskID)
	suite.Equal(1, count, "Failed to item to set an item, did we not start clean?")

	// // prepare an expired score
	timestamp := int64(time.Now().Unix()) - 100
	expiredScore := strconv.FormatInt(timestamp, 10)

	// update the count. This should return one. We set the count to a previous date
	count, _ = suite.red.zaddUpdate(set, expiredScore, taskID)
	suite.Equal(1, count, "An update should have worked")

	// try to update a non-existing item. This should fail (return zero)
	count, _ = suite.red.zaddUpdate(set, score, "nonexistent")
	suite.Equal(0, count, "We should not have set new values")

	// Check if we return the most expired item
	timestamp = int64(time.Now().Unix()) - 50
	expiredScore = strconv.FormatInt(timestamp, 10)
	suite.red.zadd(set, expiredScore, "67676767676767")
	// expiredID, _ := suite.red.checkExpired(set)
	expiredID, _ := suite.red.fetchAndUpdateExpired(set, 300)
	suite.Equal(taskID, expiredID, "Did not get the (most) expired ID from the set")

	// Check what happens with no expired items
	// Create item valid 'till 2200
	suite.red.zadd("test.sorted_sets.future", "7258118400", "9090909090909")
	// member, _ := suite.red.checkExpired(set)
	member, _ := suite.red.fetchAndUpdateExpired("test.sorted_sets.future", 300)
	suite.Equal("", member, "Got an item ?!?")

	// Check what happens if the database was not initialized and the set is empty
	suite.red.del("nonexistent")
	_, _ = suite.red.fetchAndUpdateExpired("nonexistent", 300)

	// ToDo: add move from set to set
	ok, err := suite.red.moveMemberFromSetToSet(queueID+".running", queueID+".completed", taskID)
	suite.Nil(err)
	suite.True(ok, "item didn't succeed to move from set to set")

	ok, err = suite.red.moveMemberFromSetToSet(queueID+".running", queueID+".completed", taskID)
	suite.NotNil(err)
	suite.False(ok, "moving the item again was not supposed to have succeeded")

	log.Println("### Testing Sorted Set End")

}

// complete sets
func (suite *RedClientTestSuite) TestLoadPhase() {
	log.Println("**********")
	log.Println("LOAD PHASE TEST")

	ctx := context.TODO()

	// prepare two items on the queue
	taskMessage1 := GenerateMessage(suite.sampleMsgBody)
	suite.red.lpush("test.queue", taskMessage1.ToString())

	taskMessage2 := GenerateMessage(suite.sampleMsgBody)
	suite.red.lpush("test.queue", taskMessage2.ToString())

	// load one back
	result, err := suite.red.Load(ctx, "test.queue", 300)
	suite.Nil(err, "Didn't expect error")
	suite.Equal(taskMessage1, result, "The first message put on the queue is not what came back")

	// test if it is now also in the sorted set.
	members, _ := suite.red.zrevrange("test.queue.running", 0, 0) // set is inclusive
	suite.Equal(1, len(members), "expected exactly one item to be returned from set")

	if len(members) == 1 {
		suite.Equal(result.ID, members[0], "the taskID doesn't match the one in the running set")
	}

	// test the scenario of an expired item in running set
	timestamp := int64(time.Now().Unix()) - 500
	expiredScore := strconv.FormatInt(timestamp, 10)
	// set the score as expired
	updated, _ := suite.red.zadd("test.queue.running", expiredScore, result.ID)
	suite.Equal(0, updated, "A member was added, but not that was not expected")

	expiredTaskMessage, err := suite.red.Load(ctx, "test.queue", 300)
	suite.Nil(err, "Didn't expect error")
	suite.Equal(taskMessage1, expiredTaskMessage, "The what was on the key of the expired member is not what was expected")

	log.Println("**********")
	log.Println("END LOAD PHASE TEST")
}

func (suite *RedClientTestSuite) TestHeartbeatPhase() {
	log.Println("***************")
	log.Println("HEARTBEAT phase")

	queueID := "test.queue"
	taskID := "123123123123"

	set := queueID + ".running"
	score := "999"
	member := taskID

	// set the original count
	count, _ := suite.red.zadd(set, score, member)
	suite.Equal(1, count, "Failed to zadd an item, did we not start clean?")

	var expirationSec int32 = 300 // 5 min

	expires, err := suite.red.Heartbeat(queueID, taskID, expirationSec)
	suite.Nil(err, "didn't expect an error")
	suite.NotZero(expires, "expired a non-zero expiry")

	// Here we try to update it again, which should succeed (again)
	expires2, err := suite.red.Heartbeat(queueID, taskID, expirationSec)
	suite.Nil(err, "didn't expect an error")
	suite.NotZero(expires2, "expired a non-zero expiry")

	// Here we try to update a nonexisting member, which is expected to fail
	expires, err = suite.red.Heartbeat(queueID, "nonexisting", expirationSec)
	suite.NotNil(err, "I expected an error")
	suite.Zero(expires, "expired a non-zero expiry")
}

func (suite *RedClientTestSuite) TestCompletePhase() {
	log.Println("**************")
	log.Println("COMPLETE phase")

	set := "test.queue.running"
	score := "100"
	member := "some_queue.some_task"

	suite.red.zadd(set, score, member)

	from := set
	to := "test.queue.completed"

	OK, err := suite.red.moveMemberFromSetToSet(from, to, member)
	suite.Nil(err)
	suite.Equal(true, OK, "didn't manage to move item from one set to other")
}

func (suite *RedClientTestSuite) TestFailPhase() {
	log.Println("**************")
	log.Println("FAIL phase")

	queueID := "test.queue"
	taskID := "123123123123"
	set := fmt.Sprintf("%s.running", queueID)
	score := "100"

	suite.red.zadd(set, score, taskID)

	OK, err := suite.red.Fail(queueID, taskID)
	suite.Nil(err)
	suite.Equal(true, OK, "didn't manage to move item from one set to other")
}

func (suite *RedClientTestSuite) TestRedEndToEnd() {

	log.Println("*******************")
	log.Println("END TO END TEST")
	queueID := "test.queue"
	ctx := context.TODO()

	// Prep an item on the queue
	taskMessage := GenerateMessage("testing end to end")
	suite.red.lpush(queueID, taskMessage.ToString())

	// Load it from the queue
	msg, err := suite.red.Load(ctx, queueID, 300)
	suite.Nil(err, "Didn't expect error")
	taskID := msg.ID

	// Send a heartbeat
	expires, err := suite.red.Heartbeat(queueID, taskID, 400)
	suite.Nil(err)
	suite.NotZero(expires, "expired a non-zero expiry")

	// Mark the item as complete
	OK, err := suite.red.Complete(queueID, taskID)
	suite.Nil(err, "Complete failed to complete the item...")
	suite.Equal(true, OK, "didn't get the right amount of completed items")

	log.Println("END TO END TEST END")
	log.Println("*******************")
}

func (suite *RedClientTestSuite) TestAddTasksFromFile() {

	log.Println("*******************")
	log.Println("TEST ADD FROM FILE")
	queueID := "test.queue"
	filePath := "../../test/testtextfile.txt"

	queueLength, count, err := suite.red.AddTasksFromFile(queueID, filePath)
	suite.Nil(err, "AddTasksFromFile gave an error")
	suite.Equal(queueLength, 6, "Expected the queue to have this new size")
	suite.Equal(count, 6, "We added 6 items")
}

func (suite *RedClientTestSuite) TestProgress() {

	var err error
	var result QueueProgress
	queueID := "test.queue"
	var msg TaskMessage
	ctx := context.TODO()

	// push three messages
	msg = GenerateMessage("message 1")
	suite.red.lpush(queueID, msg.ToString())
	msg = GenerateMessage("message 2")
	suite.red.lpush(queueID, msg.ToString())
	msg = GenerateMessage("message 3")
	suite.red.lpush(queueID, msg.ToString())

	// Load two items from the queue
	suite.red.Load(ctx, queueID, 300)
	msg, err = suite.red.Load(ctx, queueID, 300)
	suite.Nil(err, "Didn't expect error")

	// Complete one item
	suite.red.Complete(queueID, msg.ID)

	// Now check if we see what we expect
	result, err = suite.red.Progress("test.queue")
	suite.Nil(err, "GetProgress should not give any errors")
	suite.Equal(1, result.incomingCount)
	suite.Equal(1, result.runningCount)
	suite.Equal(1, result.completedCount)
}

func (suite *RedClientTestSuite) TestPeek() {
	var taskList []TaskMessage
	queueID := "test.queue"
	var count int
	ctx := context.TODO()

	suite.red.AddTaskFromString(queueID, "task one")
	suite.red.AddTaskFromString(queueID, "task two")
	suite.red.AddTaskFromString(queueID, "task three")
	members, err := suite.red.lrange(queueID, 0, 30)
	suite.Nil(err)
	suite.Len(members, 3)

	// Check length of incoming
	count, taskList, err = suite.red.Peek(queueID, "incoming", 30)
	suite.Nil(err)
	suite.Len(taskList, 3)
	suite.Equal(3, count)

	// Check length of running
	suite.red.Load(ctx, queueID, 50)
	count, taskList, err = suite.red.Peek(queueID, "running", 30)
	suite.Nil(err)
	suite.Len(taskList, 1)
	suite.Equal(1, count)

	// Check length of expired
	suite.red.Load(ctx, queueID, -50)
	count, taskList, err = suite.red.Peek(queueID, "expired", 30)
	suite.Nil(err)
	suite.Len(taskList, 1)
	suite.Equal(1, count)
}

func (suite *RedClientTestSuite) BeforeTest() {
}

func (suite *RedClientTestSuite) TestListQueues() {

	msg1 := GenerateMessage("message 1")
	msg2 := GenerateMessage("message 2")
	msg3 := GenerateMessage("message 3")

	suite.red.AddTask("c-list", msg1)
	// we sleep just to be sure the timestamp actually increments
	time.Sleep(2 * time.Millisecond)
	suite.red.AddTask("b-list", msg2)
	time.Sleep(2 * time.Millisecond)
	suite.red.AddTask("a-list", msg3)

	list, _ := suite.red.ListQueues("alpha")
	log.Printf("list of queues: %v", list)
	suite.Equal("a-list", list[0].QueueID)

	list, _ = suite.red.ListQueues("created")
	log.Printf("list of queues: %v", list)
	suite.Equal("c-list", list[0].QueueID)

	_, err := suite.red.ListQueues("unknown-type")
	suite.NotNil(err)
}

// test the delete queue function
func (suite *RedClientTestSuite) TestDeleteQueue() {

	var err error
	queueID := "test.queue"
	ctx := context.TODO()

	// push three messages
	msg1 := GenerateMessage("message 1")
	suite.red.AddTask(queueID, msg1)

	msg2 := GenerateMessage("message 2")
	suite.red.AddTask(queueID, msg2)

	msg3 := GenerateMessage("message 3")
	suite.red.AddTask(queueID, msg3)

	// Load two items from the queue
	suite.red.Load(ctx, queueID, 300)
	suite.red.Load(ctx, queueID, 300)

	task, _ := suite.red.get(queueID + "." + msg1.ID)
	suite.Equal(msg1.ToString(), task, "The item should exist")
	suite.red.Complete(queueID, msg1.ID)

	taskCount, err := suite.red.DeleteQueue(queueID)
	suite.Nil(err, "DeleteQueue should not give any errors")
	suite.Equal(3, taskCount, "We should have deleted 3 items")

	// check that the items are gone
	task, err = suite.red.get(queueID + msg1.ID)
	suite.NotNil(err, "The item should not exist")
	task, err = suite.red.get(queueID + msg2.ID)
	suite.NotNil(err, "The item should not exist")
}

// test the delete queue function
func (suite *RedClientTestSuite) TestDeleteQueueNotExisting() {
	_, err := suite.red.DeleteQueue("this.never.existed")
	suite.NotNil(err, "DeleteQueue should give an error")
}

// The TearDownTest method will be run after every test in the suite.
func (suite *RedClientTestSuite) TearDownTest() {
	suite.red.deleteQueue("test.queue")
	suite.red.deleteQueue("test.queue.running")
	suite.red.deleteQueue("test.queue.received")
	suite.red.deleteQueue("test.queue.completed")
	suite.red.deleteQueue("test.queue.expired")
	suite.red.deleteQueue("test.queue.failed")
	suite.red.deleteQueue("nonexistent")
	suite.red.deleteQueue("a-list")
	suite.red.deleteQueue("b-list")
	suite.red.deleteQueue("c-list")
	suite.red.deleteQueue("q1")
	suite.red.deleteQueue("foobar")
	suite.red.deleteQueue("clientTest")
}

func (suite *RedClientTestSuite) TearDownSuite() {
	// suite.red.flushdb()
	log.Println("closing suite, cleaning up Redis")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRedClientTestSuite(t *testing.T) {
	suite.Run(t, new(RedClientTestSuite))
}
