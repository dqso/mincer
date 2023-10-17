package shutdown

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCloser(t *testing.T) {
	testCases := []struct {
		name            string
		closeFuncs      []CloseFunc
		shutdownTimeout time.Duration
		wantError       string
		wantErrorNil    bool
	}{
		{
			name:            "empty",
			closeFuncs:      []CloseFunc{},
			shutdownTimeout: time.Millisecond,
			wantErrorNil:    true,
		},
		{
			name: "finished without errors",
			closeFuncs: []CloseFunc{
				func(ctx context.Context) error { return sleepNil(time.Millisecond) },
				func(ctx context.Context) error { return sleepNil(time.Millisecond) },
			},
			shutdownTimeout: time.Millisecond * 3,
			wantErrorNil:    true,
		},
		{
			name: "cancelled without errors",
			closeFuncs: []CloseFunc{
				func(ctx context.Context) error { return sleepNil(time.Millisecond) },
				func(ctx context.Context) error { return sleepNil(time.Millisecond) },
			},
			shutdownTimeout: time.Millisecond,
			wantError:       "shutdown cancelled",
		},
		{
			name: "finished with errors",
			closeFuncs: []CloseFunc{
				func(ctx context.Context) error { return sleepNil(time.Millisecond) },
				func(ctx context.Context) error { return sleepError(time.Millisecond, fmt.Errorf("1")) },
				func(ctx context.Context) error { return sleepError(time.Millisecond, fmt.Errorf("2")) },
			},
			shutdownTimeout: time.Millisecond * 4,
			wantError:       "shutdown finished with error(s):\n2\n1",
		},
		{
			name: "cancelled with errors",
			closeFuncs: []CloseFunc{
				func(ctx context.Context) error { return sleepNil(time.Millisecond) },
				func(ctx context.Context) error { return sleepNil(time.Millisecond) },
				func(ctx context.Context) error { return sleepError(time.Millisecond, fmt.Errorf("1")) },
				func(ctx context.Context) error { return sleepError(time.Millisecond, fmt.Errorf("2")) },
			},
			shutdownTimeout: time.Millisecond * 3,
			wantError:       "shutdown cancelled with error(s):\n2\n1",
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var closer Closer
			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), test.shutdownTimeout)
				defer cancel()
				err := closer.Close(ctx)
				if test.wantErrorNil {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, test.wantError)
				}
			}()

			for _, fn := range test.closeFuncs {
				closer.Add(fn)
			}

			cancel()
			<-ctx.Done()
		})
	}
}

func sleepNil(duration time.Duration) error {
	time.Sleep(duration)
	return nil
}

func sleepError(duration time.Duration, err error) error {
	time.Sleep(duration)
	return err
}
