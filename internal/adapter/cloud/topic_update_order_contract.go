package cloud

import "github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"

type UpdateOrderTopicOrderContract struct {
	State string `json:"state"`
}

type UpdateOrderTopicContract struct {
	OrderId string                        `json:"order_id"`
	Order   UpdateOrderTopicOrderContract `json:"order"`
}

func NewUpdateOrderContractFromPayment(order *order_entity.Order) *UpdateOrderTopicContract {
	order.RefreshStateTitle()

	return &UpdateOrderTopicContract{
		OrderId: order.Id,
		Order: UpdateOrderTopicOrderContract{
			State: order.StateTitle,
		},
	}
}
