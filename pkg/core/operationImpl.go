package core

import (
	"github.com/paulbellamy/ratecounter"

	"sync"
	"sync/atomic"
	"time"
)

type OperationImpl struct {
	progressChan        chan uint8
	doneChan            chan struct{}
	stats               Stats
	transportStats      TransportStats
	result              OperationResult
	TotalCards          int
	Counter             *ratecounter.RateCounter
	StartTime           time.Time
	transportStatsMutex sync.Mutex
}

func (operation *OperationImpl) ProgressChan() <-chan uint8 {
	return operation.progressChan
}

func (operation *OperationImpl) Done() chan struct{} {
	return operation.doneChan
}

func (operation *OperationImpl) Stats() (*Stats, error) {
	stat := Stats{atomic.LoadUint64(&operation.stats.TotalCardsHandled), uint64(operation.Counter.Rate())}
	return &stat, nil
}

func (operation *OperationImpl) TransportStats() (TransportStats, error) {
	var transportStats TransportStats = make(map[int64]*TransportStat)
	operation.transportStatsMutex.Lock()
	for key, value := range operation.transportStats {
		transportStats[key] = &TransportStat{LoadedWeight: value.LoadedWeight, DeliveryTime: value.DeliveryTime}
	}
	operation.transportStatsMutex.Unlock()
	return transportStats, nil
}

func (operation *OperationImpl) Result() *OperationResult {
	return &operation.result
}
