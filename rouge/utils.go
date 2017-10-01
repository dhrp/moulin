package rouge


import (
	"encoding/json"
	"log"
	"strconv"
	"time"
)

// TaskMessage is the 'message' format used on our queue
type TaskMessage struct {
	ID   string   `json:"id"`
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
