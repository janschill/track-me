package icloud

import "testing"

func TestBase62ToInt(t *testing.T) {
	tests := []struct {
		input    byte
		expected int
	}{
		{'0', 0},
		{'9', 9},
		{'A', 10},
		{'Z', 35},
		{'a', 36},
		{'z', 61},
		{'!', -1}, // Invalid character
	}

	for _, test := range tests {
		result := base62ToInt(test.input)
		if result != test.expected {
			t.Errorf("base62ToInt(%c) = %d; want %d", test.input, result, test.expected)
		}
	}
}

func TestGetPartitionFromToken(t *testing.T) {
	tests := []struct {
		token    string
		expected string
	}{
		{"A1AG4Tcsm1OG3pH", "01"},
		{"ArAG4Tcsm1OG3pH", "53"},
		{"B1A55Z2WMR2vuY", "72"},
		{"B1A5nhQST2r3K4H", "72"},
		{"B1AG4Tcsm1OG3pH", "72"},
		{"B1AG6XBub2QnCol", "72"},
		{"B1AGY8gBY1I1nym", "72"},
	}

	for _, test := range tests {
		result := getPartitionFromToken(test.token)
		if result != test.expected {
			t.Errorf("getPartitionFromToken(%s) = %s; want %s", test.token, result, test.expected)
		}
	}
}
