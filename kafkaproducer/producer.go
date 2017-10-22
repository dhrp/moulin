package kafkaproducer

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
)

// KFK is the struct that is the 'root' of Kafka services
type KFK struct {
	Producer *kafka.Producer
	Broker   string
}

// Init initializes KFK with the expected settings.
func (k *KFK) Init() *kafka.Producer {
	broker := k.Broker
	hostname, _ := os.Hostname()

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  broker,
		"batch.num.messages": 7,
		"client.id":          hostname,
		"acks":               "all",
	})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	k.Producer = p
	return k.Producer
}

func (k *KFK) readNotifFromChan(deliveryChan <-chan kafka.Event, ackChan chan<- string, count int) {
	var offset string

	for i := 0; i < count; i++ {
		e := <-deliveryChan

		// get the msg
		m, ok := e.(*kafka.Message)
		if !ok {
			log.Panic("message not received")
			continue
		}

		if m.TopicPartition.Error != nil {
			fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		} else {
			fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
				*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		}
		offset = m.TopicPartition.Offset.String()
	}
	// sleep at the end
	mstr := fmt.Sprintf("the offset of this run was %s", offset)
	ackChan <- mstr
}

// Produce is the mechanism by which we actually push messages into Kafka
func (k *KFK) Produce(topic string, values [][]byte) (string, error) {

	fmt.Printf("Using producer: %v\n", k.Producer)

	// the len(values) is the maximum number of messages on the channel
	// event notification delivery channel (commit acks)
	// we set it as big as the amount of items we've storing so that we alway get
	// a quick reply
	deliveryChan := make(chan kafka.Event, len(values))
	ackChan := make(chan string, len(values))
	go k.readNotifFromChan(deliveryChan, ackChan, len(values))

	// send n messages
	for _, value := range values {
		msg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          value,
		}
		err := k.Producer.Produce(msg, deliveryChan)
		if err != nil {
			log.Panic(err)
		}
	}

	// here we wait (blocking) to receive a message back.
	// we'll only receive the overview after completion.
	result := <-ackChan
	close(ackChan)
	log.Println(result)

	// close the channel
	close(deliveryChan)
	return result, nil
}

// ProduceFromFile is a stub for reading the uploaded file
// and then feeding it into the Kafka Producer
func (k *KFK) ProduceFromFile(filePath string) (string, error) {
	log.Println("Doing ProduceFromFile")
	if k == nil {
		log.Panic("kfk not initialized in ProduceFromFile")
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	// var line string
	scanner := bufio.NewScanner(file)

	i := 0
	var values [][]byte

	for scanner.Scan() {
		i++
		fmt.Println(i)
		fmt.Println(scanner.Text())
		values = append(values, scanner.Bytes())

		if i == 1000 {
			// finalize
		}
	}

	// check if scanner has any errors before reading.
	if err = scanner.Err(); err != nil {
		log.Panic(err)
	}

	result, err := k.Produce("foobar", values)
	if err != nil {
		return "", errors.Wrap(err, "failed ProduceFromFile")
	}

	return result, nil
}
