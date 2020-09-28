package usersRepository

import (
	"log"

	"github.com/globalsign/mgo"
	// "github.com/mongodb/mongo-go-driver/mongo"
)

// mgo
func (r *Repository) connect(settings Settings) error {
	session, err := mgo.DialWithInfo(
		&mgo.DialInfo{
			Addrs:    settings.Addresses,
			Username: settings.User,
			Password: settings.Password,
		},
	)
	if err != nil {
		log.Panic("error while connection with db:", err)
	}

	db := session.DB(settings.Database)
	dbInvitations := session.DB("invitations")

	collection := make(map[string]*mgo.Collection /*, len(settings.Collections)*/)

	collection["users"] = db.C("users")
	collection["portfolio_views"] = db.C("portfolio_views")
	collection["portfolio_likes"] = db.C("portfolio_likes")
	collection["portfolio_comments"] = db.C("portfolio_comments")

	collection["sessions"] = db.C("sessions")
	collection["reports"] = db.C("reports")
	collection["email_invitation"] = dbInvitations.C("email_invitation")

	r.collections = collection

	return nil
}

// // mongo
// func (r *Repository) connect(settings Settings) error {
// 	// opt := options.ClientOptions{
// 	// }
// 	//
// 	// opt.SetHosts(settings.Addresses)
// 	//
// 	// opt.SetAuth(options.Credential{
// 	// 	Username: settings.User,
// 	// 	Password: settings.Password,
// 	// })
//
// 	// client, err := mongo.Connect(context.TODO(), fmt.Sprintf("mongodb://%s", settings.Addresses[0]), &opt)
// 	client, err := mongo.Connect(context.TODO(), fmt.Sprintf("mongodb://%s:%s@%s", settings.User, settings.Password, settings.Addresses[0]))
// 	if err != nil {
// 		return err
// 	}
//
// 	err = client.Ping(context.TODO(), nil)
// 	if err != nil {
// 		fmt.Printf("mongodb://%s:%s@%s", settings.User, settings.Password, settings.Addresses[0])
// 		return err
// 	}
//
// 	collection := make(map[string]*mongo.Collection)
//
// 	db := client.Database(settings.Database)
//
// 	collection["users"] = db.Collection("users")
// 	collection["sessions"] = db.Collection("sessions")
//
// 	return nil
// }
