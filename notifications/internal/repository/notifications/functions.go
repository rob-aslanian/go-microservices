package notificationsRepository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/notifications/internal/notification"
)

const (
	notificationsCollection = "notifications"
	settingsCollection      = "settings"

	companySettingsCollection = "company_settings"
)

// SaveNotification ...
func (r Repository) SaveNotification(ctx context.Context, not map[string]interface{}) error {
	if not == nil {
		// log.Println("ID is nil")
		return errors.New("wrong_id")
	}

	if id, isExists := not["id"]; isExists {
		delete(not, "id")

		if value, ok := id.(string); ok {
			not["_id"] = bson.ObjectIdHex(value)
		} else {
			return errors.New("wrong_id")
		}

	} else {
		return errors.New("id_not_found")
	}

	convertStringInObjectID(not, "receiver_id")
	convertStringInObjectID(not, "user_sender_id")
	convertStringInObjectID(not, "company_id")
	convertStringInObjectID(not, "service_id")
	convertStringInObjectID(not, "request_id")
	convertStringInObjectID(not, "job_id")
	convertStringInTime(not, "created_at")

	err := r.collections[notificationsCollection].Insert(not)
	if err != nil {
		return err
	}

	// log.Printf("Notification saved %s", not["_id"])

	return nil
}

// GetNotificationsSettings ...
func (r Repository) GetNotificationsSettings(ctx context.Context, userID string) (*notification.Settings, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	set := notification.Settings{}

	err := r.collections[settingsCollection].Find(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
	).One(&set)
	if err != nil {
		if err == mgo.ErrNotFound {
			return &set, nil
		}
		return nil, err
	}

	return &set, nil
}

// GetMapNotificationsSettings ...
func (r Repository) GetMapNotificationsSettings(ctx context.Context, userID string) (map[string]*bool, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	set := map[string]*bool{}

	err := r.collections[settingsCollection].Find(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
	).One(&set)
	if err != nil {
		if err == mgo.ErrNotFound {
			return set, nil
		}
		return nil, err
	}

	return set, nil
}

// ChangeNotificationsSettings ...
func (r Repository) ChangeNotificationsSettings(ctx context.Context, userID string, parameter notification.ParameterSetting, value bool) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	_, err := r.collections[settingsCollection].Upsert(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$set": bson.M{
				string(parameter): value,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetNotifications ...
func (r Repository) GetNotifications(ctx context.Context, userID string, first uint32, after uint32) ([]map[string]interface{}, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	result := make([]map[string]interface{}, 0)

	err := r.collections[notificationsCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"receiver_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$sort": bson.M{
					"created_at": -1,
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		},
	).All(&result)
	if err != nil {
		return nil, err
	}

	for i := range result {
		convertObjecIDInString(result[i])
	}

	return result, nil
}

// GetUnseenNotifications ...
func (r Repository) GetUnseenNotifications(ctx context.Context, userID string, first uint32, after uint32) ([]map[string]interface{}, error) {
	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	result := make([]map[string]interface{}, 0)

	err := r.collections[notificationsCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"receiver_id": bson.ObjectIdHex(userID),
					"seen":        false,
				},
			},
			{
				"$sort": bson.M{
					"created_at": -1,
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		},
	).All(&result)
	if err != nil {
		return nil, err
	}

	for i := range result {
		convertObjecIDInString(result[i])
	}

	return result, nil
}

// MarkAsSeen ...
func (r Repository) MarkAsSeen(ctx context.Context, userID string, ids []string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))

	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	log.Println(idsObject)

	_, err := r.collections[notificationsCollection].UpdateAll( // TODO: probably should be upsert
		bson.M{
			"_id": bson.M{
				"$in": idsObject,
			},
		},
		bson.M{
			"$set": bson.M{
				"seen": true,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveNotification ...
func (r Repository) RemoveNotification(ctx context.Context, userID string, ids []string) error {
	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))

	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	_, err := r.collections[notificationsCollection].RemoveAll(
		bson.M{
			"_id": bson.M{
				"$in": idsObject,
			},
			"receiver_id": bson.ObjectIdHex(userID),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetAmountOfNotSeenNotifications ...
func (r Repository) GetAmountOfNotSeenNotifications(ctx context.Context, userID string) (int32, error) {
	if !bson.IsObjectIdHex(userID) {
		return 0, errors.New("wrong_id")
	}

	result := struct {
		Amount int32 `bson:"amount"`
	}{}

	r.collections[notificationsCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"receiver_id": bson.ObjectIdHex(userID),
					"seen": bson.M{
						"$ne": true,
					},
				},
			},
			{
				"$group": bson.M{
					"_id": "$receiver_id",
					"amount": bson.M{
						"$sum": 1,
					},
				},
			},
		},
	).One(&result)

	return result.Amount, nil
}

// company

// GetCompanyNotificationsSettings ...
func (r Repository) GetCompanyNotificationsSettings(ctx context.Context, companyID string) (*notification.CompanySettings, error) {
	if !bson.IsObjectIdHex(companyID) {
		return nil, errors.New("wrong_id")
	}

	set := notification.CompanySettings{}

	err := r.collections[companySettingsCollection].Find(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
	).One(&set)
	if err != nil {
		if err == mgo.ErrNotFound {
			return &set, nil
		}
		return nil, err
	}

	return &set, nil
}

// ChangeCompanyNotificationsSettings ...
func (r Repository) ChangeCompanyNotificationsSettings(ctx context.Context, companyID string, parameter notification.ParameterSetting, value bool) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	_, err := r.collections[companySettingsCollection].Upsert(
		bson.M{
			"_id": bson.ObjectIdHex(companyID),
		},
		bson.M{
			"$set": bson.M{
				string(parameter): value,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetCompanyNotifications ...
func (r Repository) GetCompanyNotifications(ctx context.Context, companyID string, first uint32, after uint32) ([]map[string]interface{}, error) {
	if !bson.IsObjectIdHex(companyID) {
		return nil, errors.New("wrong_id")
	}

	result := make([]map[string]interface{}, 0)

	err := r.collections[notificationsCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"receiver_id": bson.ObjectIdHex(companyID),
				},
			},
			{
				"$sort": bson.M{
					"created_at": -1,
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		},
	).All(&result)
	if err != nil {
		return nil, err
	}

	for i := range result {
		convertObjecIDInString(result[i])
	}

	return result, nil
}

// GetAmountOfNotSeenCompanyNotifications ...
func (r Repository) GetAmountOfNotSeenCompanyNotifications(ctx context.Context, companyID string) (int32, error) {
	if !bson.IsObjectIdHex(companyID) {
		return 0, errors.New("wrong_id")
	}

	result := struct {
		Amount int32 `bson:"amount"`
	}{}

	r.collections[notificationsCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"receiver_id": bson.ObjectIdHex(companyID),
					"seen": bson.M{
						"$ne": true,
					},
				},
			},
			{
				"$group": bson.M{
					"_id": "$receiver_id",
					"amount": bson.M{
						"$sum": 1,
					},
				},
			},
		},
	).One(&result)

	return result.Amount, nil
}

// MarkAsSeenForCompany ...
func (r Repository) MarkAsSeenForCompany(ctx context.Context, companyID string, ids []string) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))

	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	_, err := r.collections[notificationsCollection].UpdateAll( // TODO: probably should be upsert
		bson.M{
			"_id": bson.M{
				"$in": idsObject,
			},
		},
		bson.M{
			"$set": bson.M{
				"seen": true,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// RemoveNotificationForCompany ...
func (r Repository) RemoveNotificationForCompany(ctx context.Context, companyID string, ids []string) error {
	if !bson.IsObjectIdHex(companyID) {
		return errors.New("wrong_id")
	}

	idsObject := make([]bson.ObjectId, 0, len(ids))

	for i := range ids {
		if bson.IsObjectIdHex(ids[i]) {
			idsObject = append(idsObject, bson.ObjectIdHex(ids[i]))
		} else {
			return errors.New("wrong_id")
		}
	}

	_, err := r.collections[notificationsCollection].RemoveAll(
		bson.M{
			"_id": bson.M{
				"$in": idsObject,
			},
			"receiver_id": bson.ObjectIdHex(companyID),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// -------

func convertStringInObjectID(not map[string]interface{}, value string) map[string]interface{} {
	if id, isExists := not[value]; isExists {
		delete(not, value)

		if v, ok := id.(string); ok {
			if v != "" {
				not[value] = bson.ObjectIdHex(v)
			}
		} else {
			// log.Printf("value %s not a string\n", value)
		}

	} else {
		// log.Printf("value %s not found\n", value)
	}

	return not
}

func convertStringInTime(not map[string]interface{}, value string) map[string]interface{} {
	if id, isExists := not[value]; isExists {
		delete(not, value)

		if v, ok := id.(string); ok {
			var err error
			not[value], err = time.Parse(time.RFC3339, v)
			if err != nil {
				log.Printf("value %s coudn't unmarshal in time\n", value)
			}
		} else {
			log.Printf("value %s not a string\n", value)
		}

	} else {
		log.Printf("value %s not found\n", value)
	}

	return not
}

func convertObjecIDInString(not map[string]interface{}) map[string]interface{} {
	// change "_id" into "id"
	if id, isExists := not["_id"]; isExists {
		not["id"] = id
		delete(not, "_id")
	}
	// convert bson.ObjectId into string
	for key := range not {
		if bs, ok := not[key].(bson.ObjectId); ok {
			not[key] = bs.Hex()
		}
	}

	return not
}
