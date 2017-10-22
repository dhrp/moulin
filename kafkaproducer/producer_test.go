package kafkaproducer

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type KafkaProducerTestSuite struct {
	suite.Suite
	kp KFK
}

func (suite *KafkaProducerTestSuite) SetupSuite() {
	suite.kp = KFK{Broker: "localhost:9092"}
	_ = suite.kp.Init()
}

func (suite *KafkaProducerTestSuite) TestProducer() {
	var values [][]byte

	// add a value to the slice
	values = append(values, []byte("hello world"))
	result, err := suite.kp.Produce("foobartopic", values)
	suite.Nil(err, "Didn't expect and error from Produce")
	suite.Contains(result, "the offset of this run was", "result invalid")

}

func (suite *KafkaProducerTestSuite) TestProduceFromFile() {
	result, _ := suite.kp.ProduceFromFile("./test/testtextfile.txt")
	suite.Contains(result, "the offset of this run was", "result invalid")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestKafkaProducerTestSuiteSuite(t *testing.T) {
	suite.Run(t, new(KafkaProducerTestSuite))
}
