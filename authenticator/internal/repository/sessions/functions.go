package sessions

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gitlab.lan/Rightnao-site/microservices/authenticator/internal/session"
)

// InsertSession inserts session to the mongo db though the Repository struct
func (r *Repository) InsertSession(ctx context.Context, ses session.Session) error {
	_, err := r.sessionsCollection.InsertOne(ctx, ses)
	if err != nil {
		return err
	}
	return nil
}

// DeactivateSessionByToken deactivates session and deletes it from mongo db though the Repository struct
func (r *Repository) DeactivateSessionByToken(ctx context.Context, token string) error {
	_, err := r.sessionsCollection.UpdateOne(
		ctx,
		bson.M{"token": token},
		bson.M{"$set": bson.M{
			"status":       session.StatusSignedout,
			"time_signout": time.Now()},
		})
	if err != nil {
		return err
	}
	return nil
}

// DeactivateSessions deactivates sessions except specified token
func (r *Repository) DeactivateSessions(ctx context.Context, tokenException string) (string, []string, error) {
	var objUserID primitive.ObjectID
	tokens := make([]string, 0)

	// get userID
	res := r.sessionsCollection.FindOne(
		ctx,
		bson.M{
			"token": tokenException,
		},
	)
	if res.Err() != nil {
		return "", nil, res.Err()
	}

	v := struct {
		UserID primitive.ObjectID `bson:"user"`
	}{}

	err := res.Decode(&v)
	if err != nil {
		return "", nil, err
	}

	objUserID = v.UserID

	// get all active tokens
	cursor, err := r.sessionsCollection.Find(
		ctx,
		bson.M{
			"user":   objUserID,
			"status": session.StatusActive,
			"token": bson.M{
				"$ne": tokenException,
			},
		},
	)
	if err != nil {
		return "", nil, err
	}
	defer cursor.Close(ctx)

	val := struct {
		Token string `bson:"token"`
	}{}

	for cursor.Next(ctx) {
		err = cursor.Decode(&val)
		if err != nil {
			return "", nil, err
		}
		tokens = append(tokens, val.Token)
	}

	// set as signed out
	if len(tokens) > 0 {
		_, err = r.sessionsCollection.UpdateMany(
			ctx,
			bson.M{
				"user": objUserID,
				"token": bson.M{
					"$in": tokens,
				},
			},
			bson.M{"$set": bson.M{
				"status":       session.StatusSignedout,
				"time_signout": time.Now()},
			})
		if err != nil {
			return "", nil, err
		}
	}
	return objUserID.Hex(), tokens, nil
}

// DeactivateSpecificSessionByToken deactivates the specific session choosen by the user    NEW
func (r *Repository) DeactivateSpecificSessionByToken(ctx context.Context, userID, sessionID string) (string, error) {
	sesID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return "", err
	}
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", err
	}
	_, err = r.sessionsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":  sesID,
			"user": userObjID,
		},
		bson.M{"$set": bson.M{
			"status":       session.StatusSignedout,
			"time_signout": time.Now(),
		},
		})
	if err != nil {
		return "", err
	}

	t := struct {
		Token string `bson:"token"`
	}{}

	res := r.sessionsCollection.FindOne(ctx,
		bson.M{
			"_id":  sesID,
			"user": userObjID,
		},
	)
	if err != nil {
		return "", err
	}

	err = res.Decode(&t)
	if err != nil {
		return "", err
	}

	return t.Token, nil
}

// UpdateActivityTime updates the user session activity with passed token in mongo db though the Repository struct
func (r *Repository) UpdateActivityTime(ctx context.Context, token string) error {
	_, err := r.sessionsCollection.UpdateOne(
		ctx,
		bson.M{"token": token},
		bson.M{"$set": bson.M{"last_activity_time": time.Now()}},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetListOfSessions brings back user's sessions with passed ID from mongo db though the Repository struct
func (r *Repository) GetListOfSessions(ctx context.Context, id string, first, after int32) ([]session.Session, error) {
	var f int64 = 5
	var a int64

	f = int64(first)
	a = int64(after)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("wrong_id")
	}

	sesSession := make([]session.Session, 0)

	cursor, err := r.sessionsCollection.Find(
		ctx,
		bson.M{
			"user":   objID,
			"status": session.StatusActive,
		},
		&options.FindOptions{
			Limit: &f,
			Skip:  &a,
			Sort: bson.M{
				"last_activity_time": -1,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		s := session.Session{}
		err = cursor.Decode(&s)
		if err != nil {
			return nil, err
		}
		sesSession = append(sesSession, s)
	}

	return sesSession, nil
}

// GetAmountOfSessions ...
func (r *Repository) GetAmountOfSessions(ctx context.Context, userID string) (int32, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, errors.New("wrong_id")
	}

	amount := struct {
		Amount int32 `bson:"amount"`
	}{}

	cursor, err := r.sessionsCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"user":   objID,
					"status": session.StatusActive,
				},
			},
			{
				"$group": bson.M{
					"_id": "$user",
					"amount": bson.M{
						"$sum": 1,
					},
				},
			},
		},
	)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&amount)
		if err != nil {
			return 0, err
		}
	}

	return amount.Amount, nil
}

// GetTimeOfLastActivity gives back the last time the user was active. brings back info from mongo db though the Repository struct
func (r Repository) GetTimeOfLastActivity(ctx context.Context, id string) (time.Time, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return time.Time{}, errors.New("wrong_id")
	}

	bs := struct {
		LastActivity time.Time `bson:"last_activity_time"`
	}{}

	cursor, err := r.sessionsCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"user":   objID,
				"status": session.StatusActive,
			},
		},
		{
			"$sort": bson.M{
				"last_activity_time": -1,
			},
		},
		{
			"$limit": 1,
		},
	})
	if err != nil {
		return time.Time{}, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&bs)
		if err != nil {
			return time.Time{}, err
		}
	}

	return bs.LastActivity, nil
}
