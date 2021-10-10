package config

import (
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestLoader_WhenFormatIsNotSupported_ShouldReturnError(t *testing.T) {
	_, err := NewLoader(Option{
		ActiveProfiles: []string{"file_in_yaml_format"},
		ConfigPaths:    []string{"./test_assets"},
		ConfigFormat:   "ini",
	}, []Properties{new(testStore)})
	assert.ErrorIs(t, err, ErrFormatNotSupported)
}
