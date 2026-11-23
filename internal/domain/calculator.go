package domain

import (
	"errors"
	"math"
)

type ScoreStrategy interface {
	Calculate(rubric *Rubric, evaluations []Evaluation) (float64, error)
}

type WeightedScoreStrategy struct{}

func (s WeightedScoreStrategy) Calculate(
	rubric *Rubric, evaluations []Evaluation) (float64, error) {
	if rubric == nil {
		return 0, errors.New("рубрикатор не задан")
	}
	if err := rubric.Validate(); err != nil {
		return 0, err
	}
	scores := make(map[string]int)
	for _, evaluation := range evaluations {
		criterion, exists := rubric.GetCriterion(evaluation.CriterionID())
		if !exists {
			return 0, errors.New("оценка содержит неизвестный критерий")
		}
		if err := criterion.ValidateScore(evaluation.Score()); err != nil {
			return 0, err
		}
		if _, exists := scores[evaluation.CriterionID()]; exists {
			return 0, errors.New("обнаружена повторная оценка по критерию")
		}
		scores[evaluation.CriterionID()] = evaluation.Score()
	}

	var total float64
	for _, criterion := range rubric.Criteria() {
		score := scores[criterion.ID()]
		total += float64(score) / float64(criterion.MaxScore()) * criterion.Weight() * 10
	}
	return math.Round(total*100) / 100, nil
}

// доменный сервис расчёта итога
type ScoreCalculator struct {
	strategy ScoreStrategy
}

func NewScoreCalculator(strategy ScoreStrategy) (*ScoreCalculator, error) {
	if strategy == nil {
		return nil, errors.New("стратегия расчёта не задана")
	}
	return &ScoreCalculator{strategy: strategy}, nil
}

func (c *ScoreCalculator) Calculate(
	rubric *Rubric, evaluations []Evaluation) (float64, error) {
	if c == nil || c.strategy == nil {
		return 0, errors.New("калькулятор не настроен")
	}
	return c.strategy.Calculate(rubric, evaluations)
}
