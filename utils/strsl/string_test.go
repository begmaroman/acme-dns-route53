package strsl

import "testing"

func TestContainsStringSlice(t *testing.T) {
	testTable := []*struct {
		testName  string
		baseSlice []string
		subSlice  []string
		expected  bool
	}{
		{
			testName:  "same slices",
			baseSlice: []string{"1", "2", "test 123"},
			subSlice:  []string{"1", "2", "test 123"},
			expected:  true,
		},
		{
			testName:  "the base slice contains the subslice",
			baseSlice: []string{"1", "2", "3", "4"},
			subSlice:  []string{"1", "2", "3"},
			expected:  true,
		},
		{
			testName:  "the base slice is shorter than the sub-slice",
			baseSlice: []string{"1", "2"},
			subSlice:  []string{"1", "2", "test 123"},
			expected:  false,
		},
		{
			testName:  "the base slice contains part of the subslice",
			baseSlice: []string{"1", "2", "5", "50"},
			subSlice:  []string{"1", "2", "test 123"},
			expected:  false,
		},
		{
			testName:  "the base slice doesn't contain the subslice",
			baseSlice: []string{"1", "2", "5", "50"},
			subSlice:  []string{"q", "w", "e"},
			expected:  false,
		},
		{
			testName:  "empty slices with same length",
			baseSlice: make([]string, 10),
			subSlice:  make([]string, 10),
			expected:  true,
		},
		{
			testName:  "empty slices with different length",
			baseSlice: make([]string, 10),
			subSlice:  make([]string, 5),
			expected:  true,
		},
		{
			testName:  "empty slices with base slice being shorter than the sub-slice",
			baseSlice: make([]string, 10),
			subSlice:  make([]string, 50),
			expected:  false,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.testName, func(t *testing.T) {
			if actual := ContainsSub(tt.baseSlice, tt.subSlice); actual != tt.expected {
				t.Errorf("ContainsStringSlice(%#v, %#v) = %t, want %t", tt.baseSlice, tt.subSlice, actual, tt.expected)
			}
		})
	}
}
