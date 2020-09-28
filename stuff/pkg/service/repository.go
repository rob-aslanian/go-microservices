package service

import (
	"gitlab.lan/Rightnao-site/microservices/stuff/pkg/wallet"
	"context"

	additionalFeedback "gitlab.lan/Rightnao-site/microservices/stuff/pkg/feedback"
	feedback "gitlab.lan/Rightnao-site/microservices/stuff/pkg/feedback_form"
	comingsoon "gitlab.lan/Rightnao-site/microservices/stuff/pkg/coming_soon"

)

// Repository ...
type Repository interface {
	AddFileToFeedBackBug(ctx context.Context , feedBackID string ,  data *additionalFeedback.File) error
	AddFileToFeedBackSuggestion(ctx context.Context , feedBackID string ,  data *additionalFeedback.File) error
	GetAllFeedBack(ctx context.Context, first int, after int,) ([]*additionalFeedback.FeedBack , error)
	CreateWalletAccount(ctx context.Context , userID string) (error)
	GetAccoutWalletAmount(ctx context.Context , userID string) (*wallet.Wallet , error)
	ContactInvitationForWallet(ctx context.Context, userID string ,  transaction *wallet.Transaction) error
	EarnCoinsForWallet(ctx context.Context , userID string ,  transaction *wallet.Wallet) error
	GetWalletTransactions(ctx context.Context, userID string , trType wallet.TransactionCoinType , first uint32, after uint32 ) (*wallet.Transactions, error)
	AddGoldCoinsToWallet(ctx context.Context , userID string , coins int32 , transaction *wallet.Transaction) error
	VoteForComingSoon(ctx context.Context , data comingsoon.ComingSoon) error
 

	SaveFeedback(context.Context, feedback.Form) error
	SubmitFeedBack(context.Context, additionalFeedback.FeedBack) error
}
