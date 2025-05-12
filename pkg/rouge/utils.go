package rouge

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	pb "github.com/dhrp/moulin/pkg/protobuf"
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

// QueueProgress is a type meant for representing progress
type QueueProgress struct {
	incomingCount  int
	receivedCount  int
	runningCount   int
	expiredCount   int
	completedCount int
	failedCount    int
}

// ToBuff converts a QueueProgress object to it's protobuf representation
func (obj *QueueProgress) ToBuff() *pb.QueueProgress {
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

// ToString converts a QueueInfo object to string
func (obj *QueueProgress) ToString() string {

	stringFmt := `
incomingCount  %d
`

	return fmt.Sprintf(stringFmt,
		obj.incomingCount,
	)
}

// QueueInfo is a type meant for representing info about a queue
type QueueInfo struct {
	QueueID  string
	Progress QueueProgress
}

// ToString converts a QueueInfo object to string
func (obj *QueueInfo) ToString() string {

	stringFmt := `
name 		 %s
incomingCount  %d
`

	return fmt.Sprintf(stringFmt,
		obj.QueueID,
		obj.Progress.incomingCount,
	)
}

// ToBuff converts a QueueInfo object to it's protobuf representation
func (obj *QueueInfo) ToBuff() (queueInfo *pb.QueueInfo) {

	name := obj.QueueID
	qp := obj.Progress.ToBuff()

	return &pb.QueueInfo{
		QueueID:  name,
		Progress: qp,
	}
}

func newScore() string {
	// determine the score (when task should expire)
	timestamp := int64(time.Now().Unix())
	expires := timestamp + 300 // 5 min
	return strconv.FormatInt(expires, 10)
}
