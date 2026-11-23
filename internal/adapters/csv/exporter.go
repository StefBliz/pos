package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"pos/internal/application"
)

type Exporter struct{}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) Export(
	results []application.Result,
	path string,
) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	writer := csv.NewWriter(file)

	if err = writer.Write([]string{
		"ID студента",
		"Имя студента",
		"Итоговый балл",
	}); err != nil {
		return err
	}

	for _, result := range results {
		record := []string{
			result.StudentID,
			result.StudentName,
			fmt.Sprintf("%.2f", result.Score),
		}

		if err = writer.Write(record); err != nil {
			return err
		}
	}

	writer.Flush()

	if err = writer.Error(); err != nil {
		return fmt.Errorf("ошибка записи CSV: %w", err)
	}

	return nil
}
