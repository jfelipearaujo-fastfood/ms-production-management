package dummy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetName(t *testing.T) {
	t.Run("Should return the name", func(t *testing.T) {
		// Arrange
		d := NewDummy("test")

		// Act
		name := d.GetName()

		// Assert
		assert.Equal(t, "test", name)
	})
}
