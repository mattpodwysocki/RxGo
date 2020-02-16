package rxgo

import "context"

// Option handles configurable options.
type Option interface {
	apply(*funcOption)
	withBuffer() (bool, int)
	withContext() (bool, context.Context)
	withEagerObservation() bool
	withPool() (bool, int)
	buildChannel() chan Item
	buildContext() context.Context
}

type funcOption struct {
	f                func(*funcOption)
	toBuffer         bool
	buffer           int
	ctx              context.Context
	eagerObservation bool
	pool             int
}

func (fdo *funcOption) withBuffer() (bool, int) {
	return fdo.toBuffer, fdo.buffer
}

func (fdo *funcOption) withContext() (bool, context.Context) {
	return fdo.ctx != nil, fdo.ctx
}

func (fdo *funcOption) withEagerObservation() bool {
	return fdo.eagerObservation
}

func (fdo *funcOption) withPool() (bool, int) {
	return fdo.pool > 0, fdo.pool
}

func (fdo *funcOption) buildChannel() chan Item {
	if toBeBuffered, cap := fdo.withBuffer(); toBeBuffered {
		return make(chan Item, cap)
	}
	return make(chan Item)
}

func (fdo *funcOption) buildContext() context.Context {
	if withContext, c := fdo.withContext(); withContext {
		return c
	}
	return context.Background()
}

func (fdo *funcOption) apply(do *funcOption) {
	fdo.f(do)
}

func newFuncOption(f func(*funcOption)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func parseOptions(opts ...Option) Option {
	o := new(funcOption)
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

// WithBufferedChannel allows to configure the capacity of a buffered channel.
func WithBufferedChannel(capacity int) Option {
	return newFuncOption(func(options *funcOption) {
		options.toBuffer = true
		options.buffer = capacity
	})
}

// WithContext allows to pass a context.
func WithContext(ctx context.Context) Option {
	return newFuncOption(func(options *funcOption) {
		options.ctx = ctx
	})
}

// WithEagerObservation uses the eager observation mode meaning consuming the items even without subscription.
func WithEagerObservation() Option {
	return newFuncOption(func(options *funcOption) {
		options.eagerObservation = true
	})
}

// WithPool allows to specify an execution pool.
func WithPool(pool int) Option {
	return newFuncOption(func(options *funcOption) {
		options.pool = pool
	})
}
