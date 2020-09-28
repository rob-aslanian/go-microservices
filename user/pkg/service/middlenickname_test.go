package service

import (
	"testing"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

func TestMiddlenicknameValidator(t *testing.T) {

	Invalid := "sss@@!"
	LongName := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	Valid := "Valid"
	// val := "123"
	val2 := "vaxo "
	val3 := "vaxo_"
	val4 := "vaxo-vaxo"
	val5 := "vaxo da da"
	emptyString := ""

	tables := []struct {
		Name   *string
		Result error
	}{
		{
			Name:   &Invalid,
			Result: usersErrors.InvalidPatronycName,
		},
		{
			Name:   &LongName,
			Result: usersErrors.InvalidPatronycName,
		},
		{
			Name:   &Valid,
			Result: nil,
		},
		{
			Name:   &emptyString,
			Result: nil,
		},
		{
			Name:   &val2,
			Result: nil,
		},
		{
			Name:   &val3,
			Result: usersErrors.InvalidPatronycName,
		},
		{
			Name:   &val4,
			Result: nil,
		},
		{
			Name:   &val5,
			Result: nil,
		},
	}

	for _, ta := range tables {
		err := middlenicknameValidator(ta.Name)
		if err != ta.Result {
			t.Errorf("Error: middlenicknameValidator(%q) got %v, expected: %v", *ta.Name, err, ta.Result)
		}
	}
}
