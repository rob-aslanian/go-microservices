package service

//
// import (
// 	"testing"
//
// 	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/account"
// 	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/company-errors"
// )
//
// func TestRegisterValidator(t *testing.T) {
//
// 	tables := []struct {
// 		Address *account.Account
// 		Result  error
// 	}{
// 		{
// 			Address: &account.Account{
// 				Addresses: []*account.Address{
// 					{
// 						Apartment: "564",
// 						ZIPCode:   "212",
// 						Street:    "ssd",
// 					},
// 				},
// 				Name: "ss",
// 				URL:  "wwwgoogle.com",
// 				Industry: account.Industry{
// 					Main: "22",
// 				},
// 				Phones: []*account.Phone{
// 					{
// 						Number: "598661708",
// 						CountryCode: account.CountryCode{
// 							Code:      "+995",
// 							CountryID: "GE",
// 						},
// 					},
// 				},
// 			},
// 			Result: nil,
// 		},
// 		{
// 			Address: &account.Account{
// 				Addresses: []*account.Address{
// 					{
// 						Apartment: "",
// 						ZIPCode:   "212",
// 						Street:    "ssd",
// 					},
// 				},
// 				Name: "ss",
// 				URL:  "wwwgoogle1com",
// 				Industry: account.Industry{
// 					Main: "22",
// 				},
// 				Phones: []*account.Phone{
// 					{
// 						Number: "598661708",
// 						CountryCode: account.CountryCode{
// 							Code:      "+995",
// 							CountryID: "GE",
// 						},
// 					},
// 				},
// 			},
// 			Result: companyErrors.SpecificRequired,
// 		},
// 		{
// 			Address: &account.Account{
// 				Addresses: []*account.Address{
// 					{
// 						Apartment: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
// 						ZIPCode:   "212",
// 						Street:    "ssd",
// 					},
// 				},
// 				Name: "ss",
// 				URL:  "wwwgoogle2com",
// 				Industry: account.Industry{
// 					Main: "22",
// 				},
// 				Phones: []*account.Phone{
// 					{
// 						Number: "598661708",
// 						CountryCode: account.CountryCode{
// 							Code:      "+995",
// 							CountryID: "GE",
// 						},
// 					},
// 				},
// 			},
// 			Result: companyErrors.Max128,
// 		},
// 		{
// 			Address: &account.Account{
// 				Addresses: []*account.Address{
// 					{
// 						Apartment: "d",
// 						ZIPCode:   "212",
// 						Street:    "ssd",
// 					},
// 				},
// 				Name: "ss",
// 				URL:  "wwwgoogle3om",
// 				Industry: account.Industry{
// 					Main: "22",
// 				},
// 				Phones: []*account.Phone{
// 					{
// 						Number: "598661708",
// 						CountryCode: account.CountryCode{
// 							Code:      "+99",
// 							CountryID: "GE",
// 						},
// 					},
// 				},
// 			},
// 			Result: companyErrors.InValidPhone,
// 		},
// 		{
// 			Address: &account.Account{
// 				Addresses: []*account.Address{
// 					{
// 						Apartment: "d",
// 						ZIPCode:   "212",
// 						Street:    "ssd",
// 					},
// 				},
// 				Name: "ss",
// 				URL:  "wwwe",
// 				Industry: account.Industry{
// 					Main: "22",
// 				},
// 				Phones: []*account.Phone{
// 					{
// 						Number: "598661708",
// 						CountryCode: account.CountryCode{
// 							Code:      "+995",
// 							CountryID: "GE",
// 						},
// 					},
// 				},
// 			},
// 			Result: nil,
// 		},
// 		{
// 			Address: &account.Account{
// 				Addresses: []*account.Address{
// 					{
// 						Apartment: "d",
// 						ZIPCode:   "212",
// 						Street:    "ssd",
// 					},
// 				},
// 				Name: "ss",
// 				URL:  "сайтрф",
// 				Industry: account.Industry{
// 					Main: "22",
// 				},
// 				Phones: []*account.Phone{
// 					{
// 						Number: "598661708",
// 						CountryCode: account.CountryCode{
// 							Code:      "+995",
// 							CountryID: "GE",
// 						},
// 					},
// 				},
// 			},
// 			Result: companyErrors.InvalidURL,
// 		},
// 		{
// 			Address: &account.Account{
// 				Addresses: []*account.Address{
// 					{
// 						Apartment: "d",
// 						ZIPCode:   "212",
// 						Street:    "ssd",
// 					},
// 				},
// 				Name: "ss",
// 				URL:  "sdssdsds",
// 				Industry: account.Industry{
// 					Main: "1",
// 				},
// 				Phones: []*account.Phone{
// 					{
// 						Number: "5986617 08",
// 						CountryCode: account.CountryCode{
// 							Code:      "995",
// 							CountryID: "GE",
// 						},
// 					},
// 				},
// 			},
// 			Result: nil,
// 		},
// 		{
// 			Address: &account.Account{
// 				Addresses: []*account.Address{
// 					{
// 						Apartment: "d",
// 						ZIPCode:   "212",
// 						Street:    "ssd",
// 					},
// 				},
// 				Name: "ss",
// 				URL:  "wwwgoogle5om",
// 				Industry: account.Industry{
// 					Main: "22",
// 				},
// 				Phones: []*account.Phone{
// 					{
// 						Number: "5986617 08",
// 						CountryCode: account.CountryCode{
// 							Code:      "+995",
// 							CountryID: "GE",
// 						},
// 					},
// 				},
// 			},
// 			Result: nil,
// 		},
// 	}
//
// 	for _, ta := range tables {
// 		err := registerValidator(ta.Address)
// 		if err != ta.Result {
// 			t.Errorf("Error: registerValidator(%+v) got %v, expected: %v", ta.Address, err, nil)
// 		}
// 	}
// }
