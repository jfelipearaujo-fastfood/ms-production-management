package create

import (
	"github.com/go-playground/validator/v10"
	"github.com/jfelipearaujo-org/ms-production-management/internal/shared/custom_error"
)

type CreateOrderProductionItemInput struct {
	Id       string `json:"id" validate:"required,uuid4"`
	Name     string `json:"name" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,gte=1"`
}

type CreateOrderProductionInput struct {
	OrderId string `json:"order_id" validate:"required,uuid4"`

	Items []CreateOrderProductionItemInput `json:"items" validate:"required,dive"`
}

func (input *CreateOrderProductionInput) Validate() error {
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return custom_error.ErrRequestNotValid
	}

	return nil
}
