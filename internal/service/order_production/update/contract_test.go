package update

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("Should return nil when valid", func(t *testing.T) {
		// Arrange
		input := UpdateOrderProductionInput{
			OrderId: uuid.NewString(),
			State:   "Received",
		}

		// Act
		err := input.Validate()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return error when invalid", func(t *testing.T) {
		// Arrange
		input := UpdateOrderProductionInput{
			OrderId: "123",
			State:   "Received",
		}

		// Act
		err := input.Validate()

		// Assert
		assert.Error(t, err)
	})

	t.Run("Should return error when state is invalid", func(t *testing.T) {
		// Arrange
		input := UpdateOrderProductionInput{
			OrderId: uuid.NewString(),
			State:   "invalid",
		}

		// Act
		err := input.Validate()

		// Assert
		assert.Error(t, err)
	})
}
