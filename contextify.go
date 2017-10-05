// Package contextify provides a utility function to convert a context-unaware
// run function and a cancel function to a context-aware function.
package contextify

import "context"

// Contextify convert a context-unaware run function and a cancel function to
// a context-aware function.
//
// If context.Context is not cancelled before run() finishes, the return value
// function waits for run() to be finished and returns the return value of run().
//
// If context.Context is cancelled before run() finishes, the return value
// function waits for both the run() and cancel() to be finished.
//
// If pickError is nil, the first non-nil error will be retruned of the return
// value of run(), cancel(), and context.Context.Err().
// You can change this behavior with writing a function to pick a desired error
// and pass it to the pickError argument.
func Contextify(run func() error, cancel func() error,
	pickError func(errFromRun, errFromCancel, errFromContext error) error) func(context.Context) error {

	return func(ctx context.Context) error {
		var errFromRun error
		done := make(chan struct{})
		go func() {
			errFromRun = run()
			close(done)
		}()

		select {
		case <-done:
			return errFromRun
		case <-ctx.Done():
			errFromCancel := cancel()
			<-done
			if pickError == nil {
				pickError = defaultPickError
			}
			return pickError(errFromRun, errFromCancel, ctx.Err())
		}
	}
}

func defaultPickError(errFromRun, errFromCancel, errFromContext error) error {
	if errFromRun != nil {
		return errFromRun
	}
	if errFromCancel != nil {
		return errFromCancel
	}
	return errFromContext
}
