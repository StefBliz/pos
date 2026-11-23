package csv

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"

	"pos/internal/application"
)

func TestExport(t *testing.T) {
	exporter := NewExporter()
	path := filepath.Join(t.TempDir(), "results.csv")

	results := []application.Result{
		{StudentID: "ST-001", StudentName: "Анна", Score: 8.9},
		{StudentID: "ST-002", StudentName: "Мария", Score: 7.5},
	}

	if err := exporter.Export(results, path); err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 3 {
		t.Fatalf("ожидалось 3 строки, получено %d", len(records))
	}

	if records[1][0] != "ST-001" ||
		records[1][1] != "Анна" ||
		records[1][2] != "8.90" {
		t.Fatal("неверные данные первого студента")
	}

	if records[2][2] != "7.50" {
		t.Fatal("неверный итог второго студента")
	}
}

func TestExportInvalidPath(t *testing.T) {
	exporter := NewExporter()

	err := exporter.Export(
		[]application.Result{},
		filepath.Join(t.TempDir(), "missing", "results.csv"),
	)

	if err == nil {
		t.Fatal("ожидалась ошибка создания файла")
	}
}
