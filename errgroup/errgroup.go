package errgroup

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// BoundedGroup wraps an errgroup to be concurrently bounded. The behaviour is
// the same, the only addition is that the goroutines are not started until they
// acquire a semaphore.
type BoundedGroup struct {
	// The errgroup is wrapped to be bounded by a semaphore.
	eg *errgroup.Group

	// Context created by the underlying errgroup.
	ctx context.Context

	// Semaphore channel. Write to acquire, read to release. The channel will
	// *never* be closed to avoid panic. The channel will be garbage collected
	// by the runtime since none of the goroutine will block to release the
	// semaphore.
	sem chan struct{}
}

// NewBoundedGroup creates a new errgroup that is concurrently bounded by N.
func NewBoundedGroup(ctx context.Context, n uint64) (*BoundedGroup, context.Context) {
	eg, ctx := errgroup.WithContext(ctx)
	sem := make(chan struct{}, n)
	return &BoundedGroup{eg, ctx, sem}, ctx
}

// Go calls the function f in a new goroutine after acquiring a semaphore. If the
// goroutine hasn't started yet and the underlying errgroup's context is
// cancelled, the function f will never be invoked.
func (g *BoundedGroup) Go(f func() error) {
	select {
	case <-g.ctx.Done():
		return
	case g.sem <- struct{}{}:
		// There is a chance that we acquired a semaphore but the ctx is done.
		if g.ctx.Err() != nil {
			return
		}
	}

	g.eg.Go(func() error {
		err := f()
		<-g.sem
		return err
	})
}

// Wait blocks until all functions calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *BoundedGroup) Wait() error {
	return g.eg.Wait()
}
