package domain

import (
	"errors"
	"strings"
)

type Student struct {
	id   string
	name string
}

// создание студента и проверка данных
func NewStudent(id, name string) (Student, error) {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	if id == "" {
		return Student{}, errors.New("student ID cannot be empty")
	}
	if name == "" {
		return Student{}, errors.New("student name cannot be empty")
	}
	return Student{
		id:   id,
		name: name,
	}, nil
}

func (s Student) ID() string {
	return s.id
}

func (s Student) Name() string {
	return s.name
}

func (s *Student) Rename(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("student name cannot be empty")
	}
	s.name = name
	return nil
}
