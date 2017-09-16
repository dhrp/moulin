package kafkaconsumer

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type KafkaConsumerTestSuite struct {
	suite.Suite
}

func (suite *KafkaConsumerTestSuite) TestProducer() {

	topics := []string{"test", "foo"}
	consume("localhost", "main", topics)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestKafkaConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(KafkaConsumerTestSuite))
}
