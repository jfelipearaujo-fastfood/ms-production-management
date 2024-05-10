package update

import (
	"log/slog"
	"net/http"

	"github.com/jfelipearaujo-org/ms-production-management/internal/adapter/cloud"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service/order_production/update"
	"github.com/jfelipearaujo-org/ms-production-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	updateOrderProductionService service.UpdateOrderProductionService[update.UpdateOrderProductionInput]
	updateOrderTopic             cloud.TopicService
}

func NewHandler(
	updateOrderProductionService service.UpdateOrderProductionService[update.UpdateOrderProductionInput],
	updateOrderTopic cloud.TopicService,
) *Handler {
	return &Handler{
		updateOrderProductionService: updateOrderProductionService,
		updateOrderTopic:             updateOrderTopic,
	}
}

func (h *Handler) Handle(c echo.Context) error {
	var request update.UpdateOrderProductionInput

	if err := c.Bind(&request); err != nil {
		return err
	}

	ctx := c.Request().Context()

	order, err := h.updateOrderProductionService.Handle(ctx, request)
	if err != nil {
		if custom_error.IsBusinessErr(err) {
			return custom_error.NewHttpAppErrorFromBusinessError(err)
		}

		return custom_error.NewHttpAppError(http.StatusInternalServerError, "internal server error", err)
	}

	order.RefreshStateTitle()

	messageId, err := h.updateOrderTopic.PublishMessage(ctx, cloud.NewUpdateOrderContractFromPayment(order))
	if err != nil {
		slog.ErrorContext(ctx, "error publishing message to update order topic", "error", err)
	}

	if messageId != nil {
		slog.InfoContext(ctx, "message published to update order topic", "message_id", *messageId)
	}

	return c.JSON(http.StatusOK, order)
}
