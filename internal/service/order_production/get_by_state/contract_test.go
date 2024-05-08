package get_by_state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("Should return nil when valid", func(t *testing.T) {
		// Arrange
		input := GetOrderProductionByStateInput{
			State: "Received",
		}

		// Act
		err := input.Validate()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return error when invalid", func(t *testing.T) {
		// Arrange
		input := GetOrderProductionByStateInput{
			State: "",
		}

		// Act
		err := input.Validate()

		// Assert
		assert.Error(t, err)
	})

	t.Run("Should return error when invalid", func(t *testing.T) {
		// Arrange
		input := GetOrderProductionByStateInput{
			State: "Invalid",
		}

		// Act
		err := input.Validate()

		// Assert
		assert.Error(t, err)
	})
}
