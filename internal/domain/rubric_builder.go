package domain

import "errors"

type RubricBuilder struct {
	rubric *Rubric
	err    error
}

func NewRubricBuilder() *RubricBuilder {
	return &RubricBuilder{
		rubric: NewRubric(),
	}
}

func (b *RubricBuilder) Add(
	id, name string,
	maxScore int,
	weight float64,
) *RubricBuilder {
	if b.err != nil {
		return b
	}

	criterion, err := NewCriterion(id, name, maxScore, weight)
	if err != nil {
		b.err = err
		return b
	}

	b.err = b.rubric.AddCriterion(criterion)
	return b
}

func (b *RubricBuilder) Build() (*Rubric, error) {
	if b == nil || b.rubric == nil {
		return nil, errors.New("сборщик рубрикатора не создан")
	}

	if b.err != nil {
		return nil, b.err
	}

	if err := b.rubric.Validate(); err != nil {
		return nil, err
	}

	result := NewRubric()

	for _, criterion := range b.rubric.Criteria() {
		if err := result.AddCriterion(criterion); err != nil {
			return nil, err
		}
	}

	return result, nil
}
