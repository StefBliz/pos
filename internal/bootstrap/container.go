package bootstrap

import (
	"context"
	"errors"
	"pos/internal/adapters/console"
	csvadapter "pos/internal/adapters/csv"
	"pos/internal/adapters/memory"
	"pos/internal/application"
	"pos/internal/domain"
	"time"

	"go.uber.org/fx"
)

// подключает хранилище студентов к порту
func provideStudents() application.StudentRepository {
	return memory.NewStudentRepository()
}

// хранилище рубрикатора
func provideRubric() application.RubricRepository {
	return memory.NewRubricRepository()
}

// хранилище оценок
func provideEvaluations() application.EvaluationRepository {
	return memory.NewEvaluationRepository()
}

// CSV-экспорт
func provideExporter() application.ResultExporter {
	return csvadapter.NewExporter()
}

func provideCalculator() (*domain.ScoreCalculator, error) {
	return domain.NewScoreCalculator(domain.WeightedScoreStrategy{})
}

var Module = fx.Module(
	"pos",
	fx.Provide(
		provideStudents,
		provideRubric,
		provideEvaluations,
		provideExporter,
		provideCalculator,
		application.NewService,
		console.NewConsole,
	),
)

func Run() error {
	var terminal *console.Console

	app := fx.New(
		Module,
		fx.Invoke(func(c *console.Console) {
			terminal = c
		}),
		fx.NopLogger,
	)

	if err := app.Err(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		return err
	}

	if terminal == nil {
		return errors.New("консоль не создана")
	}

	runErr := terminal.Run()

	stopCtx, stopCancel := context.WithTimeout(
		context.Background(), 5*time.Second,
	)
	defer stopCancel()

	stopErr := app.Stop(stopCtx)

	return errors.Join(runErr, stopErr)
}
