package order_entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewItem(t *testing.T) {
	t.Run("Should create a new item", func(t *testing.T) {
		// Arrange
		expect := Item{
			Id:       "1",
			Name:     "name",
			Quantity: 2,
		}

		// Act
		res := NewItem("1", "name", 2)

		// Assert
		assert.Equal(t, expect, res)
	})
}
