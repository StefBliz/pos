package domain

import (
	"errors"
	"strings"
)

type Criterion struct {
	id       string
	name     string
	maxScore int
	weight   Weight
}

func NewCriterion(
	id string,
	name string,
	maxScore int,
	weight float64,
) (Criterion, error) {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)

	if id == "" {
		return Criterion{}, errors.New("ID критерия не может быть пустым")
	}

	if name == "" {
		return Criterion{}, errors.New("название критерия не может быть пустым")
	}

	if maxScore <= 0 {
		return Criterion{}, errors.New("максимальный балл должен быть положительным")
	}

	validWeight, err := NewWeight(weight)
	if err != nil {
		return Criterion{}, err
	}

	return Criterion{
		id:       id,
		name:     name,
		maxScore: maxScore,
		weight:   validWeight,
	}, nil
}

func (c Criterion) ID() string {
	return c.id
}

func (c Criterion) Name() string {
	return c.name
}

func (c Criterion) MaxScore() int {
	return c.maxScore
}

func (c Criterion) Weight() float64 {
	return c.weight.Value()
}

func (c Criterion) ValidateScore(value int) error {
	_, err := NewScore(value, c.maxScore)
	return err
}
