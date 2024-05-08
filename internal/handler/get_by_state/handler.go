package get_by_state

import (
	"net/http"

	"github.com/jfelipearaujo-org/ms-production-management/internal/service"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service/order_production/get_by_state"
	"github.com/jfelipearaujo-org/ms-production-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service service.GetOrderProductionByStateService[get_by_state.GetOrderProductionByStateInput]
}

func NewHandler(
	service service.GetOrderProductionByStateService[get_by_state.GetOrderProductionByStateInput],
) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Handle(ctx echo.Context) error {
	var request get_by_state.GetOrderProductionByStateInput

	if err := ctx.Bind(&request); err != nil {
		return err
	}

	context := ctx.Request().Context()

	orders, err := h.service.Handle(context, request)
	if err != nil {
		if custom_error.IsBusinessErr(err) {
			return custom_error.NewHttpAppErrorFromBusinessError(err)
		}

		return custom_error.NewHttpAppError(http.StatusInternalServerError, "internal server error", err)
	}

	return ctx.JSON(http.StatusOK, orders)
}
