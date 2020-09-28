package service

// import (
// 	"testing"

// 	"gitlab.lan/Rightnao-site/microservices/user/pkg/account"
// 	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
// )

// func TestAddAddressValidator(t *testing.T) {

// 	// t1, err := time.Parse(time.UnixDate, "Mon Jan 2 15:04:05 MST 2006")
// 	// if err != nil {
// 	// 	t.Fatal("wrong date")
// 	// }
// 	// t2, err := time.Parse(time.UnixDate, "Mon Jan 2 15:04:05 MST 2006")
// 	// descr := "2222"

// 	tables := []struct {
// 		Address *account.MyAddress
// 		Result  error
// 	}{
// 		{
// 			Address: &account.MyAddress{
// 				Name:      "vaxo",
// 				Apartment: "dfdfdfdfdfdfdfdfdfdfdfdfddfdfdfdfdfdffdfdfdfdsdfsdfddfdfdfdfdfdfdfdfddddddddddddsddddddddddddddfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdsdfsdfddfdfdfdfdfdfdfdfddddddddddddsddddddddddddddfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdsdfsdfddfdfdfdfdfdfdfdfddddddddddddsddddddddddddd",
// 				ZIP:       "dfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdsdfsdfddfdfdfdfdfdfdfdfddddddddddddsddddddddddddd",
// 				Street:    "dfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdfdsdfsdfddfdfdfdfdfdfdfdfddddddddddddsddddddddddddd",
// 				Location: Location.Country{
// 					"GE",
// 				},
// 			},
// 			Result: usersErrors.Max128,
// 		},
// 		{
// 			Address: &account.MyAddress{
// 				Name:      "",
// 				Apartment: "",
// 				ZIP:       "",
// 				Street:    "",
// 			},
// 			Result: nil,
// 		},
// 		{
// 			Address: &account.MyAddress{
// 				Apartment: "dfdf",
// 				ZIP:       "12",
// 				Street:    "dfd",
// 			},
// 			Result: nil,
// 		},
// 	}

// 	for _, ta := range tables {
// 		err := addressValidator(ta.Address)
// 		if err != ta.Result {
// 			t.Errorf("TestAddMyAddressValidator(%v) Error: got %v, expected: %v", ta.Address, err, nil)
// 		}
// 	}
// }
