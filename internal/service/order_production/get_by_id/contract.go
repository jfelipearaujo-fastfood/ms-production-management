package get_by_id

import (
	"github.com/go-playground/validator/v10"
	"github.com/jfelipearaujo-org/ms-production-management/internal/shared/custom_error"
)

type GetOrderProductionByIdInput struct {
	OrderId string `param:"id" json:"order_id" validate:"required,uuid4"`
}

func (input *GetOrderProductionByIdInput) Validate() error {
	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return custom_error.ErrRequestNotValid
	}

	return nil
}
