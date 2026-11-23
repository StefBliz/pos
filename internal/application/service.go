package application

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"pos/internal/domain"
)

type Service struct {
	students    StudentRepository
	rubrics     RubricRepository
	evaluations EvaluationRepository
	calculator  *domain.ScoreCalculator
	exporter    ResultExporter
}

func NewService(
	students StudentRepository,
	rubrics RubricRepository,
	evaluations EvaluationRepository,
	calculator *domain.ScoreCalculator,
	exporter ResultExporter,
) (*Service, error) {
	if students == nil || rubrics == nil ||
		evaluations == nil || calculator == nil || exporter == nil {
		return nil, errors.New("не все зависимости приложения заданы")
	}
	return &Service{
		students:    students,
		rubrics:     rubrics,
		evaluations: evaluations,
		calculator:  calculator,
		exporter:    exporter,
	}, nil
}

func (s *Service) AddStudent(id, name string) error {
	student, err := domain.NewStudent(id, name)
	if err != nil {
		return err
	}
	if _, exists := s.students.FindByID(student.ID()); exists {
		return errors.New("студент с таким ID уже существует")
	}
	return s.students.Save(student)
}

func (s *Service) Students() []domain.Student {
	students := s.students.FindAll()

	sort.Slice(students, func(i, j int) bool {
		return students[i].ID() < students[j].ID()
	})
	return students
}

func (s *Service) AddCriterion(
	id, name string,
	maxScore int,
	weight float64,
) error {
	criterion, err := domain.NewCriterion(id, name, maxScore, weight)
	if err != nil {
		return err
	}
	rubric, err := s.rubrics.Get()
	if err != nil {
		return err
	}
	if rubric.TotalWeight()+weight > 1.0+1e-9 {
		return errors.New("сумма весов не может превышать 1")
	}
	if err := rubric.AddCriterion(criterion); err != nil {
		return err
	}
	return s.rubrics.Save(rubric)
}

func (s *Service) Criteria() ([]domain.Criterion, error) {
	rubric, err := s.rubrics.Get()
	if err != nil {
		return nil, err
	}

	criteria := rubric.Criteria()

	sort.Slice(criteria, func(i, j int) bool {
		return criteria[i].ID() < criteria[j].ID()
	})

	return criteria, nil
}

func (s *Service) SetEvaluation(
	studentID, criterionID string,
	score int,
) error {
	studentID = strings.TrimSpace(studentID)
	criterionID = strings.TrimSpace(criterionID)

	if _, exists := s.students.FindByID(studentID); !exists {
		return errors.New("студент не найден")
	}

	if _, exists := s.evaluations.Find(studentID, criterionID); exists {
		return errors.New("оценка уже существует, используйте изменение")
	}

	rubric, err := s.rubrics.Get()
	if err != nil {
		return err
	}

	evaluation, err := domain.NewEvaluation(
		studentID, criterionID, score, rubric,
	)
	if err != nil {
		return err
	}

	return s.evaluations.Save(evaluation)
}

func (s *Service) ChangeEvaluation(
	studentID, criterionID string,
	score int,
) error {
	studentID = strings.TrimSpace(studentID)
	criterionID = strings.TrimSpace(criterionID)

	evaluation, exists := s.evaluations.Find(studentID, criterionID)
	if !exists {
		return errors.New("оценка не найдена")
	}

	rubric, err := s.rubrics.Get()
	if err != nil {
		return err
	}

	if err := evaluation.ChangeScore(score, rubric); err != nil {
		return err
	}

	return s.evaluations.Save(evaluation)
}

func (s *Service) DeleteEvaluation(studentID, criterionID string) error {
	return s.evaluations.Delete(
		strings.TrimSpace(studentID),
		strings.TrimSpace(criterionID),
	)
}

func (s *Service) Evaluations(studentID string) ([]domain.Evaluation, error) {
	studentID = strings.TrimSpace(studentID)

	if _, exists := s.students.FindByID(studentID); !exists {
		return nil, errors.New("студент не найден")
	}

	evaluations := s.evaluations.FindByStudent(studentID)

	sort.Slice(evaluations, func(i, j int) bool {
		return evaluations[i].CriterionID() < evaluations[j].CriterionID()
	})

	return evaluations, nil
}

func (s *Service) CalculateResult(studentID string) (Result, error) {
	studentID = strings.TrimSpace(studentID)

	student, exists := s.students.FindByID(studentID)
	if !exists {
		return Result{}, errors.New("студент не найден")
	}

	rubric, err := s.rubrics.Get()
	if err != nil {
		return Result{}, err
	}

	evaluations := s.evaluations.FindByStudent(studentID)

	score, err := s.calculator.Calculate(rubric, evaluations)
	if err != nil {
		return Result{}, err
	}

	return Result{
		StudentID:   student.ID(),
		StudentName: student.Name(),
		Score:       score,
	}, nil
}

func (s *Service) ExportResults(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("путь экспорта не может быть пустым")
	}

	students := s.Students()
	results := make([]Result, 0, len(students))

	for _, student := range students {
		result, err := s.CalculateResult(student.ID())
		if err != nil {
			return fmt.Errorf(
				"не удалось рассчитать итог студента %s: %w",
				student.ID(), err,
			)
		}

		results = append(results, result)
	}

	return s.exporter.Export(results, path)
}

func (s *Service) CreateDefaultRubric() error {
	rubric, err := s.rubrics.Get()
	if err != nil {
		return err
	}

	if len(rubric.Criteria()) != 0 {
		return errors.New("рубрикатор уже содержит критерии")
	}

	defaultRubric, err := domain.NewRubricBuilder().
		Add("CR-001", "Домашние задания", 10, 0.4).
		Add("CR-002", "Контрольная работа", 10, 0.3).
		Add("CR-003", "Экзамен", 10, 0.3).
		Build()

	if err != nil {
		return err
	}

	return s.rubrics.Save(defaultRubric)
}
