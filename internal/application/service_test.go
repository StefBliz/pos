package application_test

import (
	"math"
	"path/filepath"
	"testing"

	csvadapter "pos/internal/adapters/csv"
	"pos/internal/adapters/memory"
	"pos/internal/application"
	"pos/internal/domain"
)

func newTestService(t *testing.T) *application.Service {
	t.Helper()

	calculator, err := domain.NewScoreCalculator(
		domain.WeightedScoreStrategy{},
	)
	if err != nil {
		t.Fatal(err)
	}

	service, err := application.NewService(
		memory.NewStudentRepository(),
		memory.NewRubricRepository(),
		memory.NewEvaluationRepository(),
		calculator,
		csvadapter.NewExporter(),
	)
	if err != nil {
		t.Fatal(err)
	}

	return service
}

func TestAddStudent(t *testing.T) {
	service := newTestService(t)

	if err := service.AddStudent("ST-001", "Анна"); err != nil {
		t.Fatal(err)
	}

	students := service.Students()

	if len(students) != 1 || students[0].Name() != "Анна" {
		t.Fatal("студент добавлен некорректно")
	}
}

func TestDuplicateStudent(t *testing.T) {
	service := newTestService(t)

	if err := service.AddStudent("ST-001", "Анна"); err != nil {
		t.Fatal(err)
	}

	if err := service.AddStudent("ST-001", "Мария"); err == nil {
		t.Fatal("ожидалась ошибка повторного ID")
	}
}

func TestEvaluationLifecycle(t *testing.T) {
	service := newTestService(t)

	if err := service.AddStudent("ST-001", "Анна"); err != nil {
		t.Fatal(err)
	}

	if err := service.AddCriterion("CR-001", "Экзамен", 10, 1); err != nil {
		t.Fatal(err)
	}

	if err := service.SetEvaluation("ST-001", "CR-001", 8); err != nil {
		t.Fatal(err)
	}

	if err := service.SetEvaluation("ST-001", "CR-001", 7); err == nil {
		t.Fatal("повторная оценка должна быть запрещена")
	}

	if err := service.ChangeEvaluation("ST-001", "CR-001", 9); err != nil {
		t.Fatal(err)
	}

	result, err := service.CalculateResult("ST-001")
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(result.Score-9) > 1e-9 {
		t.Fatalf("ожидалось 9, получено %.2f", result.Score)
	}

	if err := service.DeleteEvaluation("ST-001", "CR-001"); err != nil {
		t.Fatal(err)
	}

	evaluations, err := service.Evaluations("ST-001")
	if err != nil {
		t.Fatal(err)
	}

	if len(evaluations) != 0 {
		t.Fatal("оценка не была удалена")
	}
}

func TestInvalidEvaluation(t *testing.T) {
	service := newTestService(t)

	if err := service.AddStudent("ST-001", "Анна"); err != nil {
		t.Fatal(err)
	}

	if err := service.AddCriterion("CR-001", "Экзамен", 10, 1); err != nil {
		t.Fatal(err)
	}

	if err := service.SetEvaluation("ST-001", "CR-001", 12); err == nil {
		t.Fatal("оценка выше максимума должна быть запрещена")
	}

	if err := service.SetEvaluation("UNKNOWN", "CR-001", 5); err == nil {
		t.Fatal("несуществующий студент должен быть запрещён")
	}

	if err := service.SetEvaluation("ST-001", "UNKNOWN", 5); err == nil {
		t.Fatal("несуществующий критерий должен быть запрещён")
	}
}

func TestExportResults(t *testing.T) {
	service := newTestService(t)

	if err := service.AddStudent("ST-001", "Анна"); err != nil {
		t.Fatal(err)
	}

	if err := service.AddCriterion("CR-001", "Экзамен", 10, 1); err != nil {
		t.Fatal(err)
	}

	if err := service.SetEvaluation("ST-001", "CR-001", 9); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(t.TempDir(), "results.csv")

	if err := service.ExportResults(path); err != nil {
		t.Fatal(err)
	}
}

func TestCreateDefaultRubric(t *testing.T) {
	service := newTestService(t)

	if err := service.CreateDefaultRubric(); err != nil {
		t.Fatal(err)
	}

	criteria, err := service.Criteria()
	if err != nil {
		t.Fatal(err)
	}

	if len(criteria) != 3 {
		t.Fatalf("ожидалось 3 критерия, получено %d", len(criteria))
	}

	if err := service.CreateDefaultRubric(); err == nil {
		t.Fatal("повторное создание рубрикатора должно быть запрещено")
	}
}
