package scheduler

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-currency-watcher-bot/internal/infra/tracing"
)

type TaskFn func(ctx context.Context) error

func PeriodicTaskExecutor(
	ctx context.Context,
	interval time.Duration,
	taskName string,
	task TaskFn,
) {
	log := logger.FromContext(ctx).WithField("executor_task_name", taskName)
	log.Info("start")
	defer log.Info("end")

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	executeTask(ctx, log, taskName, task)

	for {
		select {
		case <-ticker.C:
			executeTask(ctx, log, taskName, task)

		case <-ctx.Done():
			return
		}
	}
}

func executeTask(ctx context.Context, log logger.Logger, taskName string, task TaskFn) {
	ctx, span := tracing.StartSpanName(ctx, taskName)
	defer span.End()

	log.Info("execute task")

	if err := task(ctx); err != nil {
		log.WithError(err).Error("execute task error")
	}
}
