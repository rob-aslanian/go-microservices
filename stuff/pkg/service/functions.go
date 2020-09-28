package service

import (
	"gitlab.lan/Rightnao-site/microservices/stuff/pkg/coming_soon"
	"gitlab.lan/Rightnao-site/microservices/stuff/pkg/wallet"
	"strconv"
	"context"
	"errors"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	additionalFeedback "gitlab.lan/Rightnao-site/microservices/stuff/pkg/feedback"
	feedback "gitlab.lan/Rightnao-site/microservices/stuff/pkg/feedback_form"
	"google.golang.org/grpc/metadata"
)

// SaveFeedback ...
func (s Service) SaveFeedback(ctx context.Context, fb feedback.Form) error {
	err := fb.Validate()
	if err != nil {
		return err
	}

	// Extract token
	_ = s.passThroughContext(&ctx)
	token, _ := s.extractToken(ctx)
	UserID, err := s.authRPC.GetUser(ctx, token)
	if err == nil {
		fb.SetUserID(UserID)
	}

	fb.GenerateID()
	fb.CreatedAt = time.Now()

	err = s.repository.SaveFeedback(ctx, fb)
	if err != nil {
		//
	}

	return nil
}


// SubmitFeedBack ... 
func (s Service) SubmitFeedBack(ctx context.Context, feedback additionalFeedback.FeedBack) (string  , error) {
	// Extract token
	_ = s.passThroughContext(&ctx)
	token, _ := s.extractToken(ctx)
	UserID, err := s.authRPC.GetUser(ctx, token)
	if err == nil {
		feedback.SetUserID(UserID)
	}

	feedback.GenerateID()
	feedback.CreatedAt = time.Now()

	err = s.repository.SubmitFeedBack(ctx, feedback)

	return feedback.GetID() , nil
}

// AddFileToFeedBackBug ... 
func (s Service) AddFileToFeedBackBug(ctx context.Context , feedBackID string ,   data *additionalFeedback.File) (string , error) {

	id := data.GenerateID()

	err := s.repository.AddFileToFeedBackBug(ctx, feedBackID ,  data)
	if err != nil {
		return "" , nil
	}

	return id , nil
}

// AddFileToFeedBackSuggestion ... 
func (s Service) AddFileToFeedBackSuggestion(ctx context.Context , feedBackID string ,   data *additionalFeedback.File) (string , error) {

	id := data.GenerateID()

	err := s.repository.AddFileToFeedBackSuggestion(ctx, feedBackID ,  data)
	if err != nil {
		return "" , nil
	}

	return id , nil
}

// GetAllFeedBack ... 
func (s Service) GetAllFeedBack(ctx context.Context  , first uint32, after string) ([]*additionalFeedback.FeedBack , error)  {

	afterNumber, err := strconv.Atoi(after)
	if err != nil {
		return nil, errors.New("bad_after_value")
	}
	if afterNumber < 0 {
		return nil, errors.New("bad_after_value")
	}


	res , err := s.repository.GetAllFeedBack(ctx, int(first) , int(afterNumber))

	if err != nil {
		return nil , nil
	}

	return res , nil
}

// GetAccoutWalletAmount ... 
func (s Service) GetAccoutWalletAmount(ctx context.Context , userID string) (*wallet.Wallet , error) {

	res  , err := s.repository.GetAccoutWalletAmount(ctx , userID)

	if err != nil {
		return nil , err
	}

	if res == nil {
		res = &wallet.Wallet{
			WalletAmount:wallet.WalletAmount{
				GoldCoins:0,
				SilverCoins:0,
				PendingAmount:0,
			},
		}
	}

	return res, nil


}


// AddGoldCoinsToWallet ... 
func (s Service) AddGoldCoinsToWallet(ctx context.Context , userID string , coins int32) error {

	
	if coins > 0 {
		transaction := &wallet.Transaction{
				WalletAmount:wallet.WalletAmount{
					GoldCoins:coins,
					SilverCoins:0,
					PendingAmount:0,
				},
				Status:"done",
				Type:"user_registration",
				CoinType:"gold",
				TransactionTime:time.Now(),
			}
		err := s.repository.AddGoldCoinsToWallet(ctx , userID , coins , transaction)

		if err != nil {
			return err
		} 

		return nil 
	}

	return nil
}

// CreateWalletAccount ... 
func (s Service) CreateWalletAccount(ctx context.Context , userID string) (error) {

	err := s.repository.CreateWalletAccount(ctx , userID)

	if err != nil {
		return err
	} 


	return nil
}


// VoteForComingSoon ... 
func (s Service) VoteForComingSoon(ctx context.Context , email string , csType string) error {

	data := comingsoon.ComingSoon{
		Email:email,
		Type:csType,
	}

	// Extract token
	_ = s.passThroughContext(&ctx)
	token, _ := s.extractToken(ctx)
	userID, err := s.authRPC.GetUser(ctx, token)

	if err == nil {
		data.SetUserID(userID)
	}
	data.CreatedAt = time.Now()

	err = s.repository.VoteForComingSoon(ctx , data)

	if err != nil {
		return  err
	}

	return nil
}

// ContactInvitationForWallet ... 
func (s Service) ContactInvitationForWallet(ctx context.Context , transaction *wallet.Transaction) error {

	// Extract token
	_ = s.passThroughContext(&ctx)
	token, _ := s.extractToken(ctx)
	userID, err := s.authRPC.GetUser(ctx, token)

	err = s.repository.ContactInvitationForWallet(ctx , userID , transaction)

	if err != nil {
		return  err
	}

	return nil
}

// EarnCoinsForWallet ... 
func (s Service) EarnCoinsForWallet(ctx context.Context , userID string ,  transaction *wallet.Wallet) (*wallet.Transaction , error) {

	err := s.repository.EarnCoinsForWallet(ctx , userID , transaction)

	if err != nil {
		return nil , err
	}


	return &wallet.Transaction{
		Status:transaction.Transaction.Status,
		TransactionTime:transaction.Transaction.TransactionTime,
		Type:transaction.Transaction.Type,
		WalletAmount:transaction.Transaction.WalletAmount,
	} , nil 
}

// GetUserByInvitedID ... 
func (s Service) GetUserByInvitedID(ctx context.Context , userID string) (int32 , error) {

	count , err := s.userRPC.GetUserByInvitedID(ctx , userID)

	if err != nil {
		return 0 , err
	}

	return count ,  nil
}


// GetWalletTransactions ... 
func (s Service) GetWalletTransactions(ctx context.Context, userID string , trType wallet.TransactionCoinType ,  first uint32, after string ) (*wallet.Transactions, error)  {

	afterNumber, err := strconv.Atoi(after)
	if err != nil {
		return nil, errors.New("bad_after_value")
	}
	if afterNumber < 0 {
		return nil, errors.New("bad_after_value")
	}


	res , err := s.repository.GetWalletTransactions(ctx, userID , trType , first , uint32(afterNumber))

	if err != nil {
		return nil , nil
	}

	return res , nil
}



// --------------------------

func (s *Service) passThroughContext(ctx *context.Context) error {
	span := opentracing.SpanFromContext(*ctx)
	span = span.Tracer().StartSpan("passThroughContext", opentracing.ChildOf(span.Context()))
	defer span.Finish()

	md, b := metadata.FromIncomingContext(*ctx)
	if b {
		*ctx = metadata.NewOutgoingContext(*ctx, md)
	} else {
		return errors.New("token is empty") // TODO: handle error
	}
	return nil
}

func (s *Service) extractToken(ctx context.Context) (string, error) {
	span := opentracing.SpanFromContext(ctx)
	span = span.Tracer().StartSpan("extractToken", opentracing.ChildOf(span.Context()))
	defer span.Finish()

	res := make(map[string]string, 1)
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		arr := md.Get("token")
		if len(arr) > 0 {
			res["token"] = arr[0]
		}
	}

	token, ok := res["token"]
	if !ok {
		return "", errors.New("token is empty")
	}
	return token, nil
}
