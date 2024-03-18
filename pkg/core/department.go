package core

import "fmt"

var CurrentDepartment Department = nil

type Department interface {
	// WithDeliveryClient sets delivery client.
	WithDeliveryClient(client DeliveryClient) Department
	// Handle spawns an operation that serves cards. Should not block
	// the execution.
	Handle(cards []*Card) (Operation, error)
}

var errNotImpl = fmt.Errorf("not implemented")

type depNOOP struct{}

var _ Department = (*depNOOP)(nil)

func (d depNOOP) WithDeliveryClient(client DeliveryClient) Department {
	return d
}
func (depNOOP) Handle(cards []*Card) (Operation, error) {
	return nil, errNotImpl
}
func (depNOOP) Stats() (*Stats, error) {
	return nil, errNotImpl
}
func (depNOOP) TransportStats() (TransportStats, error) {
	return nil, errNotImpl
}
