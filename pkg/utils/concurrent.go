package utils

import (
	"context"
	"sync"
)

type InvokeOptions struct {
	ctx context.Context
}

func (o InvokeOptions) GetContext() context.Context {
	if o.ctx == nil {
		return context.Background()
	}

	return o.ctx
}

func WithContext(ctx context.Context) InvokeOptions {
	return InvokeOptions{
		ctx: ctx,
	}
}

func Invoke0(funcToBeRan func() error, opts ...InvokeOptions) error {
	opt := InvokeOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	resChan := make(chan struct{}, 1)
	var err error

	go func() {
		err = funcToBeRan()
		resChan <- struct{}{}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case <-opt.GetContext().Done():
			err = opt.GetContext().Err()
		case <-resChan:
		}

		wg.Done()
	}()

	wg.Wait()
	return err
}
