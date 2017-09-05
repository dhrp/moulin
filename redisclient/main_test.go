package main

import (
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type RedClientTestSuite struct {
	suite.Suite
	red           RedClient
	sampleMsgBody string
}

func (suite *RedClientTestSuite) SetupSuite() {
	suite.sampleMsgBody = "http://www.peskens.nl"

	suite.red = RedClient{host: "localhost:6379"}
	_ = suite.red.init()
}

// GenerateMessage generate a json message for piping through the system.
// most notably we add the /current/ timestamp.
// func (suite *RedClientTestSuite) GenerateMessage() string {
func GenerateMessage(body string) TaskMessage {

	timestamp := int64(time.Now().UnixNano())
	taskMessage := TaskMessage{
		ID:   timestamp,
		Body: body,
		Envs: []string{"FOO=BAR", "CHUNK=1"},
	}
	return taskMessage
}

func (suite *RedClientTestSuite) TestpopQueueAndSaveKeyToSet() {

	// prepare one item on the queue
	taskMessage := GenerateMessage(suite.sampleMsgBody)
	suite.red.lpush("test.queue", taskMessage.toString())

	queueID := "test.queue"
	destinationSet := "queue.running"
	expirationSec := 333

	msg, _ := suite.red.popQueueAndSaveKeyToSet(queueID, destinationSet, expirationSec)
	suite.Equal(taskMessage.toString(), msg, "message from popQueueAndSaveKeyToSet didn't match expectation")
}

func (suite *RedClientTestSuite) TestRPOP() {
	log.Println("*** testing TestRPOP")

	taskMessage := GenerateMessage(suite.sampleMsgBody)
	suite.red.lpush("test.queue.received", taskMessage.toString())

	msg, _ := suite.red.rpop("test.queue.received")

	suite.Equal(taskMessage.toString(), msg, "what was pushed is not what was popped")

	msg, _ = suite.red.rpop("test.queue.doesntexist")
	suite.Equal(msg, "", "what was pushed is not what was popped")

}

func (suite *RedClientTestSuite) TestPushAndPopQueue() {

	taskMessage := GenerateMessage(suite.sampleMsgBody)
	taskMessageStr := taskMessage.toString()

	suite.NotEqual("{}", taskMessageStr, "Err: json unmarshalled empty!")

	// set an item
	var listLength, _ = suite.red.lpush("test.queue", taskMessageStr)
	suite.Equal(listLength, 1, "Expected exactly 1 item in queue")

	// check if the same item is retrieved
	resp := suite.red.brpop("test.queue")
	suite.Equal(resp, taskMessageStr)

	// set an item
	listLength, _ = suite.red.lpush("test.queue", taskMessageStr)
	suite.Equal(listLength, 1, "Expected exactly 1 item in queue")

	// check if the same item is retrieved
	resp2 := suite.red.brpoplpush("test.queue", "test.queue.received")
	suite.Equal(resp2, taskMessageStr)
}

func (suite *RedClientTestSuite) TestSortedSet() {

	log.Println("### Testing Sorted Set Start")

	queueID := "test.sorted_sets"
	taskID := "123123123123"

	set := queueID + ".running"
	score := newScore()
	value := queueID + "." + taskID

	// set the original count
	count, _ := suite.red.zadd(set, score, value)
	suite.Equal(1, count, "Failed to item to set an item, did we not start clean?")

	// // prepare an expired score
	timestamp := int64(time.Now().Unix()) - 100
	expiredScore := strconv.FormatInt(timestamp, 10)

	// update the count. This should return one. We set the count to a previous date
	count, _ = suite.red.zaddUpdate(set, expiredScore, value)
	suite.Equal(1, count, "An update should have worked")

	// try to update a non-existing item. This should fail (return zero)
	count, _ = suite.red.zaddUpdate(set, score, "nonexistent")
	suite.Equal(0, count, "We should not have set new values")

	log.Println("### Testing Sorted Set End")

}

// complete sets

func (suite *RedClientTestSuite) TestLoadPhase() {
	log.Println("**********")
	log.Println("LOAD PHASE TEST")

	// prepare one item on the queue
	taskMessage1 := GenerateMessage(suite.sampleMsgBody)
	suite.red.lpush("test.queue", taskMessage1.toString())

	taskMessage2 := GenerateMessage(suite.sampleMsgBody)
	suite.red.lpush("test.queue", taskMessage2.toString())

	// load it back
	result := Load(suite.red, "test.queue")
	suite.Equal(taskMessage1, result, "The message put on the queue is not what came back")

	// test if it is now also in the sorted set.
	members, _ := suite.red.zrevrange("test.queue.running", 0, 0) // set in inclusive
	suite.Equal(1, len(members), "expected exactly one item to be returned from set")

	if len(members) == 1 {
		expectedID := strconv.FormatInt(result.ID, 10)
		suite.Equal(expectedID, members[0], "the taskID doesn't match the one in the running set")
	}

	log.Println("**********")
	log.Println("END LOAD PHASE TEST")

}

func (suite *RedClientTestSuite) TestHeartbeatPhase() {
	log.Println("***************")
	log.Println("HEARTBEAT phase")

	queueID := "test.heartbeats"
	taskID := "123123123123"

	set := queueID + ".running"
	score := "999"
	member := taskID

	// set the original count
	count, _ := suite.red.zadd(set, score, member)
	suite.Equal(1, count, "Failed to zadd an item, did we not start clean?")

	expirationSec := 300 // 5 min

	result := Heartbeat(suite.red, queueID, member, expirationSec)
	suite.True(result, "Heartbeat didn't return 1, item could not be updated")

	result = Heartbeat(suite.red, queueID, member, expirationSec)
	suite.False(result, "Heartbeat should have returned false because we updated it with the same score")
}

func (suite *RedClientTestSuite) TestCompletePhase() {
	log.Println("**************")
	log.Println("COMPLETE phase")

	set := "test.running"
	score := "100"
	member := "some_queue.some_task"

	suite.red.zadd(set, score, member)

	from := set
	to := "test.completed"

	count, _ := suite.red.moveMemberFromSetToSet(from, to, member)
	suite.Equal(1, count, "didn't manage to move item from one set to other")
}

// func (suite *RedClientTestSuite) TestRedEndToEnd() {
//
// 	log.Println("*******************")
// 	log.Println("END TO END TEST")
// 	queueID := "test.queue"
//
// 	// Prep an item on the queue
// 	taskMessage := GenerateMessage("testing end to end")
// 	suite.red.lpush(queueID, taskMessage.toString())
//
// 	// Load it from the queue
// 	returnedTaskMessage := Load(suite.red, queueID)
//
// 	print(returnedTaskMessage)
//
// 	taskID := strconv.FormatInt(returnedTaskMessage.ID, 10)
//
// 	// Send a heartbeat
// 	status := Heartbeat(suite.red, queueID, taskID, 400)
// 	suite.True(status, "Heartbeat failed to update the item...")
//
// 	// Mark the item as complete
// 	complete := Complete(suite.red, queueID, taskID)
// 	suite.True(complete, "Complete failed to complete the item...")
//
// 	log.Println("END TO END TEST END")
// 	log.Println("*******************")
// }

// The TearDownTest method will be run after every test in the suite.
func (suite *RedClientTestSuite) TearDownTest() {
	suite.red.del("test.queue")
	suite.red.del("test.queue.running")
	suite.red.del("test.queue.received")

	suite.red.del("test.sorted_sets.running")
	suite.red.del("test.heartbeats.running")
	suite.red.del("test.running")
	suite.red.del("test.completed")
}

func (suite *RedClientTestSuite) TearDownSuite() {
	log.Println("closing suite, cleaning up Redis")
	suite.red.del("test.queue")
	suite.red.del("test.sorted_sets.running")
	suite.red.del("test.heartbeats.running")

}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRedClientTestSuite(t *testing.T) {
	suite.Run(t, new(RedClientTestSuite))
}
