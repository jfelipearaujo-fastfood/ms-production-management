package update

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service/mocks"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service/order_production/update"
	"github.com/jfelipearaujo-org/ms-production-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should update the order", func(t *testing.T) {
		// Arrange
		service := mocks.NewMockUpdateOrderProductionService[update.UpdateOrderProductionInput](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(&order_entity.Order{}, nil).
			Once()

		reqBody := update.UpdateOrderProductionInput{
			OrderId: uuid.NewString(),
			State:   "Processing",
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/production/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(reqBody.OrderId)

		handler := NewHandler(service)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		service.AssertExpectations(t)
	})

	t.Run("Should return validation error", func(t *testing.T) {
		// Arrange
		service := mocks.NewMockUpdateOrderProductionService[update.UpdateOrderProductionInput](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(nil, custom_error.ErrOrderAlreadyAtState).
			Once()

		reqBody := update.UpdateOrderProductionInput{
			OrderId: uuid.NewString(),
			State:   "Processing",
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/production/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(reqBody.OrderId)

		handler := NewHandler(service)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.Error(t, err)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)

		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, custom_error.AppError{
			Code:    http.StatusBadRequest,
			Message: "unable to update order state",
			Details: "order is already at the state",
		}, he.Message)

		service.AssertExpectations(t)
	})

	t.Run("Should return internal server error", func(t *testing.T) {
		// Arrange
		service := mocks.NewMockUpdateOrderProductionService[update.UpdateOrderProductionInput](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(nil, assert.AnError).
			Once()

		reqBody := update.UpdateOrderProductionInput{
			OrderId: uuid.NewString(),
			State:   "Processing",
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/production/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(reqBody.OrderId)

		handler := NewHandler(service)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.Error(t, err)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)

		assert.Equal(t, http.StatusInternalServerError, he.Code)
		assert.Equal(t, custom_error.AppError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Details: "assert.AnError general error for testing",
		}, he.Message)

		service.AssertExpectations(t)
	})
}
