package executor_test

import (
	"context"
	"time"

	"github.com/Mereng/executor"
)

type sampleJob struct {
}

func (j *sampleJob) Execute(ctx context.Context) {
	time.Sleep(200 * time.Millisecond)
}

func ExampleExecutor() {
	ch := make(chan executor.Job, 3)

	ctx, cancel := context.WithCancel(context.Background())
	ex := executor.New(ctx, ch, 2, 100*time.Millisecond, 2)

	for i := 0; i < 3; i++ {
		ch <- &sampleJob{}
	}

	cancel()
	ex.Wait()
}
