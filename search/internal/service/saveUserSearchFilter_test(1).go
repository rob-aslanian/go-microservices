package service

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"gitlab.lan/Rightnao-site/microservices/search/internal/requests"
)

func TestSaveUserSearchFilterValidator(t *testing.T) {
	emptystr := ""
	str := "123sdfsf"
	// longStr := "sfdsdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
	objID, _ := primitive.ObjectIDFromHex(str)
	emptyObj, _ := primitive.ObjectIDFromHex(emptystr)
	// longObjID, _ := primitive.ObjectIDFromHex(longStr)

	tables := []struct {
		Data   *requests.UserSearchFilter
		Result error
	}{
		{
			Data: &requests.UserSearchFilter{
				CompanyID: &emptyObj,
				Name:      "usersearch",
				UserID:    &objID,
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
			},
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
