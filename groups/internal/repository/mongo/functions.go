package repository

import (
	"context"
	"errors"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/groups/internal/group"
	"gitlab.lan/Rightnao-site/microservices/groups/internal/location"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SaveGroup ...
func (r Repository) SaveGroup(ctx context.Context, gr *group.Group) error {
	_, err := r.groupsCollection.InsertOne(ctx, gr)
	if err != nil {
		return err
	}

	return nil
}

// ChangeTagline ...
func (r Repository) ChangeTagline(ctx context.Context, groupID string, tagline string) error {
	objID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.groupsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
		},
		bson.M{
			"$set": bson.M{
				"tagline": tagline,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeGroupDescription ...
func (r Repository) ChangeGroupDescription(ctx context.Context, groupID string, desc, rules string, loc *location.Location) error {
	objID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.groupsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
		},
		bson.M{
			"$set": bson.M{
				"description": desc,
				"rules":       rules,
				"location":    loc,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeGroupName ...
func (r Repository) ChangeGroupName(ctx context.Context, groupID string, name string) error {
	objID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.groupsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
		},
		bson.M{
			"$set": bson.M{
				"name": name,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// IsURLBusy ...
func (r Repository) IsURLBusy(ctx context.Context, url string) (bool, error) {
	res := r.groupsCollection.FindOne(ctx, bson.M{
		"url": url,
	},
		options.FindOne().SetProjection(bson.M{
			"url": 1,
		}),
	)

	if res.Err() != nil {
		return false, res.Err()
	}

	var v interface{}
	err := res.Decode(&v)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}

		return false, err
	}

	if v != nil {
		return true, nil
	}

	return false, nil
}

// ChangeGroupPrivacyType ...
func (r Repository) ChangeGroupPrivacyType(ctx context.Context, groupID string, privacyType group.PrivacyType) error {
	objID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.groupsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
		},
		bson.M{
			"$set": bson.M{
				"privacy_type": privacyType,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ChangeGroupURL ...
func (r Repository) ChangeGroupURL(ctx context.Context, groupID, url string) error {
	objID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.groupsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objID,
		},
		bson.M{
			"$set": bson.M{
				"url": url,
			},
		},
	)

	if err != nil {
		return err
	}

	return nil
}

// AddAdmin ...
func (r Repository) AddAdmin(ctx context.Context, groupID, userID string) error {
	objGroupID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New(`wrong_id`)
	}
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.groupsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":     objGroupID,
			"members": objUserID,
		},
		bson.M{
			"$set": bson.M{
				"members.$.is_admin": true,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// AddToMembers ...
func (r Repository) AddToMembers(ctx context.Context, groupID string, m group.Member) error {
	objGroupID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.groupsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objGroupID,
		},
		bson.M{
			"$push": bson.M{
				"members": m,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// LeaveGroup ...
func (r Repository) LeaveGroup(ctx context.Context, groupID, userID string) error {
	objGroupID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New(`wrong_id`)
	}
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.groupsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objGroupID,
		},
		bson.M{
			"$pull": bson.M{
				"members": bson.M{
					"id": objUserID,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	// _, err = r.groupsCollection.UpdateOne(
	// 	ctx,
	// 	bson.M{
	// 		"_id": objGroupID,
	// 	},
	// 	bson.M{
	// 		"$pull": bson.M{
	// 			"admins": objUserID,
	// 		},
	// 	},
	// )
	// if err != nil {
	// 	return err
	// }

	return nil
}

// GetGroupByID ...
func (r Repository) GetGroupByID(ctx context.Context, groupID string) (*group.Group, error) {
	objGroupID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	res := r.groupsCollection.FindOne(ctx, bson.M{
		"_id": objGroupID,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}

	gr := new(group.Group)
	err = res.Decode(gr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return gr, nil
}

// GetGroupByURL ...
func (r Repository) GetGroupByURL(ctx context.Context, url string) (*group.Group, error) {
	res := r.groupsCollection.FindOne(ctx, bson.M{
		"url": url,
	},
		&options.FindOneOptions{
			Projection: bson.M{
				"members":    0,
				"admins":     0,
				"invitaions": 0,
			},
		},
	)
	if res.Err() != nil {
		return nil, res.Err()
	}

	gr := new(group.Group)
	err := res.Decode(gr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return gr, nil
}

// IsMember ...
func (r Repository) IsMember(ctx context.Context, groupID, userID string) (bool, error) {
	objGroupID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return false, errors.New(`wrong_id`)
	}
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, errors.New(`wrong_id`)
	}

	res := r.groupsCollection.FindOne(ctx, bson.M{
		"_id":        objGroupID,
		"members.id": objUserID,
	},
		&options.FindOneOptions{
			Projection: bson.M{
				"_id": 1,
			},
		},
	)
	if res.Err() != nil {
		return false, res.Err()
	}

	gr := new(group.Group)
	err = res.Decode(gr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}

		return false, err
	}

	if gr != nil {
		return true, nil
	}

	return false, nil
}

// AddInvitations ...
func (r Repository) AddInvitations(ctx context.Context, groupID string, users []group.InvitedMember) error {
	objGroupID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.groupsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objGroupID,
		},
		bson.M{
			"$push": bson.M{
				"invitaions": bson.M{"$each": users},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// IsInvited ...
func (r Repository) IsInvited(ctx context.Context, groupID, userID string) (bool, error) {
	objGroupID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return false, errors.New(`wrong_id`)
	}
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, errors.New(`wrong_id`)
	}

	res := r.groupsCollection.FindOne(ctx, bson.M{
		"_id":           objGroupID,
		"invitaions.id": objUserID,
	},
		&options.FindOneOptions{
			Projection: bson.M{
				"_id": 1,
			},
		},
	)
	if res.Err() != nil {
		return false, res.Err()
	}

	var gr interface{}
	err = res.Decode(&gr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}

		return false, err
	}

	if gr != nil {
		return true, nil
	}

	return false, nil
}

// RemoveInvitations ...
func (r Repository) RemoveInvitations(ctx context.Context, groupID string, userIDs []string) error {
	objGroupID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objUserIDs := make([]primitive.ObjectID, 0, len(userIDs))
	for i := range userIDs {
		id, err := primitive.ObjectIDFromHex(userIDs[i])
		if err != nil {
			return err
		}
		objUserIDs = append(objUserIDs, id)
	}

	_, err = r.groupsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objGroupID,
		},
		bson.M{
			"$pull": bson.M{
				"invitaions": bson.M{
					"id": bson.M{"$in": objUserIDs},
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// AddJoinRequest ...
func (r Repository) AddJoinRequest(ctx context.Context, groupID string, user group.Member) error {
	objGroupID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.groupsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objGroupID,
		},
		bson.M{
			"$push": bson.M{
				"invitation_requests": user,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// IsRequestSend ...
func (r Repository) IsRequestSend(ctx context.Context, groupID, userID string) (bool, error) {
	objGroupID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return false, errors.New(`wrong_id`)
	}
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, errors.New(`wrong_id`)
	}

	res := r.groupsCollection.FindOne(ctx, bson.M{
		"_id":                    objGroupID,
		"invitation_requests.id": objUserID,
	},
		&options.FindOneOptions{
			Projection: bson.M{
				"_id": 1,
			},
		},
	)
	if res.Err() != nil {
		return false, res.Err()
	}

	var gr interface{}
	err = res.Decode(&gr)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}

		return false, err
	}

	if gr != nil {
		return true, nil
	}

	return false, nil
}

// RemoveInvitationRequest ...
func (r Repository) RemoveInvitationRequest(ctx context.Context, groupID, userID string) error {
	objGroupID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New(`wrong_id`)
	}

	_, err = r.groupsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objGroupID,
		},
		bson.M{
			"$pull": bson.M{
				"invitation_requests": bson.M{
					"id": objUserID,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetMembers ...
func (r Repository) GetMembers(ctx context.Context, groupID string, first, after uint32) ([]*group.Member, error) {
	objGroupID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return nil, errors.New(`wrong_id`)
	}

	cursor, err := r.groupsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"_id": objGroupID,
				},
			},
			{
				"$unwind": bson.M{
					"path": "$members",
				},
			},
			{
				"$replaceRoot": bson.M{
					"newRoot": "$members",
				},
			},
			{
				"$skip": after,
			},
			{
				"$limit": first,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	members := make([]*group.Member, 0)

	for cursor.Next(ctx) {
		m := group.Member{}
		err = cursor.Decode(&m)
		if err != nil {
			return nil, err
		}
		members = append(members, &m)
	}

	return members, nil
}
