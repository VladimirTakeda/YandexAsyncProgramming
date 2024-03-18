package core

import "time"

type DeliveryClient interface {
	// GetDeliveryInfo calculates price and estimated delivery duration
	GetDeliveryInfo(destinationID, customerID int64, weight int64) (price int64, deliveryDuration time.Duration)
	// Register the card and assign transport to it.
	// Returns an error if card already registered.
	Register(card Card) (transportID int64, err error)
}
