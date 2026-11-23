package domain

import (
	"errors"
	"strings"
)

// оценка студента по конкретному критерию
type Evaluation struct {
	studentID   string
	criterionID string
	score       int
}

func NewEvaluation(studentID string, criterionID string, score int, rubric *Rubric) (Evaluation, error) {
	studentID = strings.TrimSpace(studentID)
	criterionID = strings.TrimSpace(criterionID)
	if studentID == "" {
		return Evaluation{}, errors.New("ID студента не может быть пустым")
	}
	if rubric == nil {
		return Evaluation{}, errors.New("рубрикатор не задан")
	}
	criterion, exists := rubric.GetCriterion(criterionID)
	if !exists {
		return Evaluation{}, errors.New("критерий отсутствует в рубрикаторе")
	}
	if err := criterion.ValidateScore(score); err != nil {
		return Evaluation{}, err
	}
	return Evaluation{
		studentID:   studentID,
		criterionID: criterionID,
		score:       score,
	}, nil
}

func (e Evaluation) StudentID() string {
	return e.studentID
}

func (e Evaluation) CriterionID() string {
	return e.criterionID
}

func (e Evaluation) Score() int {
	return e.score
}

func (e *Evaluation) ChangeScore(score int, rubric *Rubric) error {
	if rubric == nil {
		return errors.New("рубрикатор не задан")
	}
	criterion, exists := rubric.GetCriterion(e.criterionID)
	if !exists {
		return errors.New("критерий отсутствует в рубрикаторе")
	}
	if err := criterion.ValidateScore(score); err != nil {
		return err
	}
	e.score = score
	return nil
}
