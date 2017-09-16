package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type KafkaProducerTestSuite struct {
	suite.Suite
}

func (suite *KafkaProducerTestSuite) TestProducer() {
	produce("localhost", "foobar")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestKafkaProducerTestSuiteSuite(t *testing.T) {
	suite.Run(t, new(KafkaProducerTestSuite))
}
