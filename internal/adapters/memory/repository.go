package memory

import (
	"errors"
	"pos/internal/domain"
	"sync"
)

type StudentRepository struct {
	mu       sync.RWMutex
	students map[string]domain.Student
}

func NewStudentRepository() *StudentRepository {
	return &StudentRepository{
		students: make(map[string]domain.Student),
	}
}

func (r *StudentRepository) Save(student domain.Student) error {
	if student.ID() == "" || student.Name() == "" {
		return errors.New("некорректные данные студента")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.students[student.ID()] = student
	return nil
}

func (r *StudentRepository) FindByID(id string) (domain.Student, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	student, exists := r.students[id]
	return student, exists
}

func (r *StudentRepository) FindAll() []domain.Student {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Student, 0, len(r.students))

	for _, student := range r.students {
		result = append(result, student)
	}

	return result
}

type RubricRepository struct {
	mu     sync.RWMutex
	rubric *domain.Rubric
}

func NewRubricRepository() *RubricRepository {
	return &RubricRepository{
		rubric: domain.NewRubric(),
	}
}

func (r *RubricRepository) Save(rubric *domain.Rubric) error {
	if rubric == nil {
		return errors.New("нельзя сохранить пустой рубрикатор")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	copyRubric, err := cloneRubric(rubric)
	if err != nil {
		return err
	}

	r.rubric = copyRubric
	return nil
}

func (r *RubricRepository) Get() (*domain.Rubric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.rubric == nil {
		return nil, errors.New("рубрикатор не найден")
	}

	return cloneRubric(r.rubric)
}

func cloneRubric(source *domain.Rubric) (*domain.Rubric, error) {
	result := domain.NewRubric()

	for _, criterion := range source.Criteria() {
		if err := result.AddCriterion(criterion); err != nil {
			return nil, err
		}
	}

	return result, nil
}

type EvaluationRepository struct {
	mu          sync.RWMutex
	evaluations map[string]map[string]domain.Evaluation
}

func NewEvaluationRepository() *EvaluationRepository {
	return &EvaluationRepository{
		evaluations: make(map[string]map[string]domain.Evaluation),
	}
}

func (r *EvaluationRepository) Save(evaluation domain.Evaluation) error {
	if evaluation.StudentID() == "" || evaluation.CriterionID() == "" {
		return errors.New("некорректная оценка")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	studentID := evaluation.StudentID()

	if r.evaluations[studentID] == nil {
		r.evaluations[studentID] = make(map[string]domain.Evaluation)
	}

	r.evaluations[studentID][evaluation.CriterionID()] = evaluation
	return nil
}

func (r *EvaluationRepository) Find(
	studentID, criterionID string,
) (domain.Evaluation, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	evaluation, exists := r.evaluations[studentID][criterionID]
	return evaluation, exists
}

func (r *EvaluationRepository) FindByStudent(
	studentID string,
) []domain.Evaluation {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Evaluation, 0)

	for _, evaluation := range r.evaluations[studentID] {
		result = append(result, evaluation)
	}

	return result
}

func (r *EvaluationRepository) Delete(
	studentID, criterionID string,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.evaluations[studentID][criterionID]; !exists {
		return errors.New("оценка не найдена")
	}

	delete(r.evaluations[studentID], criterionID)
	return nil
}
