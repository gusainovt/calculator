package main

import (
	"calculator/internal/calculationService"
	"calculator/internal/db"
	"calculator/internal/handlers"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	k "calculator/internal/kafka"
)

const (
	topic         = "calculator"
	consumerGroup = "calculator"
	numberOfKeys  = 20
)

var address = []string{"localhost:9091", "localhost:9092", "localhost:9093"}

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	calcRepo := calculationService.NewCalculationRepository(database)
	calcService := calculationService.NewCalculationService(calcRepo)
	calcHandler := handlers.NewCalculationHandler(calcService)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", calcHandler.GetCalculations)
	e.POST("/calculations", calcHandler.PostCalculations)
	e.PATCH("/calculations/:id", calcHandler.PatchCalculations)
	e.DELETE("/calculations/:id", calcHandler.DeleteCalculations)

	e.Start("localhost:8080")

	p, err := k.NewProducer(address)
	if err != nil {
		log.Fatal(err)
	}

	keys := generateUUIDString()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Kafka message %d", i)
		key := keys[i%numberOfKeys]
		if err := p.Produce(msg, topic, key); err != nil {
			log.Fatal(err)
		}
	}

	h := handlers.NewMessageHandler()
	c, err := k.NewConsumer(h, address, topic, consumerGroup)
	if err != nil {
		log.Fatal(err)
	}
	c.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Fatal(c.Stop())

}

func generateUUIDString() [numberOfKeys]string {
	var uuids [numberOfKeys]string
	for i := 0; i < numberOfKeys; i++ {
		uuids[i] = uuid.NewString()
	}
	return uuids
}
