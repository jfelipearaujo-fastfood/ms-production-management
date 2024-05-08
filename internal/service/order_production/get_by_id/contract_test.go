package get_by_id

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("Should return nil when valid", func(t *testing.T) {
		// Arrange
		input := GetOrderProductionByIdInput{
			OrderId: uuid.NewString(),
		}

		// Act
		err := input.Validate()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return error when invalid", func(t *testing.T) {
		// Arrange
		input := GetOrderProductionByIdInput{
			OrderId: "123",
		}

		// Act
		err := input.Validate()

		// Assert
		assert.Error(t, err)
	})
}
