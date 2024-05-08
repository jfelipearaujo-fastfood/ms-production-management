package get_by_state

import (
	"github.com/go-playground/validator/v10"
	"github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-production-management/internal/shared/custom_error"
)

type GetOrderProductionByStateInput struct {
	State string `json:"state" validate:"required"`
}

func (input *GetOrderProductionByStateInput) Validate() error {
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return custom_error.ErrRequestNotValid
	}

	if order_entity.NewOrderState(input.State) == order_entity.None {
		return custom_error.ErrRequestNotValid
	}

	return nil
}
