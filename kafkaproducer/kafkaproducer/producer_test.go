package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type KafkaProducerTestSuite struct {
	suite.Suite
	kp KafkaProducer
}

func (suite *KafkaProducerTestSuite) SetupSuite() {

	suite.kp = KafkaProducer{broker: "127.0.0.1:9092"}
	_ = suite.kp.init()

}

func (suite *KafkaProducerTestSuite) TestProducer() {
	suite.kp.Produce("localhost", "foobartopic", 9)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestKafkaProducerTestSuiteSuite(t *testing.T) {
	suite.Run(t, new(KafkaProducerTestSuite))
}
