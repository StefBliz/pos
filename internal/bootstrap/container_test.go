package bootstrap

import (
	"testing"

	"go.uber.org/fx"

	"pos/internal/application"
)

// Проверяем, что Uber Fx может собрать все зависимости.
func TestContainer(t *testing.T) {
	var service *application.Service

	app := fx.New(
		Module,
		fx.Invoke(func(s *application.Service) {
			service = s
		}),
		fx.NopLogger,
	)

	if err := app.Err(); err != nil {
		t.Fatal(err)
	}

	if service == nil {
		t.Fatal("DI-контейнер не создал сервис")
	}
}
