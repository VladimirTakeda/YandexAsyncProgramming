package core

import (
	"time"
)

type OperationResult struct {
	// HandledCards is a total amount of cards handled during this
	// operation
	HandledCards uint64
	// Elapsed shows total time spent by this operation.
	Elapsed time.Duration
	// TransportLoad shows aggregated information about each transport.
	TransportStats TransportStats
}

type Operation interface {
	// ProgressChan returns a channel for reporting current progress of handling
	// operation. It's expected to receive progress report each second.
	ProgressChan() <-chan uint8
	// Done report whether the operation
	Done() chan struct{}
	// Stats returns a snapshot of the current stats.
	Stats() (*Stats, error)
	// TransportStats returns a snapshot of aggregated transport stats.
	TransportStats() (TransportStats, error)

	Result() *OperationResult
}
