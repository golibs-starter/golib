package utils

import (
	"testing"
	"time"
)

func TestTime2Ms(t *testing.T) {
	testCases := []struct {
		time     time.Time
		expected int64
	}{
		{
			time:     time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			time:     time.Date(2019, 7, 5, 0, 0, 0, 0, time.UTC),
			expected: 1562284800000,
		},
		{
			time:     time.Date(2019, 7, 5, 0, 0, 0, 100, time.UTC),
			expected: 1562284800000,
		},
	}

	for _, testCase := range testCases {
		actual := Time2Ms(testCase.time)
		if actual != testCase.expected {
			t.Errorf("Expected: %v, but got: %v", testCase.expected, actual)
		}
	}
}
