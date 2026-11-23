package domain

import (
	"math"
	"testing"
)

func TestInvalidWeights(t *testing.T) {
	values := []float64{
		math.NaN(),
		math.Inf(1),
		math.Inf(-1),
		-0.5,
		0,
		1.5,
	}

	for _, value := range values {
		if _, err := NewWeight(value); err == nil {
			t.Fatalf("вес %v должен быть запрещён", value)
		}
	}
}

func TestChangeEvaluationInvalid(t *testing.T) {
	rubric := NewRubric()

	criterion, err := NewCriterion("CR-001", "Экзамен", 10, 1)
	if err != nil {
		t.Fatal(err)
	}

	if err := rubric.AddCriterion(criterion); err != nil {
		t.Fatal(err)
	}

	evaluation, err := NewEvaluation("ST-001", "CR-001", 8, rubric)
	if err != nil {
		t.Fatal(err)
	}

	if err := evaluation.ChangeScore(15, rubric); err == nil {
		t.Fatal("ожидалась ошибка изменения оценки")
	}

	if evaluation.Score() != 8 {
		t.Fatal("после ошибки исходная оценка должна сохраниться")
	}
}

func TestRubricWeightOverflow(t *testing.T) {
	rubric := NewRubric()

	first, _ := NewCriterion("CR-001", "Домашние задания", 10, 0.8)
	second, _ := NewCriterion("CR-002", "Экзамен", 10, 0.5)

	if err := rubric.AddCriterion(first); err != nil {
		t.Fatal(err)
	}

	if err := rubric.AddCriterion(second); err == nil {
		t.Fatal("сумма весов больше 1 должна быть запрещена")
	}
}
