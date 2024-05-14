package waiter

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type waiterConfig struct {
	parentContext context.Context
	catchSignals  bool
}

type WaiterOption func(c *waiterConfig)

func ParentContext(ctx context.Context) WaiterOption {
	return func(c *waiterConfig) {
		c.parentContext = ctx
	}
}

func CatchSignals() WaiterOption {
	return func(c *waiterConfig) {
		c.catchSignals = true
	}
}

type WaitFunc func(ctx context.Context) error

type Waiter interface {
	Add(fns ...WaitFunc)
	Wait() error
	Context() context.Context
	CancelFunc() context.CancelFunc
}

type waiter struct {
	ctx    context.Context
	fns    []WaitFunc
	cancel context.CancelFunc
}

func New(options ...WaiterOption) Waiter {
	cfg := &waiterConfig{
		parentContext: context.Background(),
		catchSignals:  false,
	}

	for _, option := range options {
		option(cfg)
	}

	ctx, cancel := context.WithCancel(cfg.parentContext)

	w := &waiter{
		ctx:    ctx,
		fns:    []WaitFunc{},
		cancel: cancel,
	}

	if cfg.catchSignals {
		w.ctx, w.cancel = signal.NotifyContext(w.ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	}

	return w
}

func (w *waiter) Add(fns ...WaitFunc) {
	w.fns = append(w.fns, fns...)
}

func (w *waiter) Wait() error {
	group, ctx := errgroup.WithContext(w.ctx)

	group.Go(func() error {
		<-ctx.Done()
		w.cancel()
		return nil
	})

	for _, fn := range w.fns {
		f := fn
		group.Go(func() error {
			return f(ctx)
		})
	}

	return group.Wait()
}

func (w *waiter) Context() context.Context {
	return w.ctx
}

func (w *waiter) CancelFunc() context.CancelFunc {
	return w.cancel
}
