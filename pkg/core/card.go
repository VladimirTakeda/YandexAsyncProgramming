package core

import "time"

type Card struct {
	ID                    string
	DestinationID         int64
	CustomerID            int64
	Weight                int64
	Price                 int64
	EstimatedDeliveryTime time.Duration
}
