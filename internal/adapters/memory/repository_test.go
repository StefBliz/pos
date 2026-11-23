package memory

import (
	"testing"

	"pos/internal/domain"
)

func TestStudentRepository(t *testing.T) {
	repo := NewStudentRepository()

	student, err := domain.NewStudent("ST-001", "Анна")
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.Save(student); err != nil {
		t.Fatal(err)
	}

	found, exists := repo.FindByID("ST-001")
	if !exists || found.Name() != "Анна" {
		t.Fatal("студент не найден")
	}

	if len(repo.FindAll()) != 1 {
		t.Fatal("неверное количество студентов")
	}

	if _, exists := repo.FindByID("UNKNOWN"); exists {
		t.Fatal("найден несуществующий студент")
	}

	if err := repo.Save(domain.Student{}); err == nil {
		t.Fatal("ожидалась ошибка для пустого студента")
	}
}

func TestRubricRepository(t *testing.T) {
	repo := NewRubricRepository()

	rubric, err := repo.Get()
	if err != nil {
		t.Fatal(err)
	}

	criterion, err := domain.NewCriterion(
		"CR-001", "Экзамен", 10, 1,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := rubric.AddCriterion(criterion); err != nil {
		t.Fatal(err)
	}

	if err := repo.Save(rubric); err != nil {
		t.Fatal(err)
	}

	saved, err := repo.Get()
	if err != nil {
		t.Fatal(err)
	}

	if len(saved.Criteria()) != 1 {
		t.Fatal("критерий не сохранился")
	}

	if err := repo.Save(nil); err == nil {
		t.Fatal("ожидалась ошибка для nil")
	}
}

func TestEvaluationRepository(t *testing.T) {
	repo := NewEvaluationRepository()
	rubric := domain.NewRubric()

	criterion, err := domain.NewCriterion(
		"CR-001", "Экзамен", 10, 1,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := rubric.AddCriterion(criterion); err != nil {
		t.Fatal(err)
	}

	evaluation, err := domain.NewEvaluation(
		"ST-001", "CR-001", 8, rubric,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.Save(evaluation); err != nil {
		t.Fatal(err)
	}

	found, exists := repo.Find("ST-001", "CR-001")
	if !exists || found.Score() != 8 {
		t.Fatal("оценка не сохранилась")
	}

	if len(repo.FindByStudent("ST-001")) != 1 {
		t.Fatal("неверное количество оценок")
	}

	if err := repo.Delete("ST-001", "CR-001"); err != nil {
		t.Fatal(err)
	}

	if _, exists := repo.Find("ST-001", "CR-001"); exists {
		t.Fatal("оценка не удалилась")
	}

	if err := repo.Delete("ST-001", "CR-001"); err == nil {
		t.Fatal("ожидалась ошибка повторного удаления")
	}

	if err := repo.Save(domain.Evaluation{}); err == nil {
		t.Fatal("ожидалась ошибка для пустой оценки")
	}
}
