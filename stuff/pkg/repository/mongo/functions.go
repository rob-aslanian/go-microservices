package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/globalsign/mgo/bson"
	comingsoon "gitlab.lan/Rightnao-site/microservices/stuff/pkg/coming_soon"
	"gitlab.lan/Rightnao-site/microservices/stuff/pkg/wallet"

	additionalFeedback "gitlab.lan/Rightnao-site/microservices/stuff/pkg/feedback"
	feedback "gitlab.lan/Rightnao-site/microservices/stuff/pkg/feedback_form"
)

const (
	// FeedBackCollection name of collection for storing feedbacks
	FeedBackCollection = "feedback"

	// AdditionalFeedBackCollection ...
	AdditionalFeedBackCollection = "additional_feedback"
	// ComingSoon ...
	ComingSoon = "coming_soon"

	// WalletCollection ...
	WalletCollection = "wallet"
)

// SaveFeedback ...
func (r Repository) SaveFeedback(ctx context.Context, fb feedback.Form) error {
	err := r.collections[FeedBackCollection].Insert(fb)
	if err != nil {
		log.Println("error while saving feedback:", err)
	}

	return nil
}

// SubmitFeedBack ...
func (r Repository) SubmitFeedBack(ctx context.Context, feedback additionalFeedback.FeedBack) error {
	err := r.collections[AdditionalFeedBackCollection].Insert(feedback)
	if err != nil {
		log.Println("error while saving feedback:", err)
	}
	return nil
}

// GetAllFeedBack ...
func (r Repository) GetAllFeedBack(ctx context.Context, first int, after int) ([]*additionalFeedback.FeedBack, error) {

	var m []*additionalFeedback.FeedBack

	err := r.collections[AdditionalFeedBackCollection].Find(nil).Sort("-created_at").Skip(after).Limit(first).All(&m)

	if err != nil {
		log.Println("error while getting feedback:", err)
	}

	log.Printf("FeedBacks %+v", m)

	return m, nil
}

// CreateWalletAccount ...
func (r Repository) CreateWalletAccount(ctx context.Context, userID string) error {

	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	selector := bson.M{"_id": bson.ObjectIdHex(userID)}
	update := bson.M{
		"wallet_amount": bson.M{
			"gold_coins":     0,
			"silver_coins":   0,
			"pending_amount": 0,
		},
		"created_at": time.Now(),
	}

	// adding wallet amount
	_, err := r.collections[WalletCollection].Upsert(selector, update)

	if err != nil {
		log.Println("error while getting wallet:", err)
	}

	return nil
}

// VoteForComingSoon ...
func (r Repository) VoteForComingSoon(ctx context.Context, data comingsoon.ComingSoon) error {

	// selector := bson.M{"_id": bson.ObjectIdHex(userID)}
	// update := bson.M{
	// 	"wallet_amount":bson.M{
	// 		"gold_coins":0,
	// 		"silver_coins":0,
	// 		"pending_amount":0,
	// 	},
	// 	"created_at":time.Now(),
	// }

	log.Printf("Voting data %+v", data)

	// adding wallet amount
	err := r.collections[ComingSoon].Insert(data)

	if err != nil {
		log.Println("error while voting for coming soon:", err)
	}

	return nil
}

// GetAccoutWalletAmount ...
func (r Repository) GetAccoutWalletAmount(ctx context.Context, userID string) (*wallet.Wallet, error) {

	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	var m *wallet.Wallet

	err := r.collections[WalletCollection].Find(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
	).One(&m)

	if m == nil {

		selector := bson.M{"_id": bson.ObjectIdHex(userID)}
		update := bson.M{
			"wallet_amount": bson.M{
				"gold_coins":     0,
				"silver_coins":   0,
				"pending_amount": 0,
			},
			"created_at": time.Now(),
		}

		// adding wallet amount
		_, err := r.collections[WalletCollection].Upsert(selector, update)

		if err != nil {
			log.Println("error while getting wallet:", err)
		}
	}

	if err != nil {
		log.Println("error while getting wallet:", err)
	}

	return m, nil
}

// ContactInvitationForWallet ...
func (r Repository) ContactInvitationForWallet(ctx context.Context, userID string, transaction *wallet.Transaction) error {

	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[WalletCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$inc": bson.M{
				"wallet_amount.silver_coins": transaction.WalletAmount.SilverCoins,
			},
			"$push": bson.M{
				"transactions": transaction,
			},
		})

	if err != nil {
		log.Println("error while getting wallet:", err)
	}

	return nil
}

// AddGoldCoinsToWallet ...
func (r Repository) AddGoldCoinsToWallet(ctx context.Context, userID string, coins int32, transaction *wallet.Transaction) error {

	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[WalletCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$inc": bson.M{
				"wallet_amount.gold_coins": coins,
			},
			"$push": bson.M{
				"transactions": transaction,
			},
		})

	if err != nil {
		log.Println("error while getting wallet:", err)
	}

	return nil
}

// EarnCoinsForWallet ...
func (r Repository) EarnCoinsForWallet(ctx context.Context, userID string, transaction *wallet.Wallet) error {

	if !bson.IsObjectIdHex(userID) {
		return errors.New("wrong_id")
	}

	err := r.collections[WalletCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(userID),
		},
		bson.M{
			"$inc": bson.M{
				"wallet_amount.gold_coins":     transaction.WalletAmount.GoldCoins,
				"wallet_amount.silver_coins":   transaction.WalletAmount.SilverCoins,
				"wallet_amount.pending_amount": transaction.WalletAmount.PendingAmount,
			},
			"$push": bson.M{
				"transactions": transaction.Transaction,
			},
		})

	if err != nil {
		log.Println("error while getting wallet:", err)
	}

	return nil

}

// GetWalletTransactions ...
func (r Repository) GetWalletTransactions(ctx context.Context, userID string, trType wallet.TransactionCoinType, first uint32, after uint32) (*wallet.Transactions, error) {

	if !bson.IsObjectIdHex(userID) {
		return nil, errors.New("wrong_id")
	}

	m := new(wallet.Transactions)
	// cointType := bson.M{}

	// switch trType {
	// case wallet.TransactionCoinTypeGold:
	// 	cointType = bson.M{
	// 		"transactions.coin_type": "gold",
	// 	}
	// case wallet.TransactionCoinTypeSilver:
	// 	cointType = bson.M{
	// 		"transactions.coin_type": "silver",
	// 	}
	// }

	match := bson.M{
		"wallet_amount": 1,
		"created_at":    1,
		"transactions":  1,
	}

	if trType != "all" {
		match = bson.M{
			"wallet_amount": 1,
			"created_at":    1,
			"transactions": bson.M{
				"$filter": bson.M{
					"input": "$transactions",
					"as":    "t",
					"cond": bson.M{
						"$eq": []interface{}{
							"$$t.coin_type", trType,
						},
					},
				},
			},
		}

	}

	res := r.collections[WalletCollection].Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"_id": bson.ObjectIdHex(userID),
				},
			},
			{
				"$project": match,
			},
			{
				"$addFields": bson.M{
					"transaction_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$transactions"},
							bson.M{"$size": "$transactions"},
							0,
						},
					},
				},
			},
			{
				"$facet": bson.M{
					"transactions": []bson.M{
						bson.M{"$unwind": "$transactions"},
						// bson.M{
						// 	"$match": cointType,
						// },
						bson.M{
							"$sort": bson.M{
								"transactions.transaction_at": -1,
							},
						},
						bson.M{"$skip": after},
						bson.M{"$limit": first},
					},
				},
			},
			{
				"$project": bson.M{
					"transactions": "$transactions.transactions",
					"transaction_amount": bson.M{
						"$cond": []interface{}{
							bson.M{"$isArray": "$transactions.transaction_amount"},
							bson.M{"$arrayElemAt": []interface{}{"$transactions.transaction_amount", 0}}, 0,
						},
					},
				},
			},
		})

	err := res.One(&m)

	if err != nil {
		log.Println("error while getting wallet transactions:", err)
	}

	log.Printf("Wallet %+v", m)

	return m, nil

}

// AddFileToFeedBackBug ...
func (r Repository) AddFileToFeedBackBug(ctx context.Context, feedBackID string, data *additionalFeedback.File) error {

	if !bson.IsObjectIdHex(feedBackID) {
		return errors.New("wrong_id")
	}

	err := r.collections[AdditionalFeedBackCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(feedBackID),
		},
		bson.M{
			"$push": bson.M{
				"bugs.files": data,
			},
		})

	if err != nil {
		log.Println("error while saving feedback:", err)
	}
	return nil
}

// AddFileToFeedBackSuggestion ...
func (r Repository) AddFileToFeedBackSuggestion(ctx context.Context, feedBackID string, data *additionalFeedback.File) error {

	if !bson.IsObjectIdHex(feedBackID) {
		return errors.New("wrong_id")
	}

	err := r.collections[AdditionalFeedBackCollection].Update(
		bson.M{
			"_id": bson.ObjectIdHex(feedBackID),
		},
		bson.M{
			"$push": bson.M{
				"suggestion.files": data,
			},
		})

	if err != nil {
		log.Println("error while saving feedback:", err)
	}
	return nil
}
