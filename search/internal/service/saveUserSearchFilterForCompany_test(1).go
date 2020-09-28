package service

import (
	"testing"

	"gitlab.lan/Rightnao-site/microservices/search/internal/requests"
)

func TestSaveUserSearchFilterForCompanyValidator(t *testing.T) {
	// emptystr := ""
	// str := "123sdfsf"
	// longStr := "sfdsdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
	// objID, _ := primitive.ObjectIDFromHex(str)
	// emptyObj, _ := primitive.ObjectIDFromHex(emptystr)
	// longObjID, _ := primitive.ObjectIDFromHex(longStr)

	req := requests.UserSearchFilter{
		// CompanyID:
		Name: "usersearch",
		// UserID: emptyObj,
		UserSearch: requests.UserSearch{
			Firstname: []string{
				"Vakho", "Vladimir",
			},
			Lastname: []string{
				"Nodadze", "Kopaliani",
			},
			CurrentCompany: []string{
				"Hypercube", "sdfsdfsdfds",
			},
			PastCompany: []string{
				"Respo", "sdfsfsfw",
			},
			School: []string{
				"24esajaro", "splidfgd",
			},
			Skill: []string{
				"developer", "yeaaaaa",
			},
			Language: []string{
				"english", "russian",
			},
			Nickname: []string{
				"vovchick", "russs",
			},
			FieldOfStudy: []string{
				"computer science", "russian",
			},
			Industry: []string{
				"english", "russian",
			},
			Interest: []string{
				"english", "russian",
			},
			Position: []string{
				"english", "russian",
			},
			Degree: []string{
				"english", "russian",
			},
			CountryID: []string{
				"GE", "us",
			},
			CityID: []string{
				"2221", "221",
			},
			ConnectionsOfID: []string{
				"2221", "221",
			},
		},
	}

	req.SetCompanyID("5cf64b59840a0b95101b9ce6")
	req.SetUserID("5cf64b59840a0b95101b9xd7")

	tables := []struct {
		Data   *requests.UserSearchFilter
		Result error
	}{
		{
			Data:   &req,
			Result: nil,
		},
	}

	for _, ta := range tables {
		err := saveUserSearchFilterValidatorForCompany(ta.Data)
		if err != ta.Result {
			t.Errorf("Error: userSearchValidator(%+v) got %v, expected: %v", ta.Data, err, ta.Result)
		}
	}
}
