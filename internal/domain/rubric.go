package domain

import (
	"errors"
	"math"
	"sort"
)

type Rubric struct {
	criteria map[string]Criterion
}

func NewRubric() *Rubric {
	return &Rubric{
		criteria: make(map[string]Criterion),
	}
}

func (r *Rubric) AddCriterion(criterion Criterion) error {
	if criterion.ID() == "" || criterion.Name() == "" ||
		criterion.MaxScore() <= 0 {
		return errors.New("некорректный критерий")
	}

	if _, err := NewWeight(criterion.Weight()); err != nil {
		return err
	}

	if r.criteria == nil {
		r.criteria = make(map[string]Criterion)
	}

	if _, exists := r.criteria[criterion.ID()]; exists {
		return errors.New("критерий с таким ID уже существует")
	}

	if r.TotalWeight()+criterion.Weight() > 1+1e-9 {
		return errors.New("сумма весов не может превышать 1")
	}

	r.criteria[criterion.ID()] = criterion
	return nil
}

func (r *Rubric) GetCriterion(id string) (Criterion, bool) {
	criterion, exists := r.criteria[id]
	return criterion, exists
}

func (r *Rubric) Criteria() []Criterion {
	result := make([]Criterion, 0, len(r.criteria))

	for _, criterion := range r.criteria {
		result = append(result, criterion)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID() < result[j].ID()
	})

	return result
}

func (r *Rubric) TotalWeight() float64 {
	var total float64

	for _, criterion := range r.criteria {
		total += criterion.Weight()
	}

	return total
}

func (r *Rubric) Validate() error {
	if len(r.criteria) == 0 {
		return errors.New("рубрикатор пуст")
	}

	if math.Abs(r.TotalWeight()-1) > 1e-9 {
		return errors.New("сумма весов критериев должна равняться 1")
	}

	return nil
}
