package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=1234 dbname=calculator_db port=5432 sslmode=disable"
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Calculation{}); err != nil {
		log.Fatalf("Could not migrate table: %v", err)
	}
}

type Calculation struct {
	ID         string `gorm:"primaryKey" json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

func calculateExpression(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return "", err
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), err
}

func getCalculations(c echo.Context) error {
	var calculations []Calculation
	if err := db.Find(&calculations).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get calculations"})
	}
	return c.JSON(http.StatusOK, calculations)
}

//func postCalculations(c echo.Context) error {
//	var req CalculationRequest
//
//	if err := c.Bind(&req); err != nil {
//		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
//	}
//
//	result, err := calculateExpression(req.Expression)
//
//	if err != nil {
//		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
//	}
//	calc := Calculation{
//		ID:         uuid.NewString(),
//		Expression: req.Expression,
//		Result:     result,
//	}
//	calculations = append(calculations, calc)
//	return c.JSON(http.StatusCreated, calc)
//}
//
//func patchCalculations(c echo.Context) error {
//	id := c.Param("id")
//	var req CalculationRequest
//
//	if err := c.Bind(&req); err != nil {
//		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
//	}
//
//	result, err := calculateExpression(req.Expression)
//	if err != nil {
//		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
//	}
//
//	for i, calculation := range calculations {
//		if calculation.ID == id {
//			calculations[i].Expression = req.Expression
//			calculations[i].Result = result
//			return c.JSON(http.StatusOK, calculations[i])
//		}
//	}
//	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Calculation not found"})
//}
//
//func deleteCalculations(c echo.Context) error {
//	id := c.Param("id")
//
//	for i, calculation := range calculations {
//		if calculation.ID == id {
//			calculations = append(calculations[:i], calculations[i+1:]...)
//			return c.NoContent(http.StatusNoContent)
//		}
//	}
//	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Calculation not found"})
//}

func main() {
	initDB()
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", getCalculations)
	//e.POST("/calculations", postCalculations)
	//e.PATCH("/calculations/:id", patchCalculations)
	//e.DELETE("/calculations/:id", deleteCalculations)

	e.Start("localhost:8080")
}
