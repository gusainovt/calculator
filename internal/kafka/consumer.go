package kafka

import (
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	sessionTimeout = 500
	noTimeout      = -1
)

type Handler interface {
	HandleMassage(massage []byte, offset kafka.Offset) error
}

type Consumer struct {
	consumer *kafka.Consumer
	handler  Handler
	stop     bool
}

func NewConsumer(handler Handler, address []string, topic, consumerGroup string) (*Consumer, error) {
	cfg := kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(address, ";"),
		"group.id":                 consumerGroup,
		"session.timeout.ms":       sessionTimeout,
		"enable.auto.offset.store": true,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  5000,
	}
	c, err := kafka.NewConsumer(&cfg)
	if err != nil {
		return nil, err
	}
	if err = c.Subscribe(topic, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
		handler:  handler,
	}, nil
}

func (c *Consumer) Start() {
	for {
		if c.stop {
			break
		}
		kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
		if err != nil {
			log.Fatalf(err.Error())
		}
		if kafkaMsg == nil {
			continue
		}
		if err = c.handler.HandleMassage(kafkaMsg.Value, kafkaMsg.TopicPartition.Offset); err != nil {
			log.Fatalf(err.Error())
		}
		if _, err = c.consumer.StoreMessage(kafkaMsg); err != nil {
			log.Fatalf(err.Error())
		}
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	if _, err := c.consumer.Commit(); err != nil {
		return err
	}
	return c.consumer.Close()
}
