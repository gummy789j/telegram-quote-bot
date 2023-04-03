package task

import (
	"context"
	"time"
)

type Task interface {
	Name() string
	Freq() (runTime time.Duration, tickTime time.Duration)
	Run(ctx context.Context) error
}
