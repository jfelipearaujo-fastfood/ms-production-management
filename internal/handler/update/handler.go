package update

import (
	"net/http"

	"github.com/jfelipearaujo-org/ms-production-management/internal/service"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service/order_production/update"
	"github.com/jfelipearaujo-org/ms-production-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service service.UpdateOrderProductionService[update.UpdateOrderProductionInput]
}

func NewHandler(
	service service.UpdateOrderProductionService[update.UpdateOrderProductionInput],
) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Handle(ctx echo.Context) error {
	var request update.UpdateOrderProductionInput

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

	order.RefreshStateTitle()

	return ctx.JSON(http.StatusOK, order)
}
