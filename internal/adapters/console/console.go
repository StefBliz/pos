package console

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"pos/internal/application"
)

type Console struct {
	service *application.Service
	reader  *bufio.Reader
	writer  io.Writer
}

func NewConsole(service *application.Service) *Console {
	return &Console{
		service: service,
		reader:  bufio.NewReader(os.Stdin),
		writer:  os.Stdout,
	}
}

func (c *Console) readLine(prompt string) (string, error) {
	fmt.Fprint(c.writer, prompt)

	line, err := c.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}

	if err == io.EOF && line == "" {
		return "", io.EOF
	}

	return strings.TrimSpace(line), nil
}

func (c *Console) readInt(prompt string) (int, error) {
	value, err := c.readLine(prompt)
	if err != nil {
		return 0, err
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("ожидалось целое число: %w", err)
	}

	return number, nil
}

func (c *Console) readFloat(prompt string) (float64, error) {
	value, err := c.readLine(prompt)
	if err != nil {
		return 0, err
	}

	value = strings.ReplaceAll(value, ",", ".")

	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("ожидалось число: %w", err)
	}

	return number, nil
}

func (c *Console) menu() {
	fmt.Fprintln(c.writer, "\n========== ПОС ==========")
	fmt.Fprintln(c.writer, "1. Добавить студента")
	fmt.Fprintln(c.writer, "2. Список студентов")
	fmt.Fprintln(c.writer, "3. Добавить критерий")
	fmt.Fprintln(c.writer, "4. Список критериев")
	fmt.Fprintln(c.writer, "5. Выставить оценку")
	fmt.Fprintln(c.writer, "6. Оценки студента")
	fmt.Fprintln(c.writer, "7. Изменить оценку")
	fmt.Fprintln(c.writer, "8. Удалить оценку")
	fmt.Fprintln(c.writer, "9. Рассчитать итог")
	fmt.Fprintln(c.writer, "10. Экспортировать итоги в CSV")
	fmt.Fprintln(c.writer, "11. Создать стандартный рубрикатор")
	fmt.Fprintln(c.writer, "0. Выход")
	fmt.Fprintln(c.writer, "=========================")
}

func (c *Console) Run() error {
	for {
		c.menu()

		choice, err := c.readLine("Выберите действие: ")
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if choice == "0" {
			fmt.Fprintln(c.writer, "Работа завершена.")
			return nil
		}

		err = c.execute(choice)

		if err == io.EOF {
			return nil
		}

		if err != nil {
			fmt.Fprintln(c.writer, "Ошибка:", err)
		} else {
			fmt.Fprintln(c.writer, "Операция выполнена.")
		}
	}
}

func (c *Console) execute(choice string) error {
	switch choice {
	case "1":
		return c.addStudent()
	case "2":
		return c.showStudents()
	case "3":
		return c.addCriterion()
	case "4":
		return c.showCriteria()
	case "5":
		return c.setEvaluation()
	case "6":
		return c.showEvaluations()
	case "7":
		return c.changeEvaluation()
	case "8":
		return c.deleteEvaluation()
	case "9":
		return c.calculateResult()
	case "10":
		return c.exportResults()
	case "11":
		return c.service.CreateDefaultRubric()
	default:
		return fmt.Errorf("неизвестная команда: %s", choice)
	}
}

func (c *Console) addStudent() error {
	id, err := c.readLine("ID студента: ")
	if err != nil {
		return err
	}

	name, err := c.readLine("Имя студента: ")
	if err != nil {
		return err
	}

	return c.service.AddStudent(id, name)
}

func (c *Console) showStudents() error {
	students := c.service.Students()

	if len(students) == 0 {
		fmt.Fprintln(c.writer, "Студентов пока нет.")
		return nil
	}

	for _, student := range students {
		fmt.Fprintf(
			c.writer,
			"%s — %s\n",
			student.ID(),
			student.Name(),
		)
	}

	return nil
}

func (c *Console) addCriterion() error {
	id, err := c.readLine("ID критерия: ")
	if err != nil {
		return err
	}

	name, err := c.readLine("Название: ")
	if err != nil {
		return err
	}

	maxScore, err := c.readInt("Максимальный балл: ")
	if err != nil {
		return err
	}

	weight, err := c.readFloat("Вес от 0 до 1: ")
	if err != nil {
		return err
	}

	return c.service.AddCriterion(id, name, maxScore, weight)
}

func (c *Console) showCriteria() error {
	criteria, err := c.service.Criteria()
	if err != nil {
		return err
	}

	if len(criteria) == 0 {
		fmt.Fprintln(c.writer, "Критериев пока нет.")
		return nil
	}

	for _, criterion := range criteria {
		fmt.Fprintf(
			c.writer,
			"%s — %s | максимум: %d | вес: %.2f\n",
			criterion.ID(),
			criterion.Name(),
			criterion.MaxScore(),
			criterion.Weight(),
		)
	}

	return nil
}

func (c *Console) readEvaluationData() (string, string, int, error) {
	studentID, err := c.readLine("ID студента: ")
	if err != nil {
		return "", "", 0, err
	}

	criterionID, err := c.readLine("ID критерия: ")
	if err != nil {
		return "", "", 0, err
	}

	score, err := c.readInt("Балл: ")
	if err != nil {
		return "", "", 0, err
	}

	return studentID, criterionID, score, nil
}

func (c *Console) setEvaluation() error {
	studentID, criterionID, score, err := c.readEvaluationData()
	if err != nil {
		return err
	}

	return c.service.SetEvaluation(studentID, criterionID, score)
}

func (c *Console) showEvaluations() error {
	studentID, err := c.readLine("ID студента: ")
	if err != nil {
		return err
	}

	evaluations, err := c.service.Evaluations(studentID)
	if err != nil {
		return err
	}

	if len(evaluations) == 0 {
		fmt.Fprintln(c.writer, "Оценок пока нет.")
		return nil
	}

	for _, evaluation := range evaluations {
		fmt.Fprintf(
			c.writer,
			"Критерий: %s | балл: %d\n",
			evaluation.CriterionID(),
			evaluation.Score(),
		)
	}

	return nil
}

func (c *Console) changeEvaluation() error {
	studentID, criterionID, score, err := c.readEvaluationData()
	if err != nil {
		return err
	}

	return c.service.ChangeEvaluation(studentID, criterionID, score)
}

func (c *Console) deleteEvaluation() error {
	studentID, err := c.readLine("ID студента: ")
	if err != nil {
		return err
	}

	criterionID, err := c.readLine("ID критерия: ")
	if err != nil {
		return err
	}

	return c.service.DeleteEvaluation(studentID, criterionID)
}

func (c *Console) calculateResult() error {
	studentID, err := c.readLine("ID студента: ")
	if err != nil {
		return err
	}

	result, err := c.service.CalculateResult(studentID)
	if err != nil {
		return err
	}

	fmt.Fprintf(
		c.writer,
		"Студент: %s\nИтог: %.2f из 10\n",
		result.StudentName,
		result.Score,
	)

	return nil
}

func (c *Console) exportResults() error {
	path, err := c.readLine("Имя CSV-файла (например, results.csv): ")
	if err != nil {
		return err
	}

	if err := c.service.ExportResults(path); err != nil {
		return err
	}

	fmt.Fprintln(c.writer, "Результаты сохранены:", path)
	return nil
}
