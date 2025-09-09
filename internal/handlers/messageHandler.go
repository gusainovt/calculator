package handlers

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/labstack/gommon/log"
)

type Handler struct {
}

func NewMessageHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleMassage(massage []byte, offset kafka.Offset) error {
	log.Infof("Message from kafka with ofsset: %d '%s", offset, string(massage))
	return nil
}
