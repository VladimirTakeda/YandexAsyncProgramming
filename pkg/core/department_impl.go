package core

import (
	"github.com/paulbellamy/ratecounter"

	"sync"
	"sync/atomic"
	"time"
)

const MaxPercent = 100
const MaxGoroutines = 12

type DepartmentImpl struct {
	Client DeliveryClient
}

func (department *DepartmentImpl) WithDeliveryClient(client DeliveryClient) Department {
	department.Client = client
	return department
}

func (department *DepartmentImpl) Handle(cards []*Card) (Operation, error) {
	o := &OperationImpl{
		progressChan:   make(chan uint8, 1),
		doneChan:       make(chan struct{}),
		stats:          Stats{},
		transportStats: TransportStats{},
		result:         OperationResult{},
		TotalCards:     len(cards),
		Counter:        ratecounter.NewRateCounter(1 * time.Second),
		StartTime:      time.Now(),
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			select {
			case o.progressChan <- uint8(atomic.LoadUint64(&o.stats.TotalCardsHandled) * 100 / uint64(o.TotalCards)):
			case <-o.doneChan:
				close(o.progressChan)
				return
			}
		}
	}()

	go func() {
		var wg sync.WaitGroup
		wg.Add(len(cards))

		workerPool := make(chan struct{}, MaxGoroutines)
		for i := 0; i < MaxGoroutines; i++ {
			workerPool <- struct{}{}
		}

		for _, card := range cards {
			<-workerPool

			go func(card *Card) {
				// Обновляем статистику
				defer func() {
					workerPool <- struct{}{}
					wg.Done()
				}()

				atomic.AddUint64(&o.stats.TotalCardsHandled, 1)

				card.Price, card.EstimatedDeliveryTime = department.Client.GetDeliveryInfo(card.DestinationID, card.CustomerID, card.Weight)
				transportID, err := department.Client.Register(*card)

				if err == nil {
					o.transportStatsMutex.Lock()
					stat, ok := o.transportStats[transportID]
					if !ok {
						stat = &TransportStat{}
						o.transportStats[transportID] = stat
					}

					atomic.AddInt64(&stat.LoadedWeight, card.Weight)
					atomic.AddInt64((*int64)(&stat.DeliveryTime), int64(card.EstimatedDeliveryTime))

					o.transportStatsMutex.Unlock()
				}

				o.Counter.Incr(1)
			}(card)
		}

		wg.Wait()

		o.result.Elapsed = time.Since(o.StartTime)
		o.result.TransportStats = o.transportStats

		atomic.StoreUint64(&o.result.HandledCards, atomic.LoadUint64(&o.stats.TotalCardsHandled))

		o.doneChan <- struct{}{}

		if len(o.progressChan) == 0 {
			o.progressChan <- uint8(MaxPercent)
		}

		close(o.doneChan)
	}()

	return o, nil
}

func init() {
	CurrentDepartment = &DepartmentImpl{}
}
