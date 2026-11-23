package domain

import "errors"

type Score struct {
	value int
}

func NewScore(value, maxScore int) (Score, error) {
	if maxScore <= 0 {
		return Score{}, errors.New("максимальный балл должен быть положительным")
	}

	if value < 0 || value > maxScore {
		return Score{}, errors.New("балл выходит за допустимые границы")
	}

	return Score{value: value}, nil
}

func (s Score) Value() int {
	return s.value
}
