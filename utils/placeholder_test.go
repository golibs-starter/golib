package utils

import (
	assert "github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_ReplacePlaceholder_WhenValidPlaceholderAndEnvIsPresent_ShouldReturnCorrect(t *testing.T) {
	_ = os.Setenv("ENV_EXAMPLE", "test")
	defer func() {
		_ = os.Unsetenv("ENV_EXAMPLE")
	}()
	val, err := ReplacePlaceholder("${ENV_EXAMPLE}")
	assert.Nil(t, err)
	assert.Equal(t, "test", val)
}

func Test_ReplacePlaceholder_WhenValueIsNotString_ShouldReturnCurrentValue(t *testing.T) {
	val, err := ReplacePlaceholder(10)
	assert.Nil(t, err)
	assert.Equal(t, 10, val)
}

func Test_ReplacePlaceholder_WhenItIsNotPlaceholder_ShouldReturnCurrentValue(t *testing.T) {
	val1, err := ReplacePlaceholder("TEST}")
	assert.Nil(t, err)
	assert.Equal(t, "TEST}", val1)

	val2, err := ReplacePlaceholder("${TEST")
	assert.Nil(t, err)
	assert.Equal(t, "${TEST", val2)

	val3, err := ReplacePlaceholder(" ${TEST}") //starts with space
	assert.Nil(t, err)
	assert.Equal(t, " ${TEST}", val3)
}

func Test_ReplacePlaceholder_WhenEmptyPlaceholderKey_ShouldReturnError(t *testing.T) {
	val, err := ReplacePlaceholder("${}")
	assert.NotNil(t, err)
	assert.Nil(t, val)
}

func Test_ReplacePlaceholder_WhenValidPlaceholderAndEnvNotPresent_ShouldReturnError(t *testing.T) {
	val, err := ReplacePlaceholder("${ENV_EXAMPLE}")
	assert.NotNil(t, err)
	assert.Nil(t, val)
}

func Test_ReplacePlaceholder_WhenValidPlaceholderAndEnvIsPresentAndEmpty_ShouldReturnEmptyValue(t *testing.T) {
	_ = os.Setenv("ENV_EXAMPLE", "")
	defer func() {
		_ = os.Unsetenv("ENV_EXAMPLE")
	}()
	val, err := ReplacePlaceholder("${ENV_EXAMPLE}")
	assert.Nil(t, err)
	assert.Equal(t, "", val)
}
