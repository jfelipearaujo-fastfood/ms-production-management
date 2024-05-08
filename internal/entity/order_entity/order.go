package order_entity

import (
	"time"

	"github.com/jfelipearaujo-org/ms-production-management/internal/shared/custom_error"
)

type Order struct {
	Id string `json:"id"`

	State          OrderState `json:"state"`
	StateTitle     string     `json:"state_title"`
	StateUpdatedAt time.Time  `json:"state_updated_at"`

	Items []Item `json:"items"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewOrder(orderID string, now time.Time) Order {
	return Order{
		Id: orderID,

		State:          Received,
		StateUpdatedAt: now,

		Items: make([]Item, 0),

		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (o *Order) AddItem(item Item, now time.Time) error {
	for _, i := range o.Items {
		if i.Id == item.Id {
			return custom_error.ErrOrderItemAlreadyExists
		}
	}

	o.Items = append(o.Items, item)
	o.UpdatedAt = now

	return nil
}

func (o *Order) UpdateState(toState OrderState, now time.Time) error {
	if o.State == toState {
		return nil
	}

	if !o.State.CanTransitionTo(toState) {
		return custom_error.ErrOrderInvalidStateTransition
	}

	o.State = toState
	o.StateTitle = toState.String()
	o.StateUpdatedAt = now
	o.UpdatedAt = now

	return nil
}

func (o *Order) RefreshStateTitle() {
	o.StateTitle = o.State.String()
}

func (o *Order) IsCompleted() bool {
	return o.State == Delivered || o.State == Cancelled
}

func (o *Order) HasItems() bool {
	return len(o.Items) > 0
}
