package get_by_id

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service/mocks"
	"github.com/jfelipearaujo-org/ms-production-management/internal/service/order_production/get_by_id"
	"github.com/jfelipearaujo-org/ms-production-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should return the order by id", func(t *testing.T) {
		// Arrange
		service := mocks.NewMockGetOrderProductionByIdService[get_by_id.GetOrderProductionByIdInput](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(order_entity.Order{
				Id: uuid.NewString(),
			}, nil).
			Once()

		reqBody := get_by_id.GetOrderProductionByIdInput{
			OrderId: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/orders/:id/items")
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.NewString())

		handler := NewHandler(service)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		service.AssertExpectations(t)
	})

	t.Run("Should return not found error", func(t *testing.T) {
		// Arrange
		service := mocks.NewMockGetOrderProductionByIdService[get_by_id.GetOrderProductionByIdInput](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(order_entity.Order{}, custom_error.ErrOrderNotFound).
			Once()

		reqBody := get_by_id.GetOrderProductionByIdInput{
			OrderId: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/orders/:id/items")
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.NewString())

		handler := NewHandler(service)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.Error(t, err)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)

		assert.Equal(t, http.StatusNotFound, he.Code)
		assert.Equal(t, custom_error.AppError{
			Code:    http.StatusNotFound,
			Message: "unable to find the order",
			Details: "order not found",
		}, he.Message)

		service.AssertExpectations(t)
	})

	t.Run("Should return internal server error", func(t *testing.T) {
		// Arrange
		service := mocks.NewMockGetOrderProductionByIdService[get_by_id.GetOrderProductionByIdInput](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(order_entity.Order{}, assert.AnError).
			Once()

		reqBody := get_by_id.GetOrderProductionByIdInput{
			OrderId: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/orders/:id/items")
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.NewString())

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
