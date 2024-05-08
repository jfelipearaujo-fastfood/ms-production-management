package create

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("Should return nil when valid", func(t *testing.T) {
		// Arrange
		input := CreateOrderProductionInput{
			OrderId: uuid.NewString(),
			Items: []CreateOrderProductionItemInput{
				{
					Id:       uuid.NewString(),
					Name:     "Test",
					Quantity: 1,
				},
			},
		}

		// Act
		err := input.Validate()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return error when invalid", func(t *testing.T) {
		// Arrange
		input := CreateOrderProductionInput{
			OrderId: "123",
			Items: []CreateOrderProductionItemInput{
				{
					Id:       "123",
					Name:     "",
					Quantity: 1,
				},
			},
		}

		// Act
		err := input.Validate()

		// Assert
		assert.Error(t, err)
	})
}
