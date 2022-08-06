// executor provides simple rate limiter
package executor

import (
	"context"
	"sync"
	"time"
)

// An implementation of Job can be executed by Executor
type Job interface {
	// Execute will be executed by Executor
	Execute(ctx context.Context)
}

// Executor implements rate limiter
type Executor struct {
	ctx         context.Context
	jobs        <-chan Job
	concurrency int
	rateLimit   int

	mainCh chan Job
	wg     *sync.WaitGroup
	ticker *time.Ticker

	countJobs int
}

// New creates Executor and runs it
//
// ctx is context, that can use for shutdown Executor.
// jobs is channel Executor listens it and executes jobs, if it closes then Executor will shutdown.
// concurrency is the maximum number of concurrent jobs executor may process.
// period of a limit refresh, after period executor refreshes rate.
// rateLimit sets maximum number executed jobs per period.
func New(ctx context.Context, jobs <-chan Job, concurrency int, period time.Duration, rateLimit int) *Executor {
	ex := &Executor{
		ctx:         ctx,
		jobs:        jobs,
		concurrency: concurrency,
		rateLimit:   rateLimit,
		mainCh:      make(chan Job),
		wg:          &sync.WaitGroup{},
		ticker:      time.NewTicker(period),
	}

	go ex.run()

	return ex
}

// Wait waits for Executor to shutdown
func (ex *Executor) Wait() {
	ex.wg.Wait()
}

func (ex *Executor) run() {
	for i := 0; i < ex.concurrency; i++ {
		ex.wg.Add(1)
		go ex.startWorker()
	}

loop:
	for {
		select {
		case j, ok := <-ex.jobs:
			if !ok {
				break loop
			}

			if ex.countJobs >= ex.rateLimit {
				<-ex.ticker.C
				ex.countJobs = 0
			}
			ex.mainCh <- j
			ex.countJobs++
		case <-ex.ctx.Done():
			break loop
		}
	}

	close(ex.mainCh)
}

func (ex *Executor) startWorker() {
	defer ex.wg.Done()

	for j := range ex.mainCh {
		j.Execute(ex.ctx)
	}
}
