package core_test

import (
	"testing"

	"git.yandex-academy.ru/homework/async_programming/pkg/core"
	"github.com/stretchr/testify/require"
)

// TestDeparment is a simple test that checks if basic functionality
// works well.
func TestDeparment(t *testing.T) {
	deliveryClient := newDeliveryClient(t)
	department := core.CurrentDepartment.WithDeliveryClient(deliveryClient)
	cards := generateCards(t)

	assert := require.New(t)

	operation, err := department.Handle(cards)
	assert.NoError(err)

	progress := <-operation.ProgressChan()
	assert.NotEmpty(progress)

	<-operation.Done()

	assert.EqualValues(len(cards), operation.Result().HandledCards)
}
