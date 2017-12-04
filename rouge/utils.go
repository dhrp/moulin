package rouge

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	pb "github.com/nerdalize/moulin/protobuf"
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
	incomingListLength int
	receivedListLength int
	nonExpiredCount    int
	expiredCount       int
	completedCount     int
	failedCount        int
	runningItems       []TaskMessage
}

func (obj *QueueInfo) ToString() string {

	stringFmt := `
incomingListLength %d
receivedListLength %d
nonExpiredCount    %d
expiredCount       %d
completedCount     %d
failedCount        %d
runningItems       "some"
	`
	_ = stringFmt

	return fmt.Sprintf(stringFmt,
		obj.incomingListLength,
		obj.receivedListLength,
		obj.nonExpiredCount,
		obj.expiredCount,
		obj.completedCount,
		obj.failedCount)
}

// ToBuff converts a QueueInfo object to it's protobuf representation
func (obj *QueueInfo) ToBuff() *pb.QueueProgress {
	qp := &pb.QueueProgress{
		IncomingListLength: int32(obj.incomingListLength),
		ReceivedListLength: int32(obj.receivedListLength),
		NonExpiredCount:    int32(obj.nonExpiredCount),
		ExpiredCount:       int32(obj.expiredCount),
		CompletedCount:     int32(obj.completedCount),
		FailedCount:        int32(obj.failedCount),
		// runningItems       "some"
		// expiredItems       "some"
	}
	return qp
}

func newScore() string {
	// determine the score (when task should expire)
	timestamp := int64(time.Now().Unix())
	expires := timestamp + 300 // 5 min
	return strconv.FormatInt(expires, 10)
}
