package serverRPC

import (
	"context"
	wallet "gitlab.lan/Rightnao-site/microservices/stuff/pkg/wallet"
	additionalFeedback "gitlab.lan/Rightnao-site/microservices/stuff/pkg/feedback"
	feedback "gitlab.lan/Rightnao-site/microservices/stuff/pkg/feedback_form"
)

// Service define functions inside Service
type Service interface {
	AddFileToFeedBackBug(ctx context.Context , feedBackID string , data *additionalFeedback.File) (string , error) 
	AddFileToFeedBackSuggestion(ctx context.Context  , feedBackID string ,data *additionalFeedback.File) (string , error) 
	GetAllFeedBack(ctx context.Context  , first uint32, after string) ([]*additionalFeedback.FeedBack , error) 
	CreateWalletAccount(ctx context.Context , userID string) (error)
	GetAccoutWalletAmount(ctx context.Context , userID string) (*wallet.Wallet , error)
	ContactInvitationForWallet(ctx context.Context , transaction *wallet.Transaction) error
	EarnCoinsForWallet(ctx context.Context , userID string ,  transaction *wallet.Wallet) (*wallet.Transaction , error)
	GetWalletTransactions(ctx context.Context, userID string , trType wallet.TransactionCoinType ,  first uint32, after string  ) (*wallet.Transactions, error) 
	GetUserByInvitedID(ctx context.Context , userID string) (int32 , error)
	AddGoldCoinsToWallet(ctx context.Context , userID string , coins int32) error
	VoteForComingSoon(ctx context.Context , email string , csType string) error
		
	SaveFeedback(context.Context, feedback.Form) error
	SubmitFeedBack(context.Context, additionalFeedback.FeedBack) (string , error)
}
