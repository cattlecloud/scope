package scope

import (
	"context"
	"sync"
	"time"
)

// join implements a context that is canceled when either of two
// contexts is canceled. It is a variation on the original implementation
// with race condition bug fixes.
//
// https://github.com/LK4D4/joincontext/blob/master/context.go
type join struct {
	once sync.Once
	a    C
	b    C
	done chan struct{}

	lock *sync.Mutex
	err  error
}

// Join combines two contexts into a single context that is canceled when
// either input context is canceled. The returned context's Deadline,
// Value, and Err methods delegate to the earliest of the two input
// contexts. If either context is already done, the returned context is
// immediately done with that context's error.
func Join(a, b C) (C, Cancel) {
	j := &join{
		a:    a,
		b:    b,
		done: make(chan struct{}),
		lock: new(sync.Mutex),
	}

	// check if either context is already done before spawning the goroutine
	select {
	case <-a.Done():
		j.lock.Lock()
		j.err = a.Err()
		j.lock.Unlock()

		close(j.done)
		return j, func() {}

	case <-b.Done():
		j.lock.Lock()
		j.err = b.Err()
		j.lock.Unlock()

		close(j.done)
		return j, func() {}

	default:
	}

	go j.run()
	return j, j.cancel
}

// Deadline returns the earliest deadline from either context.
// If neither context has a deadline, ok is false.
func (j *join) Deadline() (deadline time.Time, ok bool) {
	a, aok := j.a.Deadline()
	if !aok {
		return j.b.Deadline()
	}

	b, bok := j.b.Deadline()
	if !bok {
		return a, true
	}

	if b.Before(a) {
		return b, true
	}

	return a, true
}

// Done returns a channel that is closed when either context is done.
func (j *join) Done() <-chan struct{} {
	return j.done
}

// Err returns the error from whichever context was canceled first,
// or ErrCanceled if Cancel was called on the joined context.
func (j *join) Err() error {
	j.lock.Lock()
	defer j.lock.Unlock()
	return j.err
}

// Value returns the value associated with key in either context,
// prioritizing the first context's value if present.
func (j *join) Value(key any) any {
	v := j.a.Value(key)

	if v == nil {
		v = j.b.Value(key)
	}

	return v
}

func (j *join) run() {
	select {
	case <-j.a.Done():
		j.once.Do(func() {
			j.lock.Lock()
			j.err = j.a.Err()
			j.lock.Unlock()
			close(j.done)
		})
	case <-j.b.Done():
		j.once.Do(func() {
			j.lock.Lock()
			j.err = j.b.Err()
			j.lock.Unlock()
			close(j.done)
		})
	case <-j.done:
		return
	}
}

func (j *join) cancel() {
	j.once.Do(func() {
		j.lock.Lock()
		j.err = context.Canceled
		j.lock.Unlock()
		close(j.done)
	})
}
