package golib

import (
	assert "github.com/stretchr/testify/require"
	"testing"
)

type testDummyProps struct {
}

func (t testDummyProps) Prefix() string {
	return "prefix.test"
}

func Test_makeSampleProperties_WhenConstructorReturnPointer_ShouldReturnSuccess(t *testing.T) {
	f := func() (*testDummyProps, error) {
		return nil, nil
	}
	out, err := makeSampleProperties(f)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, "prefix.test", out.Prefix())
}

func Test_makeSampleProperties_WhenConstructorReturnStruct_ShouldReturnSuccess(t *testing.T) {
	f := func() (testDummyProps, error) {
		return testDummyProps{}, nil
	}
	out, err := makeSampleProperties(f)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, "prefix.test", out.Prefix())
}

func Test_makeSampleProperties_WhenConstructorReturnPointerAtLast_ShouldReturnSuccess(t *testing.T) {
	f := func() (error, *testDummyProps) {
		return nil, nil
	}
	out, err := makeSampleProperties(f)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, "prefix.test", out.Prefix())
}

func Test_makeSampleProperties_WhenNoPropertiesFound_ShouldReturnError(t *testing.T) {
	f := func() (int, error) {
		return 0, nil
	}
	out, err := makeSampleProperties(f)
	assert.Error(t, err)
	assert.Nil(t, out)
}

func Test_makeSampleProperties_WhenNotAFunction_ShouldReturnError(t *testing.T) {
	out, err := makeSampleProperties(0)
	assert.Error(t, err)
	assert.Nil(t, out)
}
