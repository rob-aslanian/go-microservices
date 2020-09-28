package service

import (
	"testing"

	"gitlab.lan/Rightnao-site/microservices/search/internal/search-errors"

	"gitlab.lan/Rightnao-site/microservices/search/internal/requests"
)

func TestCompanySearch(t *testing.T) {
	tables := []struct {
		Data   *requests.CompanySearch
		Result error
	}{
		{
			Data: &requests.CompanySearch{
				Country: []string{
					"GE",
					"US",
				},
				City: []requests.City{
					{ID: "asdsfdsfs"},
					{ID: "sdfssdfsf"},
				},
			},
			Result: nil,
		},
		{
			Data: &requests.CompanySearch{
				Country: []string{
					"Geeeee",
					"us2212312",
				},
				City: []requests.City{
					{ID: "asdsfdsfs"},
					{ID: "sdfssdfsf"},
				},
			},
			Result: searchErrors.InvalidCountry,
		},
		{
			Data: &requests.CompanySearch{
				Country: []string{
					"US",
					"GE",
				},
				City: []requests.City{
					{ID: "asdsfdsfs"},
					{ID: "sdfssdfsf"},
				},
			},
			Result: nil,
		},
	}

	for _, ta := range tables {
		err := companySearchValidator(ta.Data)
		if err != ta.Result {
			t.Errorf("Error: companySearchValidator(%+v) got %v, expected: %v", ta.Data, err, ta.Result)
		}
	}
}
