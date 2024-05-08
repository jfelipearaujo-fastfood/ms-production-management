package get_by_id

import (
	"net/http"

	"github.com/jfelipearaujo-org/ms-production-management/internal/service"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service/order_production/get_by_id"
	"github.com/jfelipearaujo-org/ms-production-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service service.GetOrderProductionByIdService[get_by_id.GetOrderProductionByIdInput]
}

func NewHandler(
	service service.GetOrderProductionByIdService[get_by_id.GetOrderProductionByIdInput],
) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Handle(ctx echo.Context) error {
	var request get_by_id.GetOrderProductionByIdInput

	if err := ctx.Bind(&request); err != nil {
		return err
	}

	context := ctx.Request().Context()

	order, err := h.service.Handle(context, request)
	if err != nil {
		if custom_error.IsBusinessErr(err) {
			return custom_error.NewHttpAppErrorFromBusinessError(err)
		}

		return custom_error.NewHttpAppError(http.StatusInternalServerError, "internal server error", err)
	}

	return ctx.JSON(http.StatusOK, order)
}
