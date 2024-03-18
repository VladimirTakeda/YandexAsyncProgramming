package core_test

import (
	"strconv"
	"testing"
	"time"

	"git.yandex-academy.ru/homework/async_programming/pkg/core"
)

type testError string

func (err testError) Error() string { return string(err) }

const (
	ErrBadCard testError = "card price or weight is empty"
)

type simpleTestClient struct {
	tb testing.TB
}

func newDeliveryClient(tb testing.TB) *simpleTestClient {
	return &simpleTestClient{
		tb: tb,
	}
}

func (s *simpleTestClient) GetDeliveryInfo(destinationID, customerID int64, weight int64) (price int64, deliveryDuration time.Duration) {
	s.tb.Logf("GetDeliveryInfo called with params destinationID=%d customerID=%d weight=%d", destinationID, customerID, weight)

	return 1, time.Second
}

func (s *simpleTestClient) Register(card core.Card) (transportID int64, err error) {
	s.tb.Logf("Register called with params card=%#v", card)

	if card.Price == 0 || card.Weight == 0 {
		return 0, ErrBadCard
	}

	return 1, nil
}

func generateCards(tb testing.TB) []*core.Card {
	const generateCardsNum = 10
	out := make([]*core.Card, generateCardsNum)
	for i := 0; i < generateCardsNum; i++ {
		card := core.Card{
			ID:                    strconv.Itoa(i + 1),
			DestinationID:         1,
			CustomerID:            1,
			Weight:                int64(10 * (i + 1)),
			Price:                 0,
			EstimatedDeliveryTime: 0,
		}

		tb.Logf("generated card %#v", card)

		out[i] = &card
	}

	return out
}
