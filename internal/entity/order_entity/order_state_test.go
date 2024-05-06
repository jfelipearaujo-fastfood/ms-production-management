package order_entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrderState(t *testing.T) {
	t.Run("Should return the correct state", func(t *testing.T) {
		// Arrange
		cases := []struct {
			title    string
			expected OrderState
		}{
			{"Received", Received},
			{"Processing", Processing},
			{"Completed", Completed},
			{"Delivered", Delivered},
			{"Cancelled", Cancelled},
		}

		for _, c := range cases {
			// Act
			res := NewOrderState(c.title)

			// Assert
			assert.Equal(t, c.expected, res)
		}
	})

	t.Run("Should return None when state is invalid", func(t *testing.T) {
		// Arrange
		title := "Invalid"

		// Act
		res := NewOrderState(title)

		// Assert
		assert.Equal(t, None, res)
	})
}

func TestCanTransitionTo(t *testing.T) {
	t.Run("Should return true when transition is allowed", func(t *testing.T) {
		// Arrange
		cases := []struct {
			from OrderState
			to   OrderState
		}{
			{None, Received},
			{Received, Processing},
			{Received, Cancelled},
			{Processing, Completed},
			{Processing, Cancelled},
			{Completed, Delivered},
		}

		// Act
		for _, c := range cases {
			res := c.from.CanTransitionTo(c.to)

			// Assert
			assert.True(t, res)
		}
	})

	t.Run("Should return false when transition is not allowed", func(t *testing.T) {
		// Arrange
		cases := []struct {
			from OrderState
			to   OrderState
		}{
			{None, Processing},
			{Received, Completed},
			{Processing, Received},
			{Completed, Received},
		}

		// Act
		for _, c := range cases {
			res := c.from.CanTransitionTo(c.to)

			// Assert
			assert.False(t, res)
		}
	})
}

func TestIsValidState(t *testing.T) {
	t.Run("Should return true when state is valid", func(t *testing.T) {
		// Arrange
		state := Received

		// Act
		res := IsValidState(state)

		// Assert
		assert.True(t, res)
	})

	t.Run("Should return false when state is invalid", func(t *testing.T) {
		// Arrange
		state := OrderState(0)

		// Act
		res := IsValidState(state)

		// Assert
		assert.False(t, res)
	})
}

func TestString(t *testing.T) {
	t.Run("Should return the string representation of the state", func(t *testing.T) {
		// Arrange
		cases := []struct {
			state    OrderState
			expected string
		}{
			{None, "None"},
			{Received, "Received"},
			{Processing, "Processing"},
			{Completed, "Completed"},
			{Delivered, "Delivered"},
			{Cancelled, "Cancelled"},
		}

		for _, c := range cases {
			// Act
			res := c.state.String()

			// Assert
			assert.Equal(t, c.expected, res)
		}
	})

	t.Run("Should return 'Unknown' when state is invalid", func(t *testing.T) {
		// Arrange
		state := OrderState(99)

		// Act
		res := state.String()

		// Assert
		assert.Equal(t, "Unknown", res)
	})
}
