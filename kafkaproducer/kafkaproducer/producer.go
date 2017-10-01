package main

import (
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	Producer *kafka.Producer
	broker   string
}

func (k *KafkaProducer) init() *kafka.Producer {
	broker := k.broker
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

func (k *KafkaProducer) readNotifFromChan(deliveryChan <-chan kafka.Event, ackChan chan<- string, count int) {
	// <-deliveryChan is essentially to receive one item from the channel (the delivery report)

	// count := 99
	// i := 0
	var offset string

	for i := 0; i < count; i++ {
		e := <-deliveryChan
		// for e := range deliveryChan {
		// time.Sleep(20 * time.Millisecond)

		// get (pop) an item
		// e := <-deliveryChan
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

func (k *KafkaProducer) Produce(broker string, topic string, count int) {

	fmt.Printf("Created Producer %v\n", k.Producer)

	// the 'count' is the maximum number of messages on the channel
	// Event notification delivery channel (commit acks)
	deliveryChan := make(chan kafka.Event, count)
	ackChan := make(chan string, count)
	go k.readNotifFromChan(deliveryChan, ackChan, count)

	// send n messages
	for i := 0; i < count; i++ {
		value := fmt.Sprintf("Hello go # %d", i)
		msg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(value),
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
}
