package calculationService

import (
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
)

type CalculationService interface {
	CreateCalculation(expression string) (Calculation, error)
	GetAllCalculations() ([]Calculation, error)
	GetCalculationById(id string) (Calculation, error)
	UpdateCalculation(id, expression string) (Calculation, error)
	DeleteCalculation(id string) error
}

type calculationService struct {
	repo CalculationRepository
}

func (s *calculationService) CreateCalculation(expression string) (Calculation, error) {
	result, err := s.calculateExpression(expression)
	if err != nil {
		return Calculation{}, err
	}

	calc := Calculation{
		ID:         uuid.NewString(),
		Expression: expression,
		Result:     result,
	}

	if err := s.repo.CreateCalculation(calc); err != nil {
		return Calculation{}, err
	}

	return calc, nil
}

func (s *calculationService) GetAllCalculations() ([]Calculation, error) {
	return s.repo.GetAllCalculations()
}

func (s *calculationService) GetCalculationById(id string) (Calculation, error) {
	return s.repo.GetCalculationById(id)
}

func (s *calculationService) UpdateCalculation(id, expression string) (Calculation, error) {
	calc, err := s.repo.GetCalculationById(id)
	if err != nil {
		return Calculation{}, err
	}

	result, err := s.calculateExpression(expression)
	if err != nil {
		return Calculation{}, err
	}

	calc.Expression = expression
	calc.Result = result

	if err := s.repo.UpdateCalculation(calc); err != nil {
		return Calculation{}, err
	}

	return calc, nil
}

func (s *calculationService) DeleteCalculation(id string) error {
	return s.repo.DeleteCalculation(id)
}

func (s *calculationService) calculateExpression(expression string) (string, error) {
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

func NewCalculationService(r CalculationRepository) CalculationService {
	return &calculationService{repo: r}
}
