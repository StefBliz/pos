package domain

import (
	"math"
	"testing"
)

func TestWeightedScore(t *testing.T) {
	rubric := NewRubric()

	data := []struct {
		id     string
		name   string
		weight float64
		score  int
	}{
		{"CR-001", "Домашние задания", 0.4, 8},
		{"CR-002", "Контрольная", 0.3, 9},
		{"CR-003", "Экзамен", 0.3, 10},
	}

	var evaluations []Evaluation

	for _, item := range data {
		criterion, err := NewCriterion(item.id, item.name, 10, item.weight)
		if err != nil {
			t.Fatal(err)
		}

		if err := rubric.AddCriterion(criterion); err != nil {
			t.Fatal(err)
		}

		evaluation, err := NewEvaluation(
			"ST-001", item.id, item.score, rubric,
		)
		if err != nil {
			t.Fatal(err)
		}

		evaluations = append(evaluations, evaluation)
	}

	calculator, err := NewScoreCalculator(WeightedScoreStrategy{})
	if err != nil {
		t.Fatal(err)
	}

	result, err := calculator.Calculate(rubric, evaluations)
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(result-8.9) > 1e-9 {
		t.Fatalf("ожидалось 8.9, получено %.2f", result)
	}
}

func TestEvaluationAboveMaximum(t *testing.T) {
	rubric := NewRubric()

	criterion, err := NewCriterion("CR-001", "Экзамен", 10, 1)
	if err != nil {
		t.Fatal(err)
	}

	if err := rubric.AddCriterion(criterion); err != nil {
		t.Fatal(err)
	}

	_, err = NewEvaluation("ST-001", "CR-001", 12, rubric)
	if err == nil {
		t.Fatal("ожидалась ошибка при превышении максимума")
	}
}

func TestUnknownCriterion(t *testing.T) {
	rubric := NewRubric()

	_, err := NewEvaluation("ST-001", "UNKNOWN", 5, rubric)
	if err == nil {
		t.Fatal("ожидалась ошибка для неизвестного критерия")
	}
}

func TestDuplicateCriterion(t *testing.T) {
	rubric := NewRubric()

	criterion, err := NewCriterion("CR-001", "Экзамен", 10, 1)
	if err != nil {
		t.Fatal(err)
	}

	if err := rubric.AddCriterion(criterion); err != nil {
		t.Fatal(err)
	}

	if err := rubric.AddCriterion(criterion); err == nil {
		t.Fatal("ожидалась ошибка при повторном добавлении")
	}
}

func TestInvalidTotalWeight(t *testing.T) {
	rubric := NewRubric()

	criterion, err := NewCriterion("CR-001", "Экзамен", 10, 0.5)
	if err != nil {
		t.Fatal(err)
	}

	if err := rubric.AddCriterion(criterion); err != nil {
		t.Fatal(err)
	}

	if err := rubric.Validate(); err == nil {
		t.Fatal("ожидалась ошибка из-за суммы весов")
	}
}
