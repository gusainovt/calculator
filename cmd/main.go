package main

import (
	"calculator/internal/calculationService"
	"calculator/internal/db"
	"calculator/internal/handlers"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	k "calculator/internal/kafka"
)

const (
	topic = "calculator"
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

	p, err := k.NewProducer(address)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Kafka message %d", i)
		if err := p.Produce(msg, topic); err != nil {
			log.Fatal(err)
		}
	}

	e.Start("localhost:8080")
}
