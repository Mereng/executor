package executor

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testJob struct {
	exec func()
}

func (j *testJob) Execute(ctx context.Context) {
	j.exec()
}

func TestExecutor(t *testing.T) {
	executed := int32(0)
	jobs := []*testJob{
		{func() {
			time.Sleep(20 * time.Millisecond)
			atomic.AddInt32(&executed, 1)
		}},
		{func() {
			time.Sleep(20 * time.Millisecond)
			atomic.AddInt32(&executed, 1)
		}},
		{func() {
			time.Sleep(20 * time.Millisecond)
			atomic.AddInt32(&executed, 1)
		}},
	}

	ch := make(chan Job, 3)
	ex := New(context.Background(), ch, 2, 100*time.Millisecond, 2)

	for _, j := range jobs {
		ch <- j
	}

	time.Sleep(30 * time.Millisecond)
	assert.Equal(t, int32(2), executed)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, int32(3), executed)
	ch <- jobs[0]
	time.Sleep(30 * time.Millisecond)
	assert.Equal(t, int32(4), executed)

	close(ch)
	ex.Wait()
}
