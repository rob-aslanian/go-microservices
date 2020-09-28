package searchrepo

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/olivere/elastic"
	"gitlab.lan/Rightnao-site/microservices/search/internal/company"
	"gitlab.lan/Rightnao-site/microservices/search/internal/requests"
)

// UserSearch ...
func (r Repository) UserSearch(ctx context.Context, data *requests.UserSearch, connections, blockedIDs []string) ([]string, int64, error) {
	must := make([]elastic.Query, 0)
	mustNot := make([]elastic.Query, 0)
	should := make([]elastic.Query, 0)

	// search only users
	must = append(must, elastic.NewTermQuery("document_type", "user"))

	// search only activated users
	must = append(must, elastic.NewTermQuery("status", "ACTIVATED"))

	if len(blockedIDs) > 0 {
		mustNot = append(mustNot, elastic.NewIdsQuery("_doc").Ids(blockedIDs...))
	}

	// search by keyword
	if len(data.Keyword) > 0 {
		keywords := make([]elastic.Query, 0, len(data.Keyword))
		for i := range data.Keyword {
			keywords = append(keywords, elastic.NewMultiMatchQuery(
				data.Keyword[i],
				"first_name", "last_name", "patronymic.name", "middlename.name",
				"native_name.name", "nickname.name", "headline", "story", "emails.email",
				"skills.skill", "educations.school", "educations.degree", "experiences.position",
				"interests.interest", "known_languages.language", "experiences.company",
				"educations.field_study",
			),
			)
		}
		must = append(must, elastic.NewBoolQuery().Should(keywords...))
	}

	// search by gender
	if data.IsMale {
		must = append(must, elastic.NewTermQuery("gender.gender", "MALE"))
	} else if data.IsFemale {
		must = append(must, elastic.NewTermQuery("gender.gender", "FEMALE"))
	}

	// search for student
	if data.IsStudent {
		must = append(must,
			elastic.NewBoolQuery().Should(
				elastic.NewRangeQuery("educations.finish_date").Gt(time.Now()),
				elastic.NewMatchQuery("educations.is_currenlty_study", true),
			),
		)
	}

	// search among certain users
	if len(connections) > 0 {
		must = append(must, elastic.NewIdsQuery("_doc").Ids(connections...))
	}

	if data.MinAge > 0 && data.MaxAge > 0 {
		must = append(must, elastic.NewRangeQuery("birthday.birthday").
			Gte(time.Date(time.Now().Year()-int(data.MaxAge), time.January, 1, 0, 0, 0, 0, time.Now().Location())). // TODO: set month and day
			Lte(time.Date(time.Now().Year()-int(data.MinAge), time.January, 1, 0, 0, 0, 0, time.Now().Location())),
		)
	} else {
		if data.MinAge > 0 {
			must = append(must, elastic.NewRangeQuery("birthday.birthday").
				Lte(time.Date(time.Now().Year()-int(data.MinAge), time.January, 1, 0, 0, 0, 0, time.Now().Location())),
			)
		} else if data.MaxAge > 0 {
			must = append(must, elastic.NewRangeQuery("birthday.birthday").
				Gte(time.Date(time.Now().Year()-int(data.MaxAge), time.January, 1, 0, 0, 0, 0, time.Now().Location())),
			)
		}
	}

	// search by skills
	if len(data.Skill) > 0 {
		skills := make([]elastic.Query, 0, len(data.Skill))
		for i := range data.Skill {
			skills = append(skills, elastic.NewMatchPhraseQuery("skills.skill", data.Skill[i]))
			skills = append(skills, elastic.NewQueryStringQuery(data.Skill[i]).Field("skills.translations.*.skill"))
		}
		must = append(must, elastic.NewBoolQuery().Should(skills...))
	}

	// search by school
	if len(data.School) > 0 {
		schools := make([]elastic.Query, 0, len(data.School))
		for i := range data.School {
			schools = append(schools, elastic.NewMatchPhraseQuery("educations.school", data.School[i]))
			schools = append(schools, elastic.NewQueryStringQuery(data.School[i]).Field("educations.translations.*.school"))
		}
		must = append(must, elastic.NewBoolQuery().Should(schools...))
	}

	// search by degree
	if len(data.Degree) > 0 {
		degrees := make([]elastic.Query, 0, len(data.Degree))
		for i := range data.Degree {
			degrees = append(degrees, elastic.NewMatchPhraseQuery("educations.degree", data.Degree[i]))
			degrees = append(degrees, elastic.NewQueryStringQuery(data.Degree[i]).Field("educations.translations.*.degree"))
		}
		must = append(must, elastic.NewBoolQuery().Should(degrees...))
	}

	// search by job position
	if len(data.Position) > 0 {
		position := make([]elastic.Query, 0, len(data.Position))
		for i := range data.Position {
			position = append(position,
				elastic.NewNestedQuery("experiences",
					elastic.NewMatchPhraseQuery("experiences.position", data.Position[i]),
				),
			)
			position = append(position, elastic.NewQueryStringQuery(data.Position[i]).Field("experiences.translations.*.position"))
		}
		must = append(must, elastic.NewBoolQuery().Should(position...))
	}

	// search by interest
	if len(data.Interest) > 0 {
		interests := make([]elastic.Query, 0, len(data.Interest))
		for i := range data.Interest {
			interests = append(interests, elastic.NewMatchPhraseQuery("interests.interest", data.Interest[i]))
			interests = append(interests, elastic.NewQueryStringQuery(data.Interest[i]).Field("interests.translations.*.interest"))
		}
		must = append(must, elastic.NewBoolQuery().Should(interests...))
	}

	// search by nickname
	if len(data.Nickname) > 0 {
		nicknames := make([]elastic.Query, 0, len(data.Nickname))
		for i := range data.Nickname {
			nicknames = append(nicknames, elastic.NewMatchQuery("nickname.name", data.Nickname[i]))
			nicknames = append(nicknames, elastic.NewQueryStringQuery(data.Nickname[i]).Field("translation.*.nickname"))
		}
		must = append(must, elastic.NewBoolQuery().Should(nicknames...))
	}

	// search by lastname
	if len(data.Lastname) > 0 {
		lastnames := make([]elastic.Query, 0, len(data.Lastname))
		for i := range data.Lastname {
			lastnames = append(lastnames, elastic.NewMatchQuery("last_name", data.Lastname[i]))
			lastnames = append(lastnames, elastic.NewQueryStringQuery(data.Lastname[i]).Field("translation.*.last_name"))
		}
		must = append(must, elastic.NewBoolQuery().Should(lastnames...))
	}

	// search by language
	if len(data.Language) > 0 {
		langs := make([]elastic.Query, 0, len(data.Language))
		for i := range data.Language {
			langs = append(langs, elastic.NewTermQuery("known_languages.language", data.Language[i]))
		}
		must = append(must, elastic.NewBoolQuery().Must(langs...))
	}

	// search by first name
	if len(data.Firstname) > 0 {
		names := make([]elastic.Query, 0, len(data.Firstname))
		for i := range data.Firstname {
			names = append(names, elastic.NewFuzzyQuery("first_name", data.Firstname[i]))
			names = append(names, elastic.NewQueryStringQuery(data.Firstname[i]).Field("translation.*.first_name"))
		}
		must = append(must, elastic.NewBoolQuery().Should(names...))
	}

	// search by past company
	if len(data.PastCompany) > 0 {
		pastCompanies := make([]elastic.Query, 0, len(data.PastCompany))
		for i := range data.PastCompany {
			pastCompanies = append(pastCompanies,
				elastic.NewBoolQuery().Must(
					elastic.NewBoolQuery().Should(
						elastic.NewMatchQuery("experiences.company", data.PastCompany[i]),
						elastic.NewQueryStringQuery(data.PastCompany[i]).Field("experiences.translations.*.company"),
					),
					elastic.NewMatchQuery("experiences.currently_work", false),
				),
			)
		}
		must = append(must,
			elastic.NewNestedQuery("experiences",
				elastic.NewBoolQuery().Must(pastCompanies...),
			),
		)
	}

	// search by current company
	if len(data.CurrentCompany) > 0 {
		currentCompanies := make([]elastic.Query, 0, len(data.CurrentCompany))
		for i := range data.CurrentCompany {
			currentCompanies = append(currentCompanies,
				elastic.NewBoolQuery().Must(
					elastic.NewBoolQuery().Should(
						elastic.NewMatchQuery("experiences.company", data.CurrentCompany[i]),
						elastic.NewQueryStringQuery(data.CurrentCompany[i]).Field("experiences.translations.*.company"),
					),
					elastic.NewMatchQuery("experiences.currently_work", true),
				),
			)
		}
		must = append(must,
			elastic.NewNestedQuery("experiences",
				elastic.NewBoolQuery().Must(currentCompanies...),
			),
		)
	}

	// search by field of study
	if len(data.FieldOfStudy) > 0 {
		fields := make([]elastic.Query, 0, len(data.FieldOfStudy))
		for i := range data.FieldOfStudy {
			fields = append(fields, elastic.NewMatchQuery("educations.field_study", data.FieldOfStudy[i]))
			fields = append(fields, elastic.NewQueryStringQuery(data.FieldOfStudy[i]).Field("educations.translations.*.field_study"))
		}
		must = append(must, elastic.NewBoolQuery().Should(fields...))
	}

	// search by country
	if len(data.CountryID) > 0 {
		countries := make([]elastic.Query, 0, len(data.CountryID))
		for i := range data.CountryID {
			countries = append(countries,
				// elastic.NewBoolQuery().Should(
				// elastic.NewTermQuery("location.country.id", data.CountryID[i]),
				elastic.NewBoolQuery().Must(
					elastic.NewMatchQuery("my_addresses.location.country.id", data.CountryID[i]),
					elastic.NewMatchQuery("my_addresses.primary", true),
				),
				// ),
			)
		}
		must = append(must, elastic.NewBoolQuery().Should(countries...))
	}

	// search by city
	if len(data.CityID) > 0 {
		cities := make([]elastic.Query, 0, len(data.CityID))
		for i := range data.CityID {
			cities = append(cities,
				elastic.NewBoolQuery().Must(
					elastic.NewMatchQuery("my_addresses.location.city.id", data.CityID[i]),
					elastic.NewMatchQuery("my_addresses.primary", true),
				// elastic.NewBoolQuery().Should(
				// 	elastic.NewTermQuery("location.city.id", data.CityID[i]),
				// 	elastic.NewTermQuery("my_addresses.location.city.id", data.CityID[i]),
				),
			)
		}
		must = append(must, elastic.NewBoolQuery().Should(cities...))
	}

	// search by past companies
	if len(data.Industry) > 0 {
		industries := make([]elastic.Query, 0, len(data.CityID))
		for i := range data.Industry {
			industries = append(industries,
				elastic.NewHasChildQuery("candidate", elastic.NewMatchPhraseQuery("career_interests.industry", data.Industry[i])),
			)
		}
		must = append(must, elastic.NewBoolQuery().Should(industries...))
	}

	if data.FullName != "" {
		names := strings.Split(data.FullName, " ")
		fullNameQuery := make([]elastic.Query, 0)

		log.Println("names:", names)

		for i := range names {
			fullNameQuery = append(fullNameQuery, elastic.NewQueryStringQuery(names[i]+`*`).Field("first_name").Field("last_name"))
			fullNameQuery = append(fullNameQuery, elastic.NewQueryStringQuery(names[i]+`*`).Field("translation.*"))
		}

		must = append(must, elastic.NewBoolQuery().Should(fullNameQuery...))
	}

	// -------------------

	search := elastic.NewBoolQuery()
	if len(must) > 0 {
		search.Must(must...)
	}
	if len(should) > 0 {
		search.Should(should...)
	}
	if len(mustNot) > 0 {
		search.MustNot(mustNot...)
	}

	searchQuery := r.client.
		Search("users_db.users").
		Query(
			search,
		).
		SortBy(
			elastic.NewFieldSort("_score").Desc(),
		).
		FetchSource(false)

	searchResultAmount, err := searchQuery.Do(ctx)
	if err != nil {
		log.Println("error: sending search:", err)
		return []string{}, 0, err
	}

	// if data.After != "" {
	// 	searchQuery.SearchAfter(data.After)
	// }
	if data.After != "" {
		// searchQuery.SearchAfter(data.GetAfter())
		num, _ := strconv.Atoi(data.After)
		searchQuery.From(num)
	}

	if data.First != 0 {
		searchQuery.Size(int(data.First))
	}

	searchResult, err := searchQuery.Do(ctx)
	if err != nil {
		log.Println("error: sending search:", err)
		return []string{}, 0, err
	}

	// -----------

	ids := make([]string, 0, len(searchResult.Hits.Hits))
	for i := range searchResult.Hits.Hits {
		ids = append(ids, searchResult.Hits.Hits[i].Id)
	}

	// // debug
	// h := make([]interface{}, 0)
	// for _, res := range searchResult.Hits.Hits {
	// 	h = append(h, *res)
	// }
	// log.Printf("Search result: %+v\n", h)
	// // ----

	return ids, searchResultAmount.TotalHits(), nil
}

// JobSearch ...
func (r Repository) JobSearch(ctx context.Context, data *requests.JobSearch, ids []string, blockedIDs []string) ([]string, int64, error) {
	must := make([]elastic.Query, 0)
	should := make([]elastic.Query, 0)
	mustNot := make([]elastic.Query, 0)

	// search only among active jobs
	// should be NewTermQuery
	must = append(must, elastic.NewMatchQuery("status", "Active"))
	must = append(must, elastic.NewTermQuery("document_type", "job"))

	// search by keyword
	if len(data.Keyword) > 0 {
		keywords := make([]elastic.Query, 0, len(data.Keyword))
		for i := range data.Keyword {
			keywords = append(keywords, elastic.NewMultiMatchQuery(
				data.Keyword[i],
				"job_details.title",
				"job_details.jobfunctions",
				"job_details.descriptions.description",
				"job_details.descriptions.whyus",
				"job_details.required_skills",
				"job_details.required_languages",
				"job_details.required_educations",
				"job_details.requiredeligibility",
				"job_details.requiredlicenses",
				"job_details.benefits",
			),
			)
		}
		must = append(must, elastic.NewBoolQuery().Should(keywords...))
	}

	// search by ids
	if len(ids) > 0 {
		must = append(must, elastic.NewHasParentQuery(
			"company",
			elastic.NewIdsQuery("_doc").Ids(ids...),
		),
		)
	}

	// search among followings
	if data.IsFollowing && len(ids) > 0 {
		must = append(must, elastic.NewHasParentQuery(
			"company",
			elastic.NewIdsQuery("_doc").Ids(ids...),
		),
		)
	}

	if len(blockedIDs) > 0 {
		mustNot = append(mustNot, elastic.NewHasParentQuery(
			"company",
			elastic.NewIdsQuery("_doc").Ids(blockedIDs...),
		),
		)
	}

	// with salary
	if data.WithSalary {
		mustNot = append(
			mustNot,
			elastic.NewBoolQuery().Must(
				elastic.NewTermQuery("job_details.salarymax", 0),
				elastic.NewTermQuery("job_details.salarymin", 0),
			),
		)
	}

	// search by currency
	if data.Currency != "" && len(ids) > 0 {
		must = append(must, elastic.NewMatchQuery("job_details.salarycurrency", data.Currency))
	}

	// search by known language
	if len(data.Language) > 0 {
		languages := make([]elastic.Query, 0, len(data.Language))
		for i := range data.Language {
			languages = append(languages, elastic.NewMatchQuery("job_details.required_languages", data.Language[i]))
		}
		must = append(must, elastic.NewBoolQuery().Must(languages...))
	}

	// search by skill
	if len(data.Skill) > 0 {
		skills := make([]elastic.Query, 0, len(data.Skill))
		for i := range data.Skill {
			skills = append(skills, elastic.NewMatchQuery("job_details.required_skills", data.Skill[i]))
		}
		must = append(must, elastic.NewBoolQuery().Must(skills...))
	}

	// search by salary interval
	if data.Period != "Any" {
		must = append(must, elastic.NewMatchQuery("job_details.salaryinterval", data.Period))
	}

	// search by degree
	if len(data.Degree) > 0 {
		degrees := make([]elastic.Query, 0, len(data.Degree))
		for i := range data.Degree {
			degrees = append(degrees, elastic.NewMatchPhraseQuery("job_details.required_educations", data.Degree[i]))
		}
		must = append(must, elastic.NewBoolQuery().Must(degrees...))
	}

	// search by job type
	if len(data.JobType) > 0 {
		jobTypes := make([]elastic.Query, 0, len(data.JobType))
		for i := range data.JobType {
			jobTypes = append(jobTypes, elastic.NewMatchQuery("job_details.employment_types", data.JobType[i]))
		}
		must = append(must, elastic.NewBoolQuery().Must(jobTypes...))
	}

	// search by range salary
	if data.MaxSalary != 0 && data.MinSalary != 0 {
		must = append(must,
			elastic.NewBoolQuery().Must(
				elastic.NewRangeQuery("job_details.salarymin").Gte(data.MinSalary),
				elastic.NewRangeQuery("job_details.salarymax").Lte(data.MaxSalary),
			),
		)
	} else {
		if data.MaxSalary != 0 {
			must = append(must, elastic.NewRangeQuery("job_details.salarymax").Lte(data.MaxSalary))
		}
		if data.MinSalary != 0 {
			must = append(must, elastic.NewRangeQuery("job_details.salarymin").Gte(data.MinSalary))
		}
	}

	if data.DatePosted != requests.DateEnumAnytime {
		switch data.DatePosted {
		case requests.DateEnumPast24Hours:
			must = append(
				must,
				elastic.NewRangeQuery("activation_date").
					Gte(time.Now().AddDate(0, 0, -1)),
			)
		case requests.DateEnumPastWeek:
			must = append(
				must,
				elastic.NewRangeQuery("activation_date").
					Gte(time.Now().AddDate(0, 0, -7)),
			)
		case requests.DateEnumPastMonth:
			must = append(
				must,
				elastic.NewRangeQuery("activation_date").
					Gte(time.Now().AddDate(0, -1, 0)),
			)
		}
	}

	// search by cover letter
	if data.WithoutCoverLetter {
		must = append(must, elastic.NewTermQuery("job_details.cover_letter", false)) // is it correct?
	}

	// search by experience level
	if data.ExperienceLevel != requests.ExperienceEnumExpericenUnknown {
		// must = append(must, elastic.NewRangeQuery("job_details.required_experience").Lte(data.ExperienceLevel))
		must = append(must, elastic.NewMatchQuery("job_details.required_experience", data.ExperienceLevel))
	}

	// search by industry
	if len(data.Industry) > 0 {
		industries := make([]elastic.Query, 0, len(data.Industry))
		for i := range data.Industry {
			industries = append(industries, elastic.NewHasParentQuery("company", elastic.NewTermQuery("industry.main", data.Industry[i])))
		}
		must = append(must, elastic.NewBoolQuery().Should(industries...))
	}

	// search by subindustry
	if len(data.Subindustry) > 0 {
		subindustries := make([]elastic.Query, 0, len(data.Subindustry))
		for i := range data.Subindustry {
			subindustries = append(subindustries, elastic.NewHasParentQuery("company", elastic.NewTermQuery("industry.sub", data.Subindustry[i])))
		}
		must = append(must, elastic.NewBoolQuery().Should(subindustries...))
	}

	// search by company name
	if len(data.CompanyName) > 0 {
		names := make([]elastic.Query, 0, len(data.CompanyName))
		for i := range data.CompanyName {
			// names = append(names, elastic.NewTermQuery("company_details.company_name", data.CompanyName[i]))
			names = append(names, elastic.NewHasParentQuery("company", elastic.NewMatchPhraseQuery("name", data.CompanyName[i])))
		}
		must = append(must, elastic.NewBoolQuery().Should(names...))
	}

	if data.CompanySize != company.SizeUnknown {
		must = append(must, elastic.NewHasParentQuery("company", elastic.NewTermQuery("size", data.CompanySize)))
		// must = append(must, elastic.NewTermQuery("company_details.company_size", data.CompanySize))
	}

	// search by city
	if len(data.CityID) > 0 {
		cities := make([]elastic.Query, 0, len(data.CityID))
		for i := range data.CityID {
			// cities = append(cities, elastic.NewMatchQuery("job_details.city", data.City[i]))
			cities = append(cities,
				elastic.NewBoolQuery().Must(
					elastic.NewMatchQuery("job_details.city", data.CityID[i]),
				),
			)
		}
		must = append(must, elastic.NewBoolQuery().Should(cities...))
	}

	// search by Country
	if len(data.Country) > 0 {
		countries := make([]elastic.Query, 0, len(data.Country))
		for i := range data.Country {
			// countries = append(countries, elastic.NewMatchQuery("job_metadata.advertisement_countries", data.Country[i]))
			countries = append(countries,
				elastic.NewBoolQuery().Must(
					elastic.NewMatchQuery("job_details.country", data.Country[i]),
				),
				// ),
			)
		}
		must = append(must, elastic.NewBoolQuery().Should(countries...))
	}

	// --------

	search := elastic.NewBoolQuery()
	if len(must) > 0 {
		search.Must(must...)
	}
	if len(should) > 0 {
		search.Should(should...)
	}
	if len(mustNot) > 0 {
		search.MustNot(mustNot...)
	}

	searchQuery := r.client.
		Search("companies-db.companies").
		Query(
			search,
		).
		SortBy(
			elastic.NewFieldSort("_score").Desc(),
		).
		FetchSource(false)

	searchResultAmount, err := searchQuery.Do(ctx)
	if err != nil {
		return []string{}, 0, err
	}

	// if data.After != "" {
	// 	searchQuery.SearchAfter(data.After)
	// }
	if data.After != "" {
		// searchQuery.SearchAfter(data.GetAfter())
		num, _ := strconv.Atoi(data.After)
		searchQuery.From(num)
	}

	if data.First != 0 {
		searchQuery.Size(int(data.First))
	}

	searchResult, err := searchQuery.Do(ctx)
	if err != nil {
		return []string{}, 0, err
	}

	// --------

	resultIDs := make([]string, 0, len(searchResult.Hits.Hits))

	for i := range searchResult.Hits.Hits {
		resultIDs = append(resultIDs, searchResult.Hits.Hits[i].Id)
	}

	return resultIDs, searchResultAmount.TotalHits(), nil
}

// CandidateSearch ...
func (r Repository) CandidateSearch(ctx context.Context, data *requests.CandidateSearch, blockedIDs []string) ([]string, []string, int64, error) {
	must := make([]elastic.Query, 0)
	should := make([]elastic.Query, 0)
	mustNot := make([]elastic.Query, 0)

	// search only among candidates who is opened
	must = append(must, elastic.NewTermQuery("is_open", true))

	// search only activated users
	must = append(must, elastic.NewHasParentQuery("user", elastic.NewTermQuery("status", "ACTIVATED")))

	// exclude blocked users
	if len(blockedIDs) > 0 {
		mustNot = append(mustNot, elastic.NewHasParentQuery(
			"user",
			elastic.NewIdsQuery("_doc").Ids(blockedIDs...),
		),
		)
	}

	// search by keyword
	if len(data.Keyword) > 0 {
		keywords := make([]elastic.Query, 0, len(data.Keyword))
		for i := range data.Keyword {
			keywords = append(keywords,
				elastic.NewMultiMatchQuery(
					data.Keyword[i],
					"jobs",
					"industry",
					"subindustries",
					"job_types",
				),
				elastic.NewHasParentQuery(
					"user",
					elastic.NewMultiMatchQuery(
						data.Keyword[i],
						"first_name", "last_name", "patronymic.name", "middlename.name",
						"native_name.name", "nickname.name", "headline", "story", "emails.email",
						"skills.skill", "educations.school", "educations.degree", "experiences.position",
						"interests.interest", "known_languages.language", "experiences.company",
						"educations.field_study",
					),
				),
			)
		}
		must = append(must, elastic.NewBoolQuery().Should(keywords...))
	}

	// search by skill
	if len(data.Skill) > 0 {
		skills := make([]elastic.Query, 0, len(data.Skill))
		for i := range data.Skill {
			skills = append(skills, elastic.NewHasParentQuery("user", elastic.NewMatchPhraseQuery("skills.skill", data.Skill[i])))
			skills = append(skills, elastic.NewHasParentQuery("user", elastic.NewQueryStringQuery(data.Skill[i]).Field("skills.translations.*.skill")))
		}
		must = append(must, elastic.NewBoolQuery().Should(skills...))
	}

	// search by degree
	if len(data.Degree) > 0 {
		degrees := make([]elastic.Query, 0, len(data.Degree))
		for i := range data.Degree {
			degrees = append(degrees, elastic.NewHasParentQuery("user", elastic.NewMatchQuery("educations.degree", data.Degree[i])))
			degrees = append(degrees, elastic.NewHasParentQuery("user", elastic.NewQueryStringQuery(data.Degree[i]).Field("educations.translations.*.degree")))
		}
		must = append(must, elastic.NewBoolQuery().Should(degrees...))
	}

	// search by field of study
	if len(data.FieldOfStudy) > 0 {
		fieldOfStudies := make([]elastic.Query, 0, len(data.FieldOfStudy))
		for i := range data.FieldOfStudy {
			fieldOfStudies = append(fieldOfStudies, elastic.NewHasParentQuery("user", elastic.NewMatchPhraseQuery("educations.field_study", data.FieldOfStudy[i])))
			fieldOfStudies = append(fieldOfStudies, elastic.NewHasParentQuery("user", elastic.NewQueryStringQuery(data.FieldOfStudy[i]).Field("educations.translations.*.field_study")))
		}
		must = append(must, elastic.NewBoolQuery().Should(fieldOfStudies...))
	}

	// search by school
	if len(data.School) > 0 {
		schools := make([]elastic.Query, 0, len(data.School))
		for i := range data.School {
			schools = append(schools, elastic.NewHasParentQuery("user", elastic.NewMatchQuery("educations.school", data.School[i])))
			schools = append(schools, elastic.NewHasParentQuery("user", elastic.NewQueryStringQuery(data.School[i]).Field("educations.translations.*.school")))
		}
		must = append(must, elastic.NewBoolQuery().Should(schools...))
	}

	// search by salary interval
	if data.Period != "Any" {
		must = append(must, elastic.NewMatchQuery("career_interests.salary_interval", data.Period))
	}

	// search by job type
	if len(data.JobType) > 0 {
		jobTypes := make([]elastic.Query, 0, len(data.JobType))
		for i := range data.JobType {
			jobTypes = append(jobTypes, elastic.NewMatchQuery("career_interests.job_types", data.JobType[i]))
		}
		must = append(must, elastic.NewBoolQuery().Should(jobTypes...))
	}

	// search by industry
	if len(data.Industry) > 0 {
		industry := make([]elastic.Query, 0, len(data.Industry))
		for i := range data.Industry {
			industry = append(industry, elastic.NewTermQuery("career_interests.industry", data.Industry[i]))
		}
		must = append(must, elastic.NewBoolQuery().Should(industry...))
	}

	// search by subindustry
	if len(data.SubIndustry) > 0 {
		subindustry := make([]elastic.Query, 0, len(data.SubIndustry))
		for i := range data.SubIndustry {
			subindustry = append(subindustry, elastic.NewTermQuery("career_interests.subindustries", data.SubIndustry[i]))
		}
		must = append(must, elastic.NewBoolQuery().Should(subindustry...))
	}

	// search by currency
	if data.Currency != "" {
		must = append(must, elastic.NewMatchQuery("career_interests.salary_currency", data.Currency))
	}

	// search by known language
	if len(data.Language) > 0 {
		languages := make([]elastic.Query, 0, len(data.Language))
		for i := range data.Language {
			languages = append(languages, elastic.NewHasParentQuery("user", elastic.NewMatchQuery("known_languages.language", data.Language[i])))
		}
		must = append(must, elastic.NewBoolQuery().Should(languages...))
	}

	// search by salary range
	if data.MaxSalary != 0 && data.MinSalary != 0 {
		must = append(must,
			elastic.NewBoolQuery().Must(
				elastic.NewRangeQuery("career_interests.salary_min").Gte(data.MinSalary),
				elastic.NewRangeQuery("career_interests.salary_max").Lte(data.MaxSalary),
			),
		)
	} else {
		if data.MaxSalary != 0 {
			must = append(must, elastic.NewRangeQuery("career_interests.salary_max").Lte(data.MaxSalary))
		}
		if data.MinSalary != 0 {
			must = append(must, elastic.NewRangeQuery("career_interests.salary_min").Gte(data.MinSalary))
		}
	}

	// search among student
	if data.IsStudent {
		must = append(must,
			elastic.NewBoolQuery().Should(
				elastic.NewHasParentQuery("user", elastic.NewRangeQuery("educations.finish_date").Gt(time.Now())),
				elastic.NewHasParentQuery("user", elastic.NewMatchQuery("educations.is_currenlty_study", true)),
			),
		)
	}

	// search by past companies
	if len(data.PastCompany) > 0 {
		pastCompanies := make([]elastic.Query, 0, len(data.PastCompany))
		for i := range data.PastCompany {
			pastCompanies = append(pastCompanies,
				elastic.NewNestedQuery("experiences",
					elastic.NewBoolQuery().Must(
						elastic.NewMatchPhraseQuery("experiences.company", data.PastCompany[i]),
						elastic.NewMatchQuery("experiences.currently_work", false),
					),
				),
			)
		}
		must = append(must, elastic.NewHasParentQuery("user", elastic.NewBoolQuery().Must(pastCompanies...)))
	}

	// search by current companies
	if len(data.CurrentCompany) > 0 {
		currentCompanies := make([]elastic.Query, 0, len(data.CurrentCompany))
		for i := range data.CurrentCompany {
			currentCompanies = append(currentCompanies,
				elastic.NewNestedQuery("experiences",
					elastic.NewBoolQuery().Must(
						elastic.NewMatchPhraseQuery("experiences.company", data.CurrentCompany[i]),
						elastic.NewMatchQuery("experiences.currently_work", true),
					),
				),
			)
		}
		must = append(must, elastic.NewHasParentQuery("user", elastic.NewBoolQuery().Must(currentCompanies...)))
	}

	// search among who is willing to travel
	if data.IsWillingToTravel {
		must = append(must, elastic.NewTermQuery("career_interests.travel", true))
	}

	// search among who can relocate
	if data.IsPossibleToRelocate {
		must = append(must, elastic.NewTermQuery("career_interests.relocate", true))
	}

	// search among who can work remotely
	if data.IsWillingToWorkRemotly {
		must = append(must, elastic.NewTermQuery("career_interests.remote", true))
	}

	// search by country
	if len(data.Country) > 0 {
		countries := make([]elastic.Query, 0, len(data.Country))
		for i := range data.Country {
			// countries = append(countries, elastic.NewHasParentQuery("user", elastic.NewTermQuery("location.country.id", data.Country[i])))
			// countries = append(countries,
			// 	elastic.NewBoolQuery().Must(
			countries = append(countries, elastic.NewMatchQuery("career_interests.locations.country", data.Country[i]))
			// 	),
			// )
		}
		must = append(must, elastic.NewBoolQuery().Should(countries...))
	}

	// search by city
	if len(data.CityID) > 0 {
		cities := make([]elastic.Query, 0, len(data.CityID))
		for i := range data.CityID {
			// cities = append(cities, elastic.NewHasParentQuery("user", elastic.NewTermQuery("location.city.id", data.CityID[i])))
			// cities = append(cities,
			// 	elastic.NewBoolQuery().Must(
			cities = append(cities, elastic.NewMatchQuery("career_interests.locations.city_id", data.CityID[i]))
			// 	),
			// )
		}
		must = append(must, elastic.NewBoolQuery().Should(cities...))
	}

	if data.ExperienceLevel != requests.ExperienceEnumExpericenUnknown {
		must = append(must, elastic.NewTermQuery("career_interests.experience", data.ExperienceLevel))
	}

	// ------------

	search := elastic.NewBoolQuery()
	if len(must) > 0 {
		search.Must(must...)
	}
	if len(should) > 0 {
		search.Should(should...)
	}
	if len(mustNot) > 0 {
		search.MustNot(mustNot...)
	}

	searchQuery := r.client.
		Search("users_db.users").
		Query(
			search,
		).
		SortBy(
			elastic.NewFieldSort("_score").Desc(),
		).
		FetchSource(false)

	searchResultAmount, err := searchQuery.Do(ctx)
	if err != nil {
		return []string{}, []string{}, 0, err
	}

	if data.After != "" {
		// searchQuery.SearchAfter(data.GetAfter())
		num, _ := strconv.Atoi(data.After)
		searchQuery.From(num)
	}

	if data.First != 0 {
		searchQuery.Size(int(data.First))
	}

	searchResult, err := searchQuery.Do(ctx)
	if err != nil {
		return []string{}, []string{}, 0, err
	}

	// -------------------

	resultIDs := make([]string, len(searchResult.Hits.Hits))
	parentsIDs := make([]string, len(searchResult.Hits.Hits))

	for i := range searchResult.Hits.Hits {
		resultIDs[i] = searchResult.Hits.Hits[i].Id
		parentsIDs[i] = searchResult.Hits.Hits[i].Routing
	}

	return resultIDs, parentsIDs, searchResultAmount.TotalHits(), nil
}

// CompanySearch ...
func (r Repository) CompanySearch(ctx context.Context, data *requests.CompanySearch, blockedIDs []string) ([]string, int64, error) {
	must := make([]elastic.Query, 0)
	should := make([]elastic.Query, 0)
	mustNot := make([]elastic.Query, 0)

	must = append(must, elastic.NewTermQuery("document_type", "company"))
	must = append(must, elastic.NewMatchQuery("status", "ACTIVATED")) // TODO: should be term query

	if len(blockedIDs) > 0 {
		mustNot = append(mustNot, elastic.NewIdsQuery("_doc").Ids(blockedIDs...))
	}

	// search by keyword
	if len(data.Keyword) > 0 {
		keywords := make([]elastic.Query, 0, len(data.Keyword))
		for i := range data.Keyword {
			keywords = append(keywords, elastic.NewMultiMatchQuery(
				data.Keyword[i],
				"name",
				"addresses.name", "addresses.street", "addresses.website",
				"emails.email",
				"vat",
				"websites.website",
				"description",
				"mission",
				// "founders."
				"awards.title",
			),
			)
		}
		must = append(must, elastic.NewBoolQuery().Should(keywords...))
	}

	// search by name
	if len(data.Name) > 0 {
		names := make([]elastic.Query, 0, len(data.Name))
		for i := range data.Name {
			names = append(names, elastic.NewMatchQuery("name", data.Name[i]))
		}
		must = append(must, elastic.NewBoolQuery().Should(names...))
	}

	// search by city
	if len(data.CityID) > 0 {
		cities := make([]elastic.Query, 0, len(data.CityID))
		for i := range data.CityID {
			// cities = append(cities, elastic.NewMatchQuery("addresses.location.city.id", data.City[i]))
			cities = append(cities,
				elastic.NewBoolQuery().Must(
					elastic.NewMatchQuery("addresses.location.city.id", data.CityID[i]),
					elastic.NewMatchQuery("addresses.is_primary", true),
				),
			)
		}
		must = append(must, elastic.NewBoolQuery().Should(cities...))
	}

	// search by country
	if len(data.Country) > 0 {
		countries := make([]elastic.Query, 0, len(data.Country))
		for i := range data.Country {
			// countries = append(countries, elastic.NewMatchQuery("addresses.location.country.id", data.Country[i]))
			countries = append(countries,
				elastic.NewBoolQuery().Must(
					elastic.NewMatchQuery("addresses.location.country.id", data.Country[i]),
					elastic.NewMatchQuery("addresses.is_primary", true),
				),
				// ),
			)
		}
		must = append(must, elastic.NewBoolQuery().Should(countries...))
	}

	if data.Type != company.TypeUnknown {
		must = append(must, elastic.NewTermQuery("type", data.Type))
	}

	if data.Size != company.SizeUnknown {
		must = append(must, elastic.NewTermQuery("size", data.Size))
	}

	// search by industry
	if len(data.Industry) > 0 {
		industries := make([]elastic.Query, 0, len(data.Industry))
		for i := range data.Industry {
			industries = append(industries, elastic.NewMatchQuery("industry.main", data.Industry[i]))
		}
		must = append(must, elastic.NewBoolQuery().Should(industries...))
	}

	// search by subindustry
	if len(data.SubIndustry) > 0 {
		industries := make([]elastic.Query, 0, len(data.SubIndustry))
		for i := range data.SubIndustry {
			industries = append(industries, elastic.NewMatchQuery("industry.sub", data.SubIndustry[i]))
		}
		must = append(must, elastic.NewBoolQuery().Should(industries...))
	}

	// search by founders id
	if len(data.FoundersID) > 0 {
		founders := make([]elastic.Query, 0, len(data.FoundersID))
		for i := range data.FoundersID {
			founders = append(founders, elastic.NewMatchQuery("founders.user_id", data.FoundersID[i]))
		}
		must = append(must, elastic.NewBoolQuery().Should(founders...))
	}

	// search by founders name
	if len(data.FoundersName) > 0 {
		founders := make([]elastic.Query, 0, len(data.FoundersName))
		for i := range data.FoundersName {
			founders = append(founders, elastic.NewMatchQuery("founders.name", data.FoundersName[i]))
		}
		must = append(must, elastic.NewBoolQuery().Should(founders...))
	}

	// search by bussiness hours
	if len(data.BusinessHours) > 0 {
		bhs := make([]elastic.Query, 0, len(data.BusinessHours))
		for i := range data.BusinessHours {
			log.Println("bs:", data.BusinessHours)
			bhs = append(bhs, elastic.NewMatchQuery("business_hours.weekdays", data.BusinessHours[i]))
		}
		must = append(must, elastic.NewBoolQuery().Must(bhs...))
	}

	if data.IsJobOffers {
		must = append(must, elastic.NewHasChildQuery("job", elastic.NewMatchAllQuery()))
	}

	if data.IsCareerCenterOpenened {
		must = append(must, elastic.NewTermQuery("career_center.is_opened", true))
	}

	// TODO:
	// data.GetIsCompany()      // bool
	// data.GetIsOrganization() // bool

	// data.GetRating()

	// ------------

	search := elastic.NewBoolQuery()
	if len(must) > 0 {
		search.Must(must...)
	}
	if len(should) > 0 {
		search.Should(should...)
	}
	if len(mustNot) > 0 {
		search.MustNot(mustNot...)
	}

	searchQuery := r.client.
		Search("companies-db.companies").
		Query(
			search,
		).
		SortBy(
			elastic.NewFieldSort("_score").Desc(),
		).
		FetchSource(false)

	searchResultAmount, err := searchQuery.Do(ctx)
	if err != nil {
		log.Println("error: sending search:", err)
		return []string{}, 0, err
	}

	if data.After != "" {
		// searchQuery.SearchAfter(data.GetAfter())
		num, _ := strconv.Atoi(data.After)
		searchQuery.From(num)
	}

	if data.First != 0 {
		searchQuery.Size(int(data.First))
	}

	searchResult, err := searchQuery.Do(ctx)
	if err != nil {
		log.Println("error: sending search:", err)
		return []string{}, 0, err
	}

	ids := make([]string, 0, len(searchResult.Hits.Hits))
	for i := range searchResult.Hits.Hits {
		ids = append(ids, searchResult.Hits.Hits[i].Id)
	}

	return ids, searchResultAmount.TotalHits(), nil
}

// ServiceSearch ...
func (r Repository) ServiceSearch(ctx context.Context, data *requests.ServiceSearch) ([]string, int64, error) {
	must := make([]elastic.Query, 0)
	should := make([]elastic.Query, 0)
	mustNot := make([]elastic.Query, 0)

	must = append(must, elastic.NewTermQuery("is_draft", false))
	must = append(must, elastic.NewTermQuery("is_paused", false))

	// must = append(must, elastic.NewTermQuery("is_remote", data.Remote))

	log.Printf("Data in Search %+v", data)

	if len(data.Keyword) > 0 {
		keywords := make([]elastic.Query, 0, len(data.Keyword))
		for i := range data.Keyword {
			keywords = append(keywords, elastic.NewMultiMatchQuery(
				data.Keyword[i],
				"title",
				"description",
			))
		}
		must = append(must, elastic.NewBoolQuery().Should(keywords...))
	}

	// search by country
	if len(data.CountryID) > 0 {
		countries := make([]elastic.Query, 0, len(data.CountryID))
		for i := range data.CountryID {
			countries = append(countries, elastic.NewMatchQuery("location.country.id", data.CountryID[i]))
		}
		must = append(must, elastic.NewBoolQuery().Should(countries...))
	}

	// search by city
	if len(data.CityID) > 0 {
		cities := make([]elastic.Query, 0, len(data.CityID))
		for i := range data.CityID {
			cities = append(cities, elastic.NewMatchQuery("locations.city.id", data.CityID[i]))
		}
		must = append(must, elastic.NewBoolQuery().Should(cities...))
	}

	if data.DeliveryTime != requests.DeliveryAny {
		must = append(must, elastic.NewMatchQuery("delivery_time", data.DeliveryTime))
	}

	if data.LocationType != requests.LocationTypeAny {
		must = append(must, elastic.NewMatchQuery("location_type", data.LocationType))
	}

	if data.Price != requests.PriceAny {
		must = append(must, elastic.NewMatchQuery("price", data.Price))
	}

	if data.FixedPrice != 0 {
		must = append(must, elastic.NewMatchQuery("fixed_price_amount", data.FixedPrice))
	}

	// search by range salary
	if data.MinPrice != 0 && data.MaxPrice != 0 {
		must = append(must,
			elastic.NewBoolQuery().Must(
				elastic.NewRangeQuery("min_price_amount").Gte(data.MinPrice),
				elastic.NewRangeQuery("max_price_amount").Lte(data.MaxPrice),
			),
		)
	} else {
		if data.MaxPrice != 0 {
			must = append(must, elastic.NewRangeQuery("max_price_amount").Lte(data.MaxPrice))
		}
		if data.MinPrice != 0 {
			must = append(must, elastic.NewRangeQuery("min_price_amount").Gte(data.MinPrice))
		}
	}

	// currency
	if data.CurrencyPrice != "" {
		must = append(must,
			elastic.NewMatchQuery("currency", data.CurrencyPrice),
		)
	}

	// skills
	if len(data.Skills) > 0 {
		for i := range data.Skills {
			must = append(must,
				elastic.NewMatchQuery("additional_details.qualifications.skills.skill", data.Skills[i]),
			)
		}
	}

	// languages
	if len(data.Languages) > 0 {
		for i := range data.Languages {
			must = append(must,
				elastic.NewMatchQuery("additional_details.qualifications.languages", data.Languages[i]),
			)
		}
	}

	// Working Hours --> Is Always open
	if data.IsAlwaysOpen {
		must = append(must,
			elastic.NewTermQuery("working_hours.is_always_open", data.IsAlwaysOpen),
		)

	}

	// Working Hours --> Weekdays
	if len(data.WeekDays) > 0 {
		for i := range data.WeekDays {
			must = append(must,
				elastic.NewMatchQuery("working_hours.working_date.week_days", data.WeekDays[i]),
			)
		}

	}

	// Working Hours --> Hour From
	if data.HourFrom != "" {
		must = append(must,
			elastic.NewMatchQuery("working_hours.working_date.hour_from", data.HourFrom),
		)

	}

	// Working Hours --> Hour To
	if data.HourTo != "" {
		must = append(must,
			elastic.NewMatchQuery("working_hours.working_date.hour_to", data.HourTo),
		)

	}

	// Service Onwner
	if data.ServiceOwner != requests.ServiceOwnerAny {
		// Service owner user
		if data.ServiceOwner == requests.ServiceOwnerUser {
			must = append(must, elastic.NewBoolQuery().Should(
				elastic.NewExistsQuery("user_id"),
			))
		} else {
			// Service owner company
			must = append(must, elastic.NewBoolQuery().Should(
				elastic.NewExistsQuery("company_id"),
			))
		}
	}

	// ------------

	search := elastic.NewBoolQuery()
	if len(must) > 0 {
		search.Must(must...)
	}
	if len(should) > 0 {
		search.Should(should...)
	}
	if len(mustNot) > 0 {
		search.MustNot(mustNot...)
	}

	searchQuery := r.client.
		Search("services_db.services").
		Query(
			search,
		).
		SortBy(
			elastic.NewFieldSort("_score").Desc(),
		).
		FetchSource(false)

	searchResultAmount, err := searchQuery.Do(ctx)
	if err != nil {
		log.Println("error: sending search:", err)
		return []string{}, 0, err
	}

	if data.After != "" {
		num, _ := strconv.Atoi(data.After)
		searchQuery.From(num)
	}

	if data.First != 0 {
		searchQuery.Size(int(data.First))
	}

	searchResult, err := searchQuery.Do(ctx)
	if err != nil {
		log.Println("error: sending search:", err)
		return []string{}, 0, err
	}

	ids := make([]string, 0, len(searchResult.Hits.Hits))
	for i := range searchResult.Hits.Hits {
		ids = append(ids, searchResult.Hits.Hits[i].Id)
	}

	return ids, searchResultAmount.TotalHits(), nil
}

// ServiceRequestSearch ...
func (r Repository) ServiceRequestSearch(ctx context.Context, data *requests.ServiceRequest) ([]string, int64, error) {
	must := make([]elastic.Query, 0)
	should := make([]elastic.Query, 0)
	mustNot := make([]elastic.Query, 0)

	must = append(must, elastic.NewTermQuery("is_draft", false))
	must = append(must, elastic.NewTermQuery("is_closed", false))
	must = append(must, elastic.NewTermQuery("is_paused", false))

	if len(data.Keyword) > 0 {
		keywords := make([]elastic.Query, 0, len(data.Keyword))
		for i := range data.Keyword {
			keywords = append(keywords, elastic.NewMultiMatchQuery(
				data.Keyword[i],
				"tittle",
				"description",
			))
		}
		must = append(must, elastic.NewBoolQuery().Should(keywords...))
	}

	// search by country
	if len(data.CountryID) > 0 {
		countries := make([]elastic.Query, 0, len(data.CountryID))
		for i := range data.CountryID {
			countries = append(countries, elastic.NewMatchQuery("location.country.id", data.CountryID[i]))
		}
		must = append(must, elastic.NewBoolQuery().Should(countries...))
	}

	// search by city
	if len(data.CityID) > 0 {
		cities := make([]elastic.Query, 0, len(data.CityID))
		for i := range data.CityID {
			cities = append(cities, elastic.NewMatchQuery("locations.city.id", data.CityID[i]))
		}
		must = append(must, elastic.NewBoolQuery().Should(cities...))
	}

	if data.DeliveryTime != requests.DeliveryAny {
		must = append(must, elastic.NewMatchQuery("delivery_time", data.DeliveryTime))
	}

	if data.LocationType != requests.LocationTypeAny {
		must = append(must, elastic.NewMatchQuery("location_type", data.LocationType))
	}

	if data.FixedPrice != 0 {
		must = append(must, elastic.NewMatchQuery("fixed_price_amount", data.FixedPrice))
	}

	if data.PriceType != requests.PriceAny {
		must = append(must, elastic.NewMatchQuery("price", data.PriceType))
	}

	// search by range salary
	if data.MinPrice != 0 && data.MaxPrice != 0 {
		must = append(must,
			elastic.NewBoolQuery().Must(
				elastic.NewRangeQuery("min_price_amount").Gte(data.MinPrice),
				elastic.NewRangeQuery("max_price_amount").Lte(data.MaxPrice),
			),
		)
	} else {
		if data.MaxPrice != 0 {
			must = append(must, elastic.NewRangeQuery("max_price_amount").Lte(data.MaxPrice))
		}
		if data.MinPrice != 0 {
			must = append(must, elastic.NewRangeQuery("min_price_amount").Gte(data.MinPrice))
		}
	}

	// currency
	if data.CurrencyPrice != "" {
		must = append(must,
			elastic.NewMatchQuery("currency", data.CurrencyPrice),
		)
	}

	// project type
	if len(data.ProjectType) > 0 {
		for i := range data.ProjectType {
			if data.ProjectType[i] != requests.ProjectTypeNotSure {
				must = append(must, elastic.NewMatchQuery("project_type", data.ProjectType[i]))
			}
		}
	}

	// skills
	if len(data.Skills) > 0 {
		for i := range data.Skills {
			must = append(must,
				elastic.NewMatchQuery("additional_details.qualifications.skills.skill", data.Skills[i]),
			)
		}
	}

	// tools
	if len(data.Tools) > 0 {
		for i := range data.Tools {
			must = append(must,
				elastic.NewMatchQuery("additional_details.qualifications.tool_technology.tool-technology", data.Tools[i]),
			)
		}
	}

	// languages
	if len(data.Languages) > 0 {
		for i := range data.Languages {
			must = append(must,
				elastic.NewMatchQuery("additional_details.qualifications.languages.language", data.Languages[i]),
			)
		}
	}

	// Service Onwner
	if data.ServiceOwner != requests.ServiceOwnerAny {
		// Service owner user
		if data.ServiceOwner == requests.ServiceOwnerUser {
			must = append(must, elastic.NewBoolQuery().Should(
				elastic.NewExistsQuery("user_id"),
			))
		} else {
			// Service owner company
			must = append(must, elastic.NewBoolQuery().Should(
				elastic.NewExistsQuery("company_id"),
			))
		}
	}

	// ------------

	search := elastic.NewBoolQuery()
	if len(must) > 0 {
		search.Must(must...)
	}
	if len(should) > 0 {
		search.Should(should...)
	}
	if len(mustNot) > 0 {
		search.MustNot(mustNot...)
	}

	searchQuery := r.client.
		Search("services_db.requests").
		Query(
			search,
		).
		SortBy(
			elastic.NewFieldSort("_score").Desc(),
		).
		FetchSource(false)

	searchResultAmount, err := searchQuery.Do(ctx)
	if err != nil {
		log.Println("error: sending search:", err)
		return []string{}, 0, err
	}

	if data.After != "" {
		num, _ := strconv.Atoi(data.After)
		searchQuery.From(num)
	}

	if data.First != 0 {
		searchQuery.Size(int(data.First))
	}

	searchResult, err := searchQuery.Do(ctx)
	if err != nil {
		log.Println("error: sending search:", err)
		return []string{}, 0, err
	}

	ids := make([]string, 0, len(searchResult.Hits.Hits))
	for i := range searchResult.Hits.Hits {
		ids = append(ids, searchResult.Hits.Hits[i].Id)
	}

	return ids, searchResultAmount.TotalHits(), nil
}
