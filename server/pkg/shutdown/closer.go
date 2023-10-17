package shutdown

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Closer struct {
	funcs []CloseFunc
	mx    sync.Mutex
}

type CloseFunc func(ctx context.Context) error

func (c *Closer) Add(fn CloseFunc) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.funcs = append(c.funcs, fn)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	chErrors := make(chan error, 1)
	go func() {
		defer close(chErrors)
		for i := len(c.funcs) - 1; i >= 0; i-- {
			err := c.funcs[i](ctx)
			if err != nil {
				chErrors <- err
			}
			c.funcs[i] = nil
			c.funcs = c.funcs[0:len(c.funcs)]
		}
	}()

	var messages []string
	const (
		asFinished  = "finished"
		asCancelled = "cancelled"
	)
	returnError := func(status string) error {
		if len(messages) > 0 {
			return fmt.Errorf("shutdown %s with error(s):\n%s", status, strings.Join(messages, "\n"))
		}
		if status != asFinished {
			return fmt.Errorf("shutdown %s", status)
		}
		return nil
	}

loop:
	for {
		select {
		case err, ok := <-chErrors:
			if !ok {
				break loop
			}
			messages = append(messages, err.Error())
		case <-ctx.Done():
			return returnError(asCancelled)
		}
	}

	return returnError(asFinished)
}
