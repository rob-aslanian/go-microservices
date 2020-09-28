package service

import (
	"testing"
)

func TestCapitalToLower(t *testing.T) {

	CAPITAL := "VAXO"
	small := "vaxo"

	tables := []struct {
		Test *string
	}{
		{
			Test: &CAPITAL,
		},
	}

	for _, ta := range tables {
		capitalToLower(ta.Test)
		if !(CAPITAL == small) {
			t.Errorf("Error: got %q, expected: %q", (*ta.Test), small)
		}
	}
}
