package service

import (
	"testing"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

func TestUsernameValidator(t *testing.T) {

	Invalid := "sss@@!"
	LongName := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	Valid := "Validss"

	tables := []struct {
		Name   string
		Result error
	}{
		{
			Name:   "datosssao",
			Result: nil,
		},
		{
			Name:   Invalid,
			Result: usersErrors.InValidUserName,
		},
		{
			Name:   LongName,
			Result: usersErrors.InValidUserName,
		},
		{
			Name:   Valid,
			Result: nil,
		},
		{
			Name:   "",
			Result: usersErrors.InValidUserName,
		},
		{
			Name:   "vaxo_123",
			Result: nil,
		},
		{
			Name:   "vaxo.1",
			Result: nil,
		},
		{
			Name:   "vaXo_.1",
			Result: nil,
		},
	}

	for _, ta := range tables {
		err := userNameValidator(ta.Name)
		if err != ta.Result {
			t.Errorf("Error: got %v, expected: %v", err, nil)
		}
	}
}
