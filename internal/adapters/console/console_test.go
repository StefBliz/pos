package console

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"pos/internal/adapters/memory"
	"pos/internal/application"
	"pos/internal/domain"
)

func newTestConsole(t *testing.T, input string) (*Console, *bytes.Buffer) {
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
		&testExporter{},
	)
	if err != nil {
		t.Fatal(err)
	}

	output := &bytes.Buffer{}

	return &Console{
		service: service,
		reader:  bufio.NewReader(strings.NewReader(input)),
		writer:  output,
	}, output
}

// testExporter заменяет запись файла в тестах консоли.
type testExporter struct{}

func (e *testExporter) Export(
	results []application.Result,
	path string,
) error {
	return nil
}

func TestConsoleExit(t *testing.T) {
	console, output := newTestConsole(t, "0\n")

	if err := console.Run(); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(output.String(), "Работа завершена") {
		t.Fatal("не найдено сообщение о завершении")
	}
}

func TestConsoleAddStudent(t *testing.T) {
	console, output := newTestConsole(
		t,
		"1\nST-001\nАнна\n2\n0\n",
	)

	if err := console.Run(); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(output.String(), "ST-001 — Анна") {
		t.Fatal("студент не отображается")
	}
}

func TestConsoleInvalidCommand(t *testing.T) {
	console, output := newTestConsole(t, "99\n0\n")

	if err := console.Run(); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(output.String(), "неизвестная команда") {
		t.Fatal("не обработана неизвестная команда")
	}
}

func TestConsoleInvalidNumber(t *testing.T) {
	console, output := newTestConsole(
		t,
		"3\nCR-001\nЭкзамен\nabc\n0\n",
	)

	if err := console.Run(); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(output.String(), "ожидалось целое число") {
		t.Fatal("не обработан неправильный числовой ввод")
	}
}

func TestConsoleDefaultRubric(t *testing.T) {
	console, output := newTestConsole(t, "11\n4\n0\n")

	if err := console.Run(); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(output.String(), "CR-003") {
		t.Fatal("стандартный рубрикатор не создан")
	}
}

func TestConsoleEOF(t *testing.T) {
	console, _ := newTestConsole(t, "")

	if err := console.Run(); err != nil {
		t.Fatal(err)
	}
}
