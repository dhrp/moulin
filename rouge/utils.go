package rouge

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	pb "github.com/dhrp/moulin/protobuf"
)

// TaskMessage is the 'message' format used on our queue
type TaskMessage struct {
	ID   string   `json:"id"`
	Body string   `json:"body"`
	Envs []string `json:"envs"`
}

// ToString converts a taskMessage to a string
func (t *TaskMessage) ToString() string {
	b, err := json.Marshal(t)
	if err != nil {
		log.Panic("error:", err)
	}
	return string(b)
}

// FromString makes a TaskMessage from a string
func (t *TaskMessage) FromString(ts string) (*TaskMessage, error) {
	tb := []byte(ts)
	err := json.Unmarshal(tb, &t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// QueueInfo is a type meant for representing progress
type QueueInfo struct {
	incomingCount  int
	receivedCount  int
	runningCount   int
	expiredCount   int
	completedCount int
	failedCount    int
}

// ToString converts a QueueInfo object to string
func (obj *QueueInfo) ToString() string {

	stringFmt := `
incomingCount  %d
receivedCount  %d
runningCount   %d
expiredCount   %d
completedCount %d
failedCount    %d`

	return fmt.Sprintf(stringFmt,
		obj.incomingCount,
		obj.receivedCount,
		obj.runningCount,
		obj.expiredCount,
		obj.completedCount,
		obj.failedCount)
}

// ToBuff converts a QueueInfo object to it's protobuf representation
func (obj *QueueInfo) ToBuff() *pb.QueueProgress {
	qp := &pb.QueueProgress{
		IncomingCount:  int32(obj.incomingCount),
		ReceivedCount:  int32(obj.receivedCount),
		RunningCount:   int32(obj.runningCount),
		ExpiredCount:   int32(obj.expiredCount),
		CompletedCount: int32(obj.completedCount),
		FailedCount:    int32(obj.failedCount),
	}
	return qp
}

func newScore() string {
	// determine the score (when task should expire)
	timestamp := int64(time.Now().Unix())
	expires := timestamp + 300 // 5 min
	return strconv.FormatInt(expires, 10)
}
