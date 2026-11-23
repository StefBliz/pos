package application

import "pos/internal/domain"

type StudentRepository interface {
	Save(student domain.Student) error
	FindByID(id string) (domain.Student, bool)
	FindAll() []domain.Student
}

type RubricRepository interface {
	Save(rubric *domain.Rubric) error
	Get() (*domain.Rubric, error)
}

type EvaluationRepository interface {
	Save(evaluation domain.Evaluation) error
	Find(studentID, criterionID string) (domain.Evaluation, bool)
	FindByStudent(studentID string) []domain.Evaluation
	Delete(studentID, criterionID string) error
}

type Result struct {
	StudentID   string
	StudentName string
	Score       float64
}

type ResultExporter interface {
	Export(results []Result, path string) error
}
