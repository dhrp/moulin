package main

import (
	"encoding/json"
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
func GenerateMessage(body string) string {

	timestamp := int64(time.Now().UnixNano())
	taskMessage := TaskMessage{
		ID:   timestamp,
		Body: body,
		Envs: []string{"FOO=BAR", "CHUNK=1"},
	}
	b, err := json.Marshal(taskMessage)
	if err != nil {
		log.Panic("error:", err)
	}
	return string(b)
}

func (suite *RedClientTestSuite) TestPushAndPopQueue() {

	task := GenerateMessage(suite.sampleMsgBody)

	suite.NotEqual("{}", task, "Err: json unmarshalled empty!")

	// check if item is set
	var listLength, _ = suite.red.lpush("test.queue", task)
	suite.Equal(listLength, 1, "Expected exactly 1 item in queue")

	// check if the same item is retrieved
	resp := suite.red.brpop("test.queue")
	suite.Equal(resp, task)
}

func (suite *RedClientTestSuite) TestUnpackAndStore() {

	// first store a message in the queue
	task := GenerateMessage(suite.sampleMsgBody)
	// _ = suite.red.lpush("test.queue", task)

	// ToDo: Here it's time to decode the json
	// feth the key and store it.
	// https://golang.org/pkg/encoding/json/#Decoder

	// validJson := json.V
	var taskMessage TaskMessage
	err := json.Unmarshal([]byte(task), &taskMessage)
	if err != nil {
		panic(err)
	}

	log.Println("unpacked taskMesage body: " + taskMessage.Body)
	suite.Equal(suite.sampleMsgBody, taskMessage.Body)
}

func (suite *RedClientTestSuite) TestSortedSet() {
	GenerateMessage("testing the heartbeat")

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

}

// complete sets

func (suite *RedClientTestSuite) TestLoadPhase() {
	log.Println("LOAD phase: testing entire 'load' phase")

	// prepare one item on the queue
	task := GenerateMessage(suite.sampleMsgBody)
	suite.red.lpush("test.load-queue", task)

	// load it
	result := Load(suite.red, "test.load-queue")
	suite.True(result, "Err: the load process returned false")

}

func (suite *RedClientTestSuite) TestHeartbeatPhase() {
	log.Println("HEARTBEAT phase")
	GenerateMessage("testing the heartbeat")

	queueID := "test.heartbeats"
	taskID := "123123123123"

	set := queueID + ".running"
	score := "999"
	value := queueID + "." + taskID

	// set the original count
	count, _ := suite.red.zadd(set, score, value)
	suite.Equal(1, count, "Failed to zadd an item, did we not start clean?")

	score = newScore()

	result := Heartbeat(suite.red, queueID, taskID, score)
	suite.True(result, "Heartbeat didn't return 1, item could not be updated")

	result = Heartbeat(suite.red, queueID, taskID, score)
	suite.False(result, "Heartbeat should have returned 1 because we updated it with the same score")
}

func (suite *RedClientTestSuite) TestCompletePhase() {
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

// The TearDownTest method will be run after every test in the suite.
func (suite *RedClientTestSuite) TearDownTest() {
	suite.red.del("test.queue")
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
