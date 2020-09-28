package network

import (
	"fmt"
	"strings"

	"gitlab.lan/Rightnao-site/microservices/network/pkg/model"
)

func filter_collection(coll, objectName string) string {
	return fmt.Sprintf("FILTER IS_SAME_COLLECTION('%s', %s)", coll, objectName)
}

func filter_not_blocked(name1, name2 string) string {
	return fmt.Sprintf("FILTER COUNT(for b in blocks filter (b._from == %[1]s._id AND b._to == %[2]s._id) OR (b._to == %[1]s._id AND b._from == %[2]s._id) return 1) == 0", name1, name2)
}

func filter_users_by_companies(userName string, companies []string) string {
	return fmt.Sprintf("FILTER COUNT(for w in works_at filter w._from == %s._id AND w._to in @companies return 1) > 0 ", userName)
}

func filter_users_by_query(userName string) string {
	return fmt.Sprintf("FILTER (CONTAINS(LOWER(%[1]s.first_name), @query) OR CONTAINS(LOWER(%[1]s.last_name), @query))", userName)
}

func filter_companies_by_query(companyName string) string {
	return fmt.Sprintf("FILTER CONTAINS(LOWER(%s.name), @query)", companyName)
}

func filter_by_category(me, other string) string {
	return fmt.Sprintf("FILTER LENGTH(for v IN have_in_category filter v._from == %s._id AND v._to == %s._id AND v.name == @category return 1) > 0", me, other)
}

func filter_by_followings_category(me, other string) string {
	return fmt.Sprintf("FILTER LENGTH(for v IN have_in_category_followings filter v._from == %s._id AND v._to == %s._id AND v.name == @category return 1) > 0", me, other)
}

func filter_users_by_first_letter(userName string) string {
	return fmt.Sprintf("FILTER (LOWER(LEFT(%[1]s.first_name, 1)) == @letter OR LOWER(LEFT(%[1]s.last_name, 1)) == @letter)", userName)
}
func filter_companies_by_first_letter(companyName string) string {
	return fmt.Sprintf("FILTER LOWER(LEFT(%s.name, 1)) == @letter", companyName)
}

func create_users_filter(me, other string, filterParams *model.UserFilter) (string, map[string]interface{}) {
	filter := ""
	params := make(map[string]interface{})
	if len(filterParams.Companies) > 0 {
		filter = fmt.Sprint(filter, "\n", filter_users_by_companies(other, filterParams.Companies))
		params["companies"] = filterParams.Companies
	}
	if len(filterParams.Query) > 0 {
		filter = fmt.Sprint(filter, "\n", filter_users_by_query(other))
		params["query"] = strings.ToLower(filterParams.Query)
	}
	if len(filterParams.Category) > 0 {
		filter = fmt.Sprint(filter, "\n", filter_by_category(me, other))
		params["category"] = filterParams.Category
	}

	if len(filterParams.Letter) > 0 {
		filter = fmt.Sprint(filter, "\n", filter_users_by_first_letter(other))
		params["letter"] = strings.ToLower(filterParams.Letter)
	}

	return filter, params
}

func create_companies_filter(company, other string, filterParams *model.CompanyFilter) (string, map[string]interface{}) {
	filter := ""
	params := make(map[string]interface{})
	if len(filterParams.Query) > 0 {
		filter = fmt.Sprint(filter, "\n", filter_companies_by_query(other))
		params["query"] = strings.ToLower(filterParams.Query)
	}
	if len(filterParams.Category) > 0 {
		filter = fmt.Sprint(filter, "\n", filter_by_followings_category(company, other))
		params["category"] = filterParams.Category
	}

	if len(filterParams.Letter) > 0 {
		filter = fmt.Sprint(filter, "\n", filter_companies_by_first_letter(other))
		params["letter"] = strings.ToLower(filterParams.Letter)
	}

	return filter, params
}

func count_friends(varName, objName string) string {
	return fmt.Sprintf("let %s = COUNT(for v, f in 1..1 ANY %s friendship filter f.status == 'Approved' FILTER v.status == 'ACTIVATED' return 1)", varName, objName)
}

func count_followers(varName, objName string) string {
	return fmt.Sprintf("let %s = COUNT(for v, u IN 1..1 INBOUND %s follow FILTER v.status == 'ACTIVATED' return 1)", varName, objName)
}

func count_followings(varName, objName string) string {
	return fmt.Sprintf("let %s = COUNT(for v, u IN 1..1 OUTBOUND %s follow FILTER v.status == 'ACTIVATED' return 1)", varName, objName)
}

func count_received_recommendations(varName, objName string) string {
	return fmt.Sprintf("let %s = COUNT(for v IN 1..1 INBOUND %s recommendations return 1)", varName, objName)
}

func new_document(name, param string) string {
	return fmt.Sprintf("let %s = DOCUMENT(@%s)", name, param)
}

func is_following(varName, followerName, followingName string) string {
	return fmt.Sprintf("let %s = COUNT(for f in follow filter f._from == %s._id AND f._to == %s._id return 1) > 0", varName, followerName, followingName)
}

func is_friend(varName, myName, othersName string) string {
	return fmt.Sprintf("let %[1]s = COUNT(for f in friendship filter (f._from == %[2]s._id AND f._to == %[3]s._id AND f.status == 'Approved') or (f._to == %[2]s._id AND f._from == %[3]s._id AND f.status == 'Approved') return 1) > 0", varName, myName, othersName)
}

func is_blocked(varName, myName, othersName string) string {
	return fmt.Sprintf("let %s = COUNT(for f in blocks filter f._from == %s._id AND f._to == %s._id return 1) > 0", varName, myName, othersName)
}

func have_in_category(varName, myName, othersName string) string {
	return fmt.Sprintf("let %s = COUNT(for v IN have_in_category filter v._from == %s._id AND v._to == %s._id AND v.name == @category return 1) > 0", varName, myName, othersName)
}

func sort_users(userName, relationName, sortBy string) string {
	if sortBy == "first_name" || sortBy == "last_name" {
		return fmt.Sprintf("SORT %s.%s", userName, sortBy)
	} else if sortBy == "recently_added" {
		return fmt.Sprintf("SORT %s.created_at desc", relationName)
	}
	return ""
}

// TODO: add sort by "recently_added"
func sort_companies(companyName, sortBy string) string {
	if sortBy == "name" {
		return fmt.Sprintf("SORT %s.name", companyName)
	} else if sortBy == "rating" || sortBy == "size" || sortBy == "number_of_followers" {
		return fmt.Sprintf("SORT %s DESC", sortBy)
	}
	return ""
}

func following_of_followings(varName, myName string) string {
	return fmt.Sprintf(`let %s = (
for user in 2..2 OUTBOUND %s follow
	return user
)`, varName, myName)
}

func workers(varName, companyName string) string {
	return fmt.Sprintf(`let %s = (
for u in 1..1 INBOUND %s works_at
	return u
)`, varName, companyName)
}

func paginate(pagination *model.Pagination) string {
	if pagination.Amount > 0 {
		return fmt.Sprint("limit ", pagination.After, ", ", pagination.Amount)
	}
	return ""
}
