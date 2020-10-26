package fs

import (
	"errors"
	"os"
	"testing"

	"github.com/aos-dev/go-storage/v2/services"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	c, err := NewStorager()
	assert.NotNil(t, c)
	assert.NoError(t, err)
}

func TestFormatOsError(t *testing.T) {
	testErr := errors.New("test error")
	tests := []struct {
		name     string
		input    error
		expected error
	}{
		{
			"not found",
			os.ErrNotExist,
			services.ErrObjectNotExist,
		},
		{
			"not supported error",
			testErr,
			testErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := formatError(tt.input)
			assert.True(t, errors.Is(err, tt.expected))
		})
	}
}
