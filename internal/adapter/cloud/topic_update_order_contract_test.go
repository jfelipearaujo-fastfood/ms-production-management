package cloud

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-production-management/internal/entity/order_entity"
	"github.com/stretchr/testify/assert"
)

func TestNewContract(t *testing.T) {
	t.Run("Should return a valid contract", func(t *testing.T) {
		// Arrange
		order := order_entity.NewOrder(uuid.NewString(), time.Now())

		// Act
		contract := NewUpdateOrderContractFromPayment(&order)

		// Assert
		assert.Equal(t, order.Id, contract.OrderId)
		assert.Equal(t, order.StateTitle, contract.Order.State)
	})
}
