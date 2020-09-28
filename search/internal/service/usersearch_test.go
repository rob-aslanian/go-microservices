package service

import (
	"testing"

	"gitlab.lan/Rightnao-site/microservices/search/internal/requests"
	"gitlab.lan/Rightnao-site/microservices/search/internal/search-errors"
)

func TestUserSearchValidator(t *testing.T) {
	tables := []struct {
		Data   *requests.UserSearch
		Result error
	}{
		{
			Data: &requests.UserSearch{
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
				MinAge: 18,
				MaxAge: 32,
			},
			Result: nil,
		},
		{
			Data: &requests.UserSearch{
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
					"ge", "us",
				},
				CityID: []string{
					"2221", "221",
				},
				ConnectionsOfID: []string{
					"2221", "221",
				},
				MinAge: 18,
				MaxAge: 160,
			},
			Result: searchErrors.InvalidAge,
		},
		{
			Data: &requests.UserSearch{
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
					"ge", "us",
				},
				CityID: []string{
					"2221", "221",
				},
				ConnectionsOfID: []string{
					"2221", "221",
				},
				MinAge: 160,
				MaxAge: 90,
			},
			Result: searchErrors.InvalidAge,
		},
		{
			Data: &requests.UserSearch{
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
					"ge", "us",
				},
				CityID: []string{
					"2221", "221",
				},
				ConnectionsOfID: []string{
					"2221", "221",
				},
				MinAge: 60,
				MaxAge: 20,
			},
			Result: searchErrors.InvalidAge,
		},
		{
			Data: &requests.UserSearch{
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
					"ge22222", "s",
				},
				CityID: []string{
					"2221", "221",
				},
				ConnectionsOfID: []string{
					"2221", "221",
				},
			},
			Result: searchErrors.InvalidCountry,
		},
		{
			Data: &requests.UserSearch{
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
					"ge", "us",
				},
				CityID: []string{
					"sdfsfds", "sfsd",
				},
				ConnectionsOfID: []string{
					"sssss", "221",
				},
			},
			Result: searchErrors.InvalidCity,
		},
		{
			Data: &requests.UserSearch{
				Firstname: []string{
					"Vakho", "vladimir",
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
					"englishenglishenglishenglishenglishenglishenglishenglishenglishenglishenglishenglishenglishenglishenglishenglishenglishenglishenglishenglishenglishenglishenglishenglishenglishenglish", "russian",
				},
				CountryID: []string{
					"ge", "us",
				},
				CityID: []string{
					"2221", "221",
				},
				ConnectionsOfID: []string{
					"2221", "221",
				},
			},
			Result: searchErrors.Max100,
		},
		{
			Data: &requests.UserSearch{
				Firstname: []string{
					"Vakho", "123123",
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
					"shs", "russian",
				},
				CountryID: []string{
					"ge", "us",
				},
				CityID: []string{
					"2221", "221",
				},
				ConnectionsOfID: []string{
					"2221", "221",
				},
			},
			Result: searchErrors.InvalidName,
		},
	}

	for _, ta := range tables {
		err := userSearchValidator(ta.Data)
		if err != ta.Result {
			t.Errorf("Error: userSearchValidator(%+v) got %v, expected: %v", ta.Data, err, ta.Result)
		}
	}
}
