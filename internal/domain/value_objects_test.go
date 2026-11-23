package domain

import "testing"

func TestScore(t *testing.T) {
	tests := []struct {
		name    string
		score   int
		max     int
		wantErr bool
	}{
		{"корректный балл", 8, 10, false},
		{"нулевой балл", 0, 10, false},
		{"максимальный балл", 10, 10, false},
		{"отрицательный балл", -1, 10, true},
		{"превышение максимума", 11, 10, true},
		{"некорректный максимум", 5, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, err := NewScore(tt.score, tt.max)

			if (err != nil) != tt.wantErr {
				t.Fatalf("неожиданный результат: %v", err)
			}

			if err == nil && score.Value() != tt.score {
				t.Fatal("балл сохранён неправильно")
			}
		})
	}
}

func TestWeight(t *testing.T) {
	tests := []struct {
		name    string
		weight  float64
		wantErr bool
	}{
		{"корректный вес", 0.4, false},
		{"полный вес", 1, false},
		{"нулевой вес", 0, true},
		{"отрицательный вес", -0.2, true},
		{"вес больше единицы", 1.2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewWeight(tt.weight)

			if (err != nil) != tt.wantErr {
				t.Fatalf("неожиданный результат: %v", err)
			}
		})
	}
}

func TestRubricBuilder(t *testing.T) {
	rubric, err := NewRubricBuilder().
		Add("CR-001", "Домашние задания", 10, 0.4).
		Add("CR-002", "Контрольная", 10, 0.3).
		Add("CR-003", "Экзамен", 10, 0.3).
		Build()

	if err != nil {
		t.Fatal(err)
	}

	if len(rubric.Criteria()) != 3 {
		t.Fatal("ожидалось три критерия")
	}
}

func TestRubricBuilderInvalid(t *testing.T) {
	_, err := NewRubricBuilder().
		Add("CR-001", "Экзамен", 10, 0.4).
		Build()

	if err == nil {
		t.Fatal("ожидалась ошибка неполного рубрикатора")
	}
}
