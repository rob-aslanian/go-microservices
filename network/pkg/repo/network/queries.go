package network

import (
	"fmt"
	"strings"

	"gitlab.lan/Rightnao-site/microservices/network/pkg/model"
)

const GET_FRIENDSHIP_QUERY = `
let friendship = DOCUMENT(@id)
let friend = DOCUMENT(@myId == friendship._from ? friendship._to : friendship._from)
return {friend: friend, _key: friendship._key, status: friendship.status, my_request: friendship._from == @myId, created_at: friendship.created_at}
`

var IS_FRIEND_QUERY = fmt.Sprint(
	"let user1 = DOCUMENT(@user1)\n",
	"let user2 = DOCUMENT(@user2)\n",
	is_friend("isFriend", "user1", "user2"), "\n",
	`return isFriend`,
)

func GET_FILTERED_FRIENDSHIP_REQUEST_QUERY(id string, filter *model.FriendshipRequestFilter) (string, map[string]interface{}) {
	params := make(map[string]interface{})

	var direction string
	if filter.Sent {
		direction = "OUTBOUND"
	} else {
		direction = "INBOUND"
	}

	if filter.Status == "" {
		filter.Status = string(model.FriendshipStatus_Requested)
	}

	params["id"] = id
	params["status"] = filter.Status

	return fmt.Sprintf(`
WITH users
let me = DOCUMENT(@id)
for friend,friendship in 1..1 %s me friendship
  FILTER friend.status == "ACTIVATED"
	FILTER COUNT(for b in blocks filter (b._from == friend._id AND b._to == me._id) OR (b._to == friend._id AND b._from == me._id) return 1) == 0
	FILTER friendship.status == @status
	return {friend,  _key: friendship._key, status: friendship.status, my_request: friendship._from == @id, description: friendship.description, created_at: friendship.created_at, responded_at: friendship.responded_at}
`, direction), params
}

func GET_ALL_FRIENDSHIP_QUERY(id string, filter *model.FriendshipRequestFilter) (string, map[string]interface{}) {
	params := make(map[string]interface{})

	// var direction string
	// if filter.Sent {
	// 	direction = "OUTBOUND"
	// } else {
	// 	direction = "INBOUND"
	// }
	//
	// if filter.Status == "" {
	// 	filter.Status = string(model.FriendshipStatus_Requested)
	// }

	params["id"] = id

	return fmt.Sprintf(`
WITH users
let me = DOCUMENT(@id)
for friend,friendship in 1..1 ANY me friendship
	FILTER friend.status == "ACTIVATED"
	FILTER COUNT(for b in blocks filter (b._from == friend._id AND b._to == me._id) OR (b._to == friend._id AND b._from == me._id) return 1) == 0
	return {friend,  _key: friendship._key, status: friendship.status, my_request: friendship._from == @id, description: friendship.description, created_at: friendship.created_at, responded_at: friendship.responded_at}
`), params
}

func GET_ALL_FRIENDSHIPID_QUERY(id string) (string, map[string]interface{}) {
	params := make(map[string]interface{})
	params["id"] = id

	return fmt.Sprintf(`
	WITH users
	LET me = DOCUMENT(@id)
	FOR friend,friendship IN 1..1 ANY me friendship
		FILTER friend.status == "ACTIVATED"
	 	FILTER COUNT(FOR b IN blocks FILTER (b._from == friend._id AND b._to == me._id) OR (b._to == friend._id AND b._from == me._id) return 1) == 0
	 	FILTER friendship.status == "Approved"
	 	RETURN {ids: friend._key}
`), params
}

func GET_FILTERED_FRIENDSHIP_QUERY(id, query, category, letter, sortBy string, companies []string) (string, map[string]interface{}) {
	filter := ""
	params := make(map[string]interface{})
	if len(companies) > 0 {
		filter = " COUNT(for w in works_at filter w._from == friend._id AND w._to in @companies return 1) > 0 "
		params["companies"] = companies
	}
	if len(query) > 0 {
		if len(filter) > 0 {
			filter = fmt.Sprint(filter, " AND ")
		}
		// filter = fmt.Sprint(filter, "(CONTAINS(LOWER(friend.first_name), @query) OR CONTAINS(LOWER(friend.last_name), @query)) ")
		filter = fmt.Sprint(filter, `(CONTAINS(LOWER(CONCAT(friend.first_name, " " ,friend.last_name)), @query)) OR (CONTAINS(LOWER(CONCAT(friend.last_name, " " ,friend.first_name)), @query))`)
		params["query"] = strings.ToLower(query)
	}
	if len(category) > 0 {
		if len(filter) > 0 {
			filter = fmt.Sprint(filter, " AND ")
		}
		if category == "not_categorised" {
			filter = fmt.Sprint(filter, "LENGTH(for v IN have_in_category filter v._from == @id AND v._to == friend._id return 1) == 0 ")
		} else {
			filter = fmt.Sprint(filter, "LENGTH(for v IN have_in_category filter v._from == @id AND v._to == friend._id AND CONTAINS(v.name, @category) return 1) > 0 ")
			params["category"] = category
		}
	}
	if len(letter) > 0 {
		if len(filter) > 0 {
			filter = fmt.Sprint(filter, " AND ")
		}
		filter = fmt.Sprint(filter, "(LOWER(LEFT(friend.first_name, 1)) == @letter OR LOWER(LEFT(friend.last_name, 1)) == @letter)")
		params["letter"] = strings.ToLower(letter)
	}
	if len(filter) > 0 {
		filter = fmt.Sprint("FILTER ", filter)
	}
	sortQuery := ""
	if len(sortBy) > 0 {
		if sortBy == "first_name" || sortBy == "last_name" {
			sortQuery = fmt.Sprint("SORT friend.", sortBy)
		} else if sortBy == "recently_added" {
			sortQuery = fmt.Sprint("SORT friendship.responded_at desc")
		}
	}

	params["id"] = id

	return fmt.Sprintf(`
WITH users, companies
let me = DOCUMENT(@id)
for friend,friendship in 1..1 ANY me friendship
	FILTER friendship.status == "Approved"
	FILTER friend.status == "ACTIVATED"
	FILTER COUNT(for b in blocks filter (b._from == friend._id AND b._to == me._id) OR (b._to == friend._id AND b._from == me._id) return 1) == 0
	%s
	let cats = (for cat in have_in_category FILTER cat._from == me._id AND cat._to == friend._id return cat.name)
	let isFollowing = COUNT(for fol in follow filter fol._from == me._id AND fol._to == friend._id return 1) > 0
	%s
	return {friend,
			_key: friendship._key,
			status: friendship.status,
			my_request: friendship._from == @id,
			description: friendship.description,
			categories: cats,
			following: isFollowing,
			created_at: friendship.created_at,
			responded_at: friendship.responded_at
	}
`, filter, sortQuery), params
}

const REMOVE_UNIDIRECTIONAL_RELATION_QUERY = `
WITH users
for rel in @@relation
	filter (rel._from == @user1 && rel._to == @user2) || (rel._from == @user2 && rel._to == @user1)
	remove rel in @@relation
`

const REMOVE_DIRECTIONAL_RELATION_QUERY = `
WITH users
for rel in @@relation
	filter rel._from == @from && rel._to == @to
	remove rel in @@relation
`

const REMOVE_FROM_CATEGORY_QUERY = `
WITH users
for rel in have_in_category
	filter rel._from == @from && rel._to == @to && rel.name == @category
	remove rel in have_in_category
`

const REMOVE_FROM_FOLLOWINGS_CATEGORY_QUERY = `
WITH users
for rel in have_in_category_followings
	filter rel._from == @from && rel._to == @to && rel.name == @category
	remove rel in have_in_category_followings
`

func BATCH_REMOVE_FROM_CATEGORY_QUERY(key string, userIds []string, category string, all bool) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"id":      ToUserId(key),
		"userIds": userIds,
	}
	query := fmt.Sprint(
		"WITH users\n",
		"for rel in have_in_category\n",
		"filter rel._from == @id && rel._to in @userIds\n",
	)
	if !all {
		params["category"] = category
		query = fmt.Sprint(query,
			"filter rel.name == @category\n",
		)
	}
	query = fmt.Sprint(query,
		"remove rel in have_in_category",
	)
	return query, params
}

func BATCH_REMOVE_FROM_CATEGORY_FOR_FOLLOWINGS_QUERY(key string, referalIds []string, category string, all bool) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"id":         ToUserId(key),
		"referalIds": referalIds,
	}
	query := fmt.Sprint(
		"WITH users\n",
		"for rel in have_in_category_followings\n",
		"filter rel._from == @id && rel._to in @referalIds\n",
	)
	if !all {
		params["category"] = category
		query = fmt.Sprint(query,
			"filter rel.name == @category\n",
		)
	}
	query = fmt.Sprint(query,
		"remove rel in have_in_category_followings",
	)
	return query, params
}

func GET_FILTERED_USER_FOLLOWINGS(key string, filter *model.UserFilter) (string, map[string]interface{}) {
	usersFilter, params := create_users_filter("me", "user", filter)
	params["id"] = ToUserId(key)
	result := fmt.Sprint(
		"WITH users, companies", "\n",
		new_document("me", "id"), "\n",
		"for user, fol in 1..1 OUTBOUND me follow", "\n",
		filter_collection(UsersName, "user"), "\n",
		`FILTER user.status == "ACTIVATED"`, "\n",
		filter_not_blocked("me", "user"), "\n",
		usersFilter, "\n",
		count_followers("followers", "user"), "\n",
		is_friend("is_friend", "me", "user"), "\n",
		sort_users("user", "fol", filter.SortBy), "\n",
		"return MERGE(fol, {user, followers, following: true, is_friend})",
	)
	return result, params
}

func GET_FILTERED_USER_FOLLOWERS(key string, filter *model.UserFilter) (string, map[string]interface{}) {
	usersFilter, params := create_users_filter("me", "user", filter)
	params["id"] = ToUserId(key)
	result := fmt.Sprint(
		"WITH users, companies", "\n",
		new_document("me", "id"), "\n",
		"for user, fol in 1..1 INBOUND me follow", "\n",
		filter_collection(UsersName, "user"), "\n",
		`FILTER user.status == "ACTIVATED"`, "\n",
		filter_not_blocked("me", "user"), "\n",
		usersFilter, "\n",
		count_followers("followers", "user"), "\n",
		is_friend("is_friend", "me", "user"), "\n",
		is_following("following", "me", "user"), "\n",
		sort_users("user", "fol", filter.SortBy), "\n",
		"return MERGE(fol, {user, followers, following, is_friend})",
	)
	return result, params
}

func GET_FILTERED_COMPANY_FOLLOWINGS(key string, filter *model.CompanyFilter) (string, map[string]interface{}) {
	compsFilter, params := create_companies_filter("me", "comp", filter)
	params["id"] = ToUserId(key)
	result := fmt.Sprint(
		"WITH users, companies", "\n",
		new_document("me", "id"), "\n",
		"for comp, fol in 1..1 OUTBOUND me follow", "\n",
		filter_collection(CompaniesName, "comp"), "\n",
		`FILTER comp.status == "ACTIVATED"`, "\n",
		filter_not_blocked("me", "comp"), "\n",
		compsFilter, "\n",
		count_followers("followers", "comp"), "\n",
		"let cats = (for cat in have_in_category_followings FILTER cat._from == me._id AND cat._to == comp._id return cat.name)\n",
		"let number_of_followers = followers\n",
		"let rating = 0\n",
		"let size = COUNT(for e in 1..1 INBOUND comp works_at return 1)\n",
		sort_companies("comp", filter.SortBy), "\n",
		"return MERGE(fol, {company:comp, followers, following: true, categories: cats})",
	)
	return result, params
}

func GET_FILTERED_COMPANY_FOLLOWERS(key string, filter *model.CompanyFilter) (string, map[string]interface{}) {
	compsFilter, params := create_companies_filter("me", "comp", filter)
	params["id"] = ToUserId(key)
	result := fmt.Sprint(
		"WITH users, companies", "\n",
		new_document("me", "id"), "\n",
		"for comp, fol in 1..1 INBOUND me follow", "\n",
		filter_collection(CompaniesName, "comp"), "\n",
		`FILTER comp.status == "ACTIVATED"`, "\n",
		filter_not_blocked("me", "comp"), "\n",
		compsFilter, "\n",
		count_followers("followers", "comp"), "\n",
		is_following("isFollowing", "me", "comp"), "\n",
		"let number_of_followers = followers\n",
		"let rating = 0\n",
		"let size = COUNT(for e in 1..1 INBOUND comp works_at return 1)\n",
		sort_companies("comp", filter.SortBy), "\n",
		"return MERGE(fol, {company:comp, followers, following: isFollowing})",
	)
	return result, params
}

func friends_of_friends(varName, myName string) string {
	return fmt.Sprintf(`let %s = (
for user, fr, ch in 2..2 ANY %s friendship
	filter ch.edges[0].status == 'Approved' AND ch.edges[1].status == 'Approved'
	FILTER user.status == "ACTIVATED"
	return user
)`, varName, myName)
}

func following_followers(varName, myName string) string {
	return fmt.Sprintf(`let %s = (
for user in 1..1 ANY %s follow
	return user
)`, varName, myName)
}

func followers(varName, myName string) string {
	return fmt.Sprintf(`let %s = (
for user in 1..1 OUTBOUND %s follow
	return user
)`, varName, myName)
}

func same_experience(varName, myName string) string {
	return fmt.Sprintf(`let %s = (
for comp in 1..1 OUTBOUND %s works_at
	for user in 1..1 INBOUND comp works_at
		return user
)`, varName, myName)
}

func GET_FRIEND_SUGGESTIONS(key string, pagination *model.Pagination) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"id": ToUserId(key),
	}

	return fmt.Sprint(
		"WITH users, companies", "\n",
		new_document("me", "id"), "\n",
		friends_of_friends("fofs", "me"), "\n",
		following_followers("fols", "me"), "\n",
		same_experience("sameExp", "me"), "\n",
		"for u in UNION(fofs, fols, sameExp)\n",
		filter_collection(UsersName, "u"), "\n",
		"collect user = u WITH COUNT INTO count\n",
		"sort count desc\n",
		"FILTER user != me\n",
		`FILTER user.status == "ACTIVATED"`, "\n",
		filter_not_blocked("me", "user"), "\n",
		"filter COUNT(for f in friendship filter (f._from == me._id AND f._to == user._id) or (f._to == me._id AND f._from == user._id) return 1) == 0 \n",
		is_following("following", "me", "user"), "\n",
		count_followers("followers", "user"), "\n",
		paginate(pagination), "\n",
		"return {user, following, followers}",
	), params
}

const GET_ALL_EXPERIENCE = `
WITH users
let me = DOCUMENT(@id)
for company, exp in 1..1 OUTBOUND me works_at
	return MERGE(exp, {company})
`

const GET_COMPANY_OWNERSHIP = `
WITH companies
for o in owns_company
	filter o._from == @userId AND o._to == @companyId
	return o
`

const GET_COMPANY_ADMIN = `
WITH companies
let me = DOCUMENT(@id)
for comp, admin in 1..1 OUTBOUND me admins
	filter admin._to == @companyId
	let created_by = DOCUMENT(admin.created_by_id)
	return MERGE(admin, {user: me, company: comp, created_by})
`

const GET_ALL_ADMINED_COMPANIES = `
WITH companies
let me = DOCUMENT(@id)
for comp, admin in 1..1 OUTBOUND me admins
	let created_by = DOCUMENT(admin.created_by_id)
	return MERGE(admin, {user: me, company: comp, created_by})
`

const GET_ALL_COMPANY_ADMIN = `
WITH users, companies
let company = DOCUMENT(@id)
for user, admin in 1..1 INBOUND company admins
	let created_by = DOCUMENT(admin.created_by_id)
	return MERGE(admin, {user, company, created_by})
`

func GET_SUGGESTED_COMPANIES(pagination *model.Pagination) string {
	return fmt.Sprint(`
WITH users, companies
let me = DOCUMENT(@id)
for comp in 2..3 ANY me friendship, owns_company, follow, works_at OPTIONS {uniqueVertices: 'global', bfs: true}
	filter IS_SAME_COLLECTION('companies', comp)
	let followers = LENGTH(for v IN 1..1 INBOUND comp follow return 1)`, "\n",
		paginate(pagination), "\n",
		`return {company: comp, followers}`, "\n")
}

const GET_BLOCKED_USERS_OR_COMPANIES = `
WITH users, companies
let me = DOCUMENT(@id)
for user in 1..1 OUTBOUND me blocks
  FILTER user.status == "ACTIVATED"
	let item = {id: user._key, avatar: user.avatar}
	return {
		id: user._key,
		avatar: user.avatar,
		name: IS_SAME_COLLECTION('users', user) ? CONCAT(user.first_name, " ", user.last_name) : user.name,
		is_company: IS_SAME_COLLECTION('companies', user)
	}
`

const GET_BLOCKED_USERS = `
WITH users, copmanies
let me = DOCUMENT(@id)
for user in 1..1 OUTBOUND me blocks
  FILTER user.status == "ACTIVATED"
	filter IS_SAME_COLLECTION('users', user)
	return user
`

//const GET_BLOCKED_COMPANIES = `
//let me = DOCUMENT(@id)
//for comp in 1..1 OUTBOUND me blocks
//	filter IS_SAME_COLLECTION('companies', comp)
//	return comp
//`

func GET_FILTERED_USER_FOLLOWINGS_FOR_COMPANY(key string, filter *model.UserFilter) (string, map[string]interface{}) {
	usersFilter, params := create_users_filter("company", "user", filter)
	params["id"] = ToCompanyId(key)
	result := fmt.Sprint(
		"WITH users, companies", "\n",
		new_document("company", "id"), "\n",
		"for user, fol in 1..1 OUTBOUND company follow", "\n",
		filter_collection(UsersName, "user"), "\n",
		`FILTER user.status == "ACTIVATED"`, "\n",
		filter_not_blocked("company", "user"), "\n",
		usersFilter, "\n",
		count_followers("followers", "user"), "\n",
		sort_users("user", "fol", filter.SortBy), "\n",
		"return MERGE(fol, {user, followers, following: true})",
	)
	return result, params
}

func GET_FILTERED_USER_FOLLOWERS_FOR_COMPANY(key string, filter *model.UserFilter) (string, map[string]interface{}) {
	usersFilter, params := create_users_filter("company", "user", filter)
	params["id"] = ToCompanyId(key)
	result := fmt.Sprint(
		"WITH users, companies", "\n",
		new_document("company", "id"), "\n",
		"for user,fol in 1..1 INBOUND company follow", "\n",
		filter_collection(UsersName, "user"), "\n",
		`FILTER user.status == "ACTIVATED"`, "\n",
		filter_not_blocked("company", "user"), "\n",
		usersFilter, "\n",
		count_followers("followers", "user"), "\n",
		is_following("isFollowing", "company", "user"), "\n",
		sort_users("user", "fol", filter.SortBy), "\n",
		"return MERGE(fol, {user, followers, following: isFollowing})",
	)
	return result, params
}

func GET_FILTERED_COMPANY_FOLLOWINGS_FOR_COMPANY(key string, filter *model.CompanyFilter) (string, map[string]interface{}) {
	compsFilter, params := create_companies_filter("company", "comp", filter)
	params["id"] = ToCompanyId(key)
	result := fmt.Sprint(
		"WITH users, companies", "\n",
		new_document("company", "id"), "\n",
		"for comp, fol in 1..1 OUTBOUND company follow", "\n",
		filter_collection(CompaniesName, "comp"), "\n",
		`FILTER comp.status == "ACTIVATED"`, "\n",
		filter_not_blocked("company", "comp"), "\n",
		compsFilter, "\n",
		count_followers("followers", "comp"), "\n",
		"let number_of_followers = followers\n",
		"let rating = 0\n",
		"let size = COUNT(for e in 1..1 INBOUND company works_at return 1)\n",
		sort_companies("comp", filter.SortBy), "\n",
		`LET cats = (
    FOR cat IN have_in_category_followings
    FILTER cat._from == company._id AND cat._to == comp._id
    RETURN cat.name
    )
		`, // TODO: verify that it is correct
		"return MERGE(fol, {company:comp, followers, following: true, categories: cats})",
	)
	fmt.Println("RESULT:", result)
	fmt.Println("PARAMS:", params)
	return result, params
}

func GET_FILTERED_COMPANY_FOLLOWERS_FOR_COMPANY(key string, filter *model.CompanyFilter) (string, map[string]interface{}) {
	compsFilter, params := create_companies_filter("company", "comp", filter)
	params["id"] = ToCompanyId(key)
	result := fmt.Sprint(
		"WITH companies", "\n",
		new_document("company", "id"), "\n",
		"for comp, fol in 1..1 INBOUND company follow", "\n",
		filter_collection(CompaniesName, "comp"), "\n",
		`FILTER comp.status == "ACTIVATED"`, "\n",
		filter_not_blocked("company", "comp"), "\n",
		compsFilter, "\n",
		count_followers("followers", "comp"), "\n",
		is_following("isFollowing", "company", "comp"), "\n",
		"let number_of_followers = followers\n",
		"let rating = 0\n",
		"let size = COUNT(for e in 1..1 INBOUND company works_at return 1)\n",
		sort_companies("comp", filter.SortBy), "\n",
		"return MERGE(fol, {company:comp, followers, following: isFollowing})",
	)
	return result, params
}

func GET_SUGGESTED_PEOPLE_FOR_COMPANY(key string, limit int) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"id":          ToCompanyId(key),
		"limit":       limit,
		"upper_limit": 2 * limit,
	}

	return fmt.Sprint(
		"WITH users, companies", "\n",
		new_document("company", "id"), "\n",
		following_of_followings("fofs", "company"), "\n",
		followers("fols", "company"), "\n",
		workers("workers", "company"), "\n",
		"for u in UNION(fofs, fols, workers)\n",
		filter_collection(UsersName, "u"), "\n",
		"collect user = u WITH COUNT INTO count\n",
		"sort count desc\n",
		"limit @upper_limit \n",
		is_following("following", "company", "user"), "\n",
		"filter !following \n",
		`FILTER user.status == "ACTIVATED"`, "\n",
		filter_not_blocked("company", "user"), "\n",
		"limit @limit \n",
		count_followers("followers", "user"), "\n",
		"return {user, followers}",
	), params
}

func GET_SUGGESTED_COMPANIES_FOR_COMPANY(key string, pagination *model.Pagination) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"id": ToCompanyId(key),
	}

	return fmt.Sprint(
		"WITH users, companies", "\n",
		new_document("company", "id"), "\n",
		following_of_followings("fofs", "company"), "\n",
		followers("fols", "company"), "\n",
		"for u in UNION(fofs, fols)\n",
		filter_collection(CompaniesName, "u"), "\n",
		"collect comp = u WITH COUNT INTO count\n",
		"sort count desc\n",
		is_following("following", "company", "comp"), "\n",
		"filter !following \n",
		`FILTER comp.status == "ACTIVATED"`, "\n",
		filter_not_blocked("company", "comp"), "\n",
		count_followers("followers", "comp"), "\n",
		paginate(pagination), "\n",
		"return {company: comp, followers}",
	), params
}

func BATCH_REMOVE_FROM_CATEGORY_FOR_FOLLOWINGS_QUERY_FOR_COMPANY(key string, referalIds []string, category string, all bool) (string, map[string]interface{}) {
	companyReferalIds := make([]string, len(referalIds))

	for i := range referalIds {
		companyReferalIds[i] = ToCompanyId(referalIds[i])
	}

	params := map[string]interface{}{
		"id":         ToCompanyId(key),
		"referalIds": companyReferalIds,
	}
	query := fmt.Sprint(
		"WITH users, companies", "\n",
		"for rel in have_in_category_followings\n",
		"filter rel._from == @id && rel._to in @referalIds\n",
	)
	if !all {
		params["category"] = category
		query = fmt.Sprint(query,
			"filter rel.name == @category\n",
		)
	}
	query = fmt.Sprint(query,
		"remove rel in have_in_category_followings",
	)

	fmt.Println("QUERY:", query)
	fmt.Println("PARAMS:", params)

	return query, params
}

func GET_NUMBER_OF_FOLLOWERS(id string) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"id": id,
	}
	return fmt.Sprint(
		"WITH users", "\n",
		new_document("me", "id"), "\n",
		count_followers("followers", "me"), "\n",
		"return followers",
	), params
}

func GET_REQUESTED_RECOMMENDATIONS(userKey string, pagination *model.Pagination) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"id": ToUserId(userKey),
	}
	return fmt.Sprint(
		"WITH users", "\n",
		new_document("me", "id"), "\n",
		"for user,req in 1..1 OUTBOUND me recommendation_requests\n",
		paginate(pagination), "\n",
		"return MERGE(req, {requestor: me, requested: user})",
	), params
}

func GET_RECEIVED_RECOMMENDATION_REQUESTS(userKey string, pagination *model.Pagination) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"id": ToUserId(userKey),
	}
	return fmt.Sprint(
		"WITH users", "\n",
		new_document("me", "id"), "\n",
		"for user,req in 1..1 INBOUND me recommendation_requests\n",
		"filter !req.is_ignored\n",
		paginate(pagination), "\n",
		"return MERGE(req, {requestor: user, requested: me})",
	), params
}

func GET_RECEIVED_RECOMMENDATIONS(userKey string, pagination *model.Pagination /*, isMe bool*/) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"id": ToUserId(userKey),
		// "isMe": isMe,
	}
	return fmt.Sprint(
		"WITH users", "\n",
		new_document("target_user", "id"), "\n",
		"for user,rec in 1..1 INBOUND target_user recommendations\n",
		"FILTER rec.hidden == false OR rec.hidden == null \n", // "FILTER @isMe ? true : rec.hidden == false\n",
		paginate(pagination), "\n",
		"return MERGE(rec, {recommendator: user, receiver: target_user})",
	), params
}

func GET_GIVEN_RECOMMENDATIONS(userKey string, pagination *model.Pagination) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"id": ToUserId(userKey),
	}
	return fmt.Sprint(
		"WITH users", "\n",
		new_document("target_user", "id"), "\n",
		"for user,rec in 1..1 OUTBOUND target_user recommendations\n",
		paginate(pagination), "\n",
		"return MERGE(rec, {recommendator: target_user, receiver: user})",
	), params
}

func GET_HIDDEN_RECOMMENDATIONS(userKey string, pagination *model.Pagination) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"id": ToUserId(userKey),
	}
	return fmt.Sprint(
		"WITH users", "\n",
		new_document("me", "id"), "\n",
		"for user,rec in 1..1 INBOUND me recommendations\n",
		"FILTER rec.hidden == true\n",
		paginate(pagination), "\n",
		"return MERGE(rec, {recommendator: user, receiver: me})",
	), params
}

func IS_BLOCKED(id1, id2 string) (string, map[string]interface{}) {
	// params := map[string]interface{}{
	// 	"id1": id1,
	// 	"id2": id2,
	// }
	// return fmt.Sprint(
	// 	new_document("u1", "id1"), "\n",
	// 	new_document("u2", "id2"), "\n",
	// 	is_blocked("blocked", "u1", "u2"), "\n",
	// 	"return blocked",
	// ), params

	return `
WITH users, companies
LET is_blocked = COUNT(for f in blocks
FILTER f._from == @sender_id AND f._to == @blocked_id
RETURN 1) > 0

RETURN is_blocked`, map[string]interface{}{
			"sender_id":  id1,
			"blocked_id": id2,
		}
}

// func IS_BLOCKED_COMPANY(id1, id2 string) (string, map[string]interface{}) {
// 	params := map[string]interface{}{
// 		"id1": id1,
// 		"id2": id2,
// 	}
// 	return fmt.Sprint(
// 		new_document("u1", "id1"), "\n",
// 		new_document("u2", "id2"), "\n",
// 		is_blocked("blocked", "u1", "u2"), "\n",
// 		"return blocked",
// 	), params
// }

func IS_FOLLOWING(id1, id2 string) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"id1": id1,
		"id2": id2,
	}
	return fmt.Sprint(
		"WITH users, companies", "\n",
		new_document("u1", "id1"), "\n",
		new_document("u2", "id2"), "\n",
		is_following("fol", "u1", "u2"), "\n",
		"return fol",
	), params
}

func IS_FAVOURITE(id1, id2 string) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"id1":      id1,
		"id2":      id2,
		"category": "favourites",
	}
	return fmt.Sprint(
		"WITH users, companies", "\n",
		new_document("u1", "id1"), "\n",
		new_document("u2", "id2"), "\n",
		have_in_category("fav", "u1", "u2"), "\n",
		"return fav",
	), params
}

func GET_USER_COUNTINGS(id string) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"id": id,
	}
	return fmt.Sprint(
		"WITH users, companies", "\n",
		new_document("user", "id"), "\n",
		count_friends("num_of_friends", "user"), "\n",
		count_followers("num_of_followers", "user"), "\n",
		count_followings("num_of_followings", "user"), "\n",
		count_received_recommendations("num_of_received_recommendations", "user"), "\n",
		`return {"num_of_friends": num_of_friends, "num_of_followers": num_of_followers, "num_of_followings": num_of_followings,"num_of_received_recommendations": num_of_received_recommendations}`,
	), params
}

func CLEAR_CATEGORIES(userKey, uniqueName string) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"userId":      ToUserId(userKey),
		"unique_name": uniqueName,
	}
	return "WITH users for cat in have_in_category FILTER cat._from == @userId AND cat.unique_name == @unique_name remove cat in have_in_category", params
}

func CLEAR_CATEGORIES_FOR_FOLLOWINGS(ownerId, uniqueName string) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"ownerId":     ownerId,
		"unique_name": uniqueName,
	}
	return "WITH users for cat in have_in_category_followings FILTER cat._from == @ownerId AND cat.unique_name == @unique_name remove cat in have_in_category_followings", params
}

func GET_COMPANY_COUNTINGS(id string) (string, map[string]interface{}) {
	params := map[string]interface{}{
		"id": id,
	}
	return `
	WITH users, companies
	LET comp = DOCUMENT(@id)

LET amount_of_followers = COUNT(
    FOR v, f in 1..1 INBOUND comp follow
		FILTER v.status == "ACTIVATED"
    RETURN 1
)

LET amount_of_followings = COUNT(
    FOR v, f in 1..1 OUTBOUND comp follow
		FILTER v.status == "ACTIVATED"
    RETURN 1
)

LET amount_of_employees = COUNT(
    FOR v, f in 1..1 INBOUND comp works_at
		FILTER v.status == "ACTIVATED"
    RETURN 1
)

RETURN {
    amount_of_followers,
    amount_of_followings,
    amount_of_employees
}
	`, params
}

const REMOVE_FRIENDSHIP_BY_ID = `
FOR f IN friendship
	// FILTER f.status == "Ignored"
	FILTER f._key == @friendship_id
	REMOVE f IN friendship
`

func GET_ALL_FOLLOWERS_IDS(id string) (string, map[string]interface{}) {
	params := make(map[string]interface{})
	params["id"] = id

	return fmt.Sprintf(`
		WITH users, companies
		LET me = DOCUMENT(@id)

		FOR v IN 1..1 OUTBOUND me follow
		    FILTER v.status == "ACTIVATED"
			FILTER COUNT(
			    FOR b IN blocks FILTER (b._from == v._id AND b._to == me._id) OR (b._to == v._id AND b._from == me._id) return 1
			    ) == 0
		    RETURN v._key

`), params
}
