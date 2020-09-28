package network

import "fmt"

func ToUserId(key string) string {
	return fmt.Sprint(UsersName, "/", key)
}

func ToFriendshipId(key string) string {
	return fmt.Sprint(FriendshipName, "/", key)
}

func ToCompanyId(key string) string {
	return fmt.Sprint(CompaniesName, "/", key)
}

func (r *NetworkRepo) executeBoolQuery(query string, params map[string]interface{}) (bool, error) {
	cursor, err := r.db.Query(nil, query, params)
	if err != nil {
		return false, err
	}

	var value bool
	_, err = cursor.ReadDocument(nil, &value)
	if err != nil {
		return false, nil
	}
	return value, nil
}
