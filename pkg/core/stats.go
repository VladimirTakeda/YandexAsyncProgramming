package core

import (
	"time"
)

type Stats struct {
	TotalCardsHandled uint64
	CardsPerSecond    uint64
}

type TransportStat struct {
	LoadedWeight int64
	DeliveryTime time.Duration
}

type TransportStats map[int64]*TransportStat
