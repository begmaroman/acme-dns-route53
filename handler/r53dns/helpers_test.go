package r53dns

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildQuotedValue(t *testing.T) {
	testTable := []*struct {
		testName      string
		actualValue   string
		expectedValue string
	}{
		{
			testName:      "simple string",
			actualValue:   "some string",
			expectedValue: `"some string"`,
		},
		{
			testName:      "string with quote",
			actualValue:   "some \" string",
			expectedValue: `"some " string"`,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.testName, func(t *testing.T) {
			actVal := buildQuotedValue(tt.actualValue)

			require.Equal(t, tt.expectedValue, actVal)
		})
	}
}
