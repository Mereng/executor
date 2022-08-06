# executor
provides simple rate limiter

## Example

```go
package executor_test

import (
	"context"
	"time"

	"github.com/Mereng/executor"
)

type sampleJob struct {
}

// Implements interface executor.Job
func (j *sampleJob) Execute(ctx context.Context) {
	time.Sleep(200 * time.Millisecond)
}

func main() {
	ch := make(chan executor.Job, 3)

	ctx, cancel := context.WithCancel(context.Background())
	ex := executor.New(ctx, ch, 2, 1*time.Minute, 2)

	for i := 0; i < 3; i++ {
		ch <- &sampleJob{}
	}

	cancel()
	ex.Wait()
}
```