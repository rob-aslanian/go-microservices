package serverRPC

import (
	"context"
	"time"

	"gitlab.lan/Rightnao-site/microservices/stuff/pkg/wallet"

	"github.com/globalsign/mgo/bson"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/stuffRPC"
	additionalFeedback "gitlab.lan/Rightnao-site/microservices/stuff/pkg/feedback"
	feedback "gitlab.lan/Rightnao-site/microservices/stuff/pkg/feedback_form"
)

// AddFileToFeedBackBug ...
func (s Server) AddFileToFeedBackBug(ctx context.Context, data *stuffRPC.File) (*stuffRPC.ID, error) {

	id, err := s.service.AddFileToFeedBackBug(ctx, data.GetTargetID(), fileRPCToProfileFile(data))

	if err != nil {
		return nil, err
	}

	return &stuffRPC.ID{ID: id}, nil
}

// AddFileToFeedBackSuggestion ...
func (s Server) AddFileToFeedBackSuggestion(ctx context.Context, data *stuffRPC.File) (*stuffRPC.ID, error) {

	id, err := s.service.AddFileToFeedBackSuggestion(ctx, data.GetTargetID(), fileRPCToProfileFile(data))

	if err != nil {
		return nil, err
	}

	return &stuffRPC.ID{ID: id}, nil
}

// SaveFeedback ...
func (s Server) SaveFeedback(ctx context.Context, data *stuffRPC.FeedbackForm) (*stuffRPC.Empty, error) {
	s.service.SaveFeedback(ctx, rpcToFeedbackForm(data))
	return &stuffRPC.Empty{}, nil
}

// SubmitFeedBack ...
func (s Server) SubmitFeedBack(ctx context.Context, data *stuffRPC.FeedBack) (*stuffRPC.ID, error) {

	id, err := s.service.SubmitFeedBack(ctx, rpcToAdditionalFeedback(data))

	if err != nil {
		return nil, err
	}

	return &stuffRPC.ID{
		ID: id,
	}, nil
}

// GetAllFeedBack ...
func (s Server) GetAllFeedBack(ctx context.Context, data *stuffRPC.FeedBackRequest) (*stuffRPC.FeedBacks, error) {

	res, err := s.service.GetAllFeedBack(
		ctx,
		data.GetFirst(),
		data.GetAfter(),
	)

	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}

	return &stuffRPC.FeedBacks{
		FeedBack: feedBacksToRPC(res),
	}, nil

}

// GetAccoutWalletAmount ...
func (s Server) GetAccoutWalletAmount(ctx context.Context, data *stuffRPC.WalletAmountRequest) (*stuffRPC.WalletAmountResponse, error) {
	res, err := s.service.GetAccoutWalletAmount(
		ctx,
		data.GetUserID(),
	)

	if err != nil {
		return nil, err
	}

	return &stuffRPC.WalletAmountResponse{
		GoldCoins:     res.WalletAmount.GoldCoins,
		SilverCoins:   res.WalletAmount.SilverCoins,
		PendingAmount: res.WalletAmount.PendingAmount,
	}, nil

}

// ContactInvitationForWallet ...
func (s Server) ContactInvitationForWallet(ctx context.Context, data *stuffRPC.InvitationWalletRequest) (*stuffRPC.WalletResponse, error) {
	err := s.service.ContactInvitationForWallet(
		ctx,
		&wallet.Transaction{
			Type: wallet.TransactionTypeInvitation,
			WalletAmount: wallet.WalletAmount{
				GoldCoins:     0,
				SilverCoins:   data.GetSilverCoins(),
				PendingAmount: 0,
			},
			CoinType:        "silver",
			Status:          wallet.WalletStatusDone,
			TransactionTime: time.Now(),
		},
	)

	if err != nil {
		return nil, err
	}

	return &stuffRPC.WalletResponse{
		Amount: &stuffRPC.WalletAmountResponse{
			GoldCoins:     0,
			PendingAmount: 0,
			SilverCoins:   data.GetSilverCoins(),
		},
		Status: stuffRPC.WalletStatusEnum_DONE,
	}, nil
}

// EarnCoinsForWallet ...
func (s Server) EarnCoinsForWallet(ctx context.Context, data *stuffRPC.WalletRequest) (*stuffRPC.WalletResponse, error) {
	res, err := s.service.EarnCoinsForWallet(
		ctx,
		data.GetUserID(),
		walletRPCToWallet(data),
	)

	if err != nil {
		return nil, err
	}

	return &stuffRPC.WalletResponse{
		Amount: walletAmountToRPC(res.WalletAmount),
		Status: walletStatusToRPC(res.Status),
	}, nil
}

// GetWalletTransactions ...
func (s Server) GetWalletTransactions(ctx context.Context, data *stuffRPC.WalletTransactionRequest) (*stuffRPC.WalletTransactionResponse, error) {
	res, err := s.service.GetWalletTransactions(
		ctx,
		data.GetUserID(),
		walletTranctionCoinTypeRPCToCoinType(data.GetTransactionType()),
		data.GetFirst(),
		data.GetAfter(),
	)

	if err != nil {
		return nil, err
	}

	return &stuffRPC.WalletTransactionResponse{
		TransactionAmount: res.TransactionAmount,
		Transactions:      walletTransactionsToRPC(res.Transactions),
	}, nil
}

// GetUserByInvitedID ...
func (s Server) GetUserByInvitedID(ctx context.Context, data *stuffRPC.UserId) (*stuffRPC.WalletInvitedByCount, error) {

	count, err := s.service.GetUserByInvitedID(
		ctx,
		data.GetID(),
	)

	if err != nil {
		return nil, err
	}

	return &stuffRPC.WalletInvitedByCount{
		Count: count,
	}, nil
}

// AddGoldCoinsToWallet ...
func (s Server) AddGoldCoinsToWallet(ctx context.Context, data *stuffRPC.WalletAddGoldCoins) (*stuffRPC.Empty, error) {
	err := s.service.AddGoldCoinsToWallet(
		ctx,
		data.GetUserID(),
		data.GetCoins(),
	)

	if err != nil {
		return nil, err
	}

	return &stuffRPC.Empty{}, nil
}

// CreateWalletAccount ...
func (s Server) CreateWalletAccount(ctx context.Context, data *stuffRPC.UserId) (*stuffRPC.Empty, error) {
	err := s.service.CreateWalletAccount(
		ctx,
		data.GetID(),
	)

	if err != nil {
		return nil, err
	}

	return &stuffRPC.Empty{}, nil
}

// VoteForComingSoon ...
func (s Server) VoteForComingSoon(ctx context.Context, data *stuffRPC.VoteForComingSoonRequest) (*stuffRPC.Empty, error) {
	err := s.service.VoteForComingSoon(
		ctx,
		data.GetEmail(),
		data.GetType(),
	)

	if err != nil {
		return nil, err
	}

	return &stuffRPC.Empty{}, nil
}

func walletTranctionCoinTypeRPCToCoinType(data stuffRPC.TransactionType) wallet.TransactionCoinType {

	switch data {
	case stuffRPC.TransactionType_GOLD:
		return wallet.TransactionCoinTypeGold
	case stuffRPC.TransactionType_SILVER:
		return wallet.TransactionCoinTypeSilver
	}
	return wallet.TransactionCoinTypeAll
}

func walletTransactionsToRPC(data []*wallet.Transaction) []*stuffRPC.WalletTransaction {
	if data == nil {
		return nil
	}

	transactions := make([]*stuffRPC.WalletTransaction, 0, len(data))

	for _, t := range data {
		transactions = append(transactions, walletTransactionToRPC(t))
	}

	return transactions
}

func walletTransactionToRPC(data *wallet.Transaction) *stuffRPC.WalletTransaction {
	if data == nil {
		return nil
	}

	return &stuffRPC.WalletTransaction{
		TransactionType:   walletTypeToRPC(data.Type),
		TransactionStatus: walletStatusToRPC(data.Status),
		WalletAmount:      walletAmountToRPC(data.WalletAmount),
		TransactionAt:     data.TransactionTime.String(),
		CoinType:          walletCoinTypeToRPC(data.CoinType),
	}
}

func walletCoinTypeToRPC(data wallet.TransactionCoinType) stuffRPC.TransactionType {
	switch data {
	case wallet.TransactionCoinTypeGold:
		return stuffRPC.TransactionType_GOLD
	case wallet.TransactionCoinTypeSilver:
		return stuffRPC.TransactionType_SILVER
	}
	return stuffRPC.TransactionType_ALL
}

func walletRPCToWallet(data *stuffRPC.WalletRequest) *wallet.Wallet {
	if data == nil {
		return nil
	}

	return &wallet.Wallet{
		WalletAmount: walletAmountRPCToAmount(data.GetAmount()),
		Transaction:  walletTransactionRPCToTransaction(data),
	}
}

func walletTransactionRPCToTransaction(data *stuffRPC.WalletRequest) *wallet.Transaction {
	if data == nil {
		return nil
	}

	coinType := wallet.TransactionCoinTypeAll

	if data.GetAmount() != nil {
		if data.GetAmount().GetSilverCoins() > 0 {
			coinType = wallet.TransactionCoinTypeSilver
		} else {
			coinType = wallet.TransactionCoinTypeGold
		}
	}

	return &wallet.Transaction{
		Type:            walletTypeRPCToType(data.GetActionType()),
		WalletAmount:    walletAmountRPCToAmount(data.GetAmount()),
		CoinType:        coinType,
		Status:          wallet.WalletStatusDone,
		TransactionTime: time.Now(),
	}
}

func walletTypeRPCToType(data stuffRPC.WalletActionEnum) wallet.TransactionType {

	switch data {
	case stuffRPC.WalletActionEnum_SHARE:
		return wallet.TransactionTypeShare
	case stuffRPC.WalletActionEnum_APPLY_JOB:
		return wallet.TransactionTypeApplyJob
	case stuffRPC.WalletActionEnum_INVITATION:
		return wallet.TransactionTypeInvitation
	case stuffRPC.WalletActionEnum_USER_REGISTATION:
		return wallet.TransactionTypeUserRegistation
	case stuffRPC.WalletActionEnum_COMPANY_REGISTATION:
		return wallet.TransactionTypeCompanyRegistation
	case stuffRPC.WalletActionEnum_BECOME_CANDIDATE:
		return wallet.TransactionTypeBecomeCandidate
	case stuffRPC.WalletActionEnum_CREATE_POST:
		return wallet.TransactionTypeCreatePost
	case stuffRPC.WalletActionEnum_JOB_SHARE:
		return wallet.TransactionJobShare
	}

	return wallet.TransactionTypeUnknown

}

func walletTypeToRPC(data wallet.TransactionType) stuffRPC.WalletActionEnum {

	switch data {
	case wallet.TransactionTypeShare:
		return stuffRPC.WalletActionEnum_SHARE
	case wallet.TransactionTypeApplyJob:
		return stuffRPC.WalletActionEnum_APPLY_JOB
	case wallet.TransactionTypeInvitation:
		return stuffRPC.WalletActionEnum_INVITATION
	case wallet.TransactionTypeUserRegistation:
		return stuffRPC.WalletActionEnum_USER_REGISTATION
	case wallet.TransactionTypeCompanyRegistation:
		return stuffRPC.WalletActionEnum_COMPANY_REGISTATION
	case wallet.TransactionTypeBecomeCandidate:
		return stuffRPC.WalletActionEnum_BECOME_CANDIDATE
	case wallet.TransactionTypeCreatePost:
		return stuffRPC.WalletActionEnum_CREATE_POST
	case wallet.TransactionJobShare:
		return stuffRPC.WalletActionEnum_JOB_SHARE
	}

	return stuffRPC.WalletActionEnum_UNKNOWN_ACTION
}

func walletStatusToRPC(data wallet.Status) stuffRPC.WalletStatusEnum {

	switch data {
	case wallet.WalletStatusDone:
		return stuffRPC.WalletStatusEnum_DONE
	case wallet.WalletStatusPending:
		return stuffRPC.WalletStatusEnum_PENDING
	}

	return stuffRPC.WalletStatusEnum_REJECTED
}

func walletAmountToRPC(data wallet.WalletAmount) *stuffRPC.WalletAmountResponse {
	return &stuffRPC.WalletAmountResponse{
		GoldCoins:     data.GoldCoins,
		SilverCoins:   data.SilverCoins,
		PendingAmount: data.PendingAmount,
	}
}

func walletAmountRPCToAmount(data *stuffRPC.WalletAmountResponse) wallet.WalletAmount {
	if data == nil {
		return wallet.WalletAmount{
			GoldCoins:     0,
			SilverCoins:   0,
			PendingAmount: 0,
		}
	}

	return wallet.WalletAmount{
		GoldCoins:     data.GetGoldCoins(),
		SilverCoins:   data.GetSilverCoins(),
		PendingAmount: data.GetPendingAmount(),
	}
}

func feedBacksToRPC(data []*additionalFeedback.FeedBack) []*stuffRPC.FeedBackResponse {

	if data == nil {
		return nil
	}

	feedbacks := make([]*stuffRPC.FeedBackResponse, 0, len(data))

	for _, f := range data {
		feedbacks = append(feedbacks, feedBackToRPC(f))
	}

	return feedbacks
}

func feedBackToRPC(data *additionalFeedback.FeedBack) *stuffRPC.FeedBackResponse {
	if data == nil {
		return nil
	}

	return &stuffRPC.FeedBackResponse{
		FeedBackReaction:   feedBackReactionToRPC(data.Reactions),
		FeedBackCompliment: feedBackComplimentToRPC(data.Compliment),
		FeedBackComplaint:  feedBackComplaintToRPC(data.Complaint),
		FeedBackOther:      feedBackOtherToRPC(data.Other),
		FeedBackBugs:       feedBackBugsToRPC(data.Bugs),
		FeedBackSuggestion: feedBackSuggestionToRPC(data.Suggestion),
		UserID:             data.UserID.Hex(),
		CouldNotFind:       data.CouldNotFind,
		CreatedAt:          data.CreatedAt.String(),
	}
}

func feedBackSuggestionToRPC(data additionalFeedback.FeedBackSuggestion) *stuffRPC.FeedBackSuggestion {
	return &stuffRPC.FeedBackSuggestion{
		Idea:     data.Idea,
		Proposal: data.Proposal,
		File:     filesToRPC(data.Files),
	}
}

func feedBackBugsToRPC(data additionalFeedback.FeedBackBugs) *stuffRPC.FeedBackBugs {
	return &stuffRPC.FeedBackBugs{
		Description: data.Description,
		File:        filesToRPC(data.Files),
	}
}

func feedBackOtherToRPC(data additionalFeedback.FeedBackOther) *stuffRPC.FeedBackOther {
	return &stuffRPC.FeedBackOther{
		Description: data.Description,
		Subject:     data.Subject,
	}
}

func feedBackComplaintToRPC(data additionalFeedback.FeedBackComplaint) *stuffRPC.FeedBackComplaint {
	return &stuffRPC.FeedBackComplaint{
		ImproveExperience: data.ImproveExperience,
		MissingOrWrong:    data.MissingOrWrong,
		TellUsMore:        data.TellUsMore,
	}
}
func feedBackComplimentToRPC(data additionalFeedback.FeedBackCompliment) *stuffRPC.FeedBackCompliment {
	return &stuffRPC.FeedBackCompliment{
		FavoriteFeatures:  data.FavoriteFeatures,
		ImproveExperience: data.ImproveExperience,
		ServicesToHave:    data.ServicesToHave,
	}
}

func feedBackReactionToRPC(data additionalFeedback.FeedBackReaction) stuffRPC.FeedBackReaction {
	reaction := stuffRPC.FeedBackReaction_Very_Bad

	switch data {
	case additionalFeedback.FeedBackReactionBad:
		reaction = stuffRPC.FeedBackReaction_Bad
	case additionalFeedback.FeedBackReactionOkey:
		reaction = stuffRPC.FeedBackReaction_Okey
	case additionalFeedback.FeedBackReactionGreat:
		reaction = stuffRPC.FeedBackReaction_Great
	case additionalFeedback.FeedBackReactionVeryGood:
		reaction = stuffRPC.FeedBackReaction_Good
	case additionalFeedback.FeedBackReactionExcellent:
		reaction = stuffRPC.FeedBackReaction_Excellent

	}

	return reaction
}

func filesToRPC(data []additionalFeedback.File) []*stuffRPC.File {

	if len(data) == 0 {
		return nil
	}

	files := make([]*stuffRPC.File, 0, len(data))

	for _, f := range data {
		files = append(files, fileToRPC(f))
	}

	return files
}

func fileToRPC(data additionalFeedback.File) *stuffRPC.File {
	return &stuffRPC.File{
		ID:       data.ID.Hex(),
		Name:     data.Name,
		URL:      data.URL,
		MimeType: data.MimeType,
	}
}

// --------------------------

func rpcToFeedbackForm(data *stuffRPC.FeedbackForm) feedback.Form {
	return feedback.Form{
		Name:    data.GetName(),
		Email:   data.GetEmail(),
		Message: data.GetMessage(),
	}
}

func rpcToAdditionalFeedback(data *stuffRPC.FeedBack) additionalFeedback.FeedBack {

	if data == nil {
		return additionalFeedback.FeedBack{}
	}

	return additionalFeedback.FeedBack{
		CompanyId:    stringToID(data.CompanyId),
		Reactions:    rpcToFeedbackReaction(data.GetFeedBackReaction()),
		Compliment:   rpcToFeedbackCompliment(data.GetFeedBackCompliment()),
		Complaint:    rpcToFeedbackComplaint(data.GetFeedBackComplaint()),
		Bugs:         rpcToFeedbackBug(data.GetFeedBackBugs()),
		CouldNotFind: data.GetCouldNotFind(),
		Suggestion:   rpcToFeedbackSuggestion(data.GetFeedBackSuggestion()),
		Other:        rpcToFeedbackOther(data.GetFeedBackOther()),
	}
}

func stringToID(data string) *bson.ObjectId {
	if data == "" {
		return nil
	}

	res := bson.ObjectIdHex(data)
	return &res
}

func rpcToFeedbackReaction(data stuffRPC.FeedBackReaction) additionalFeedback.FeedBackReaction {

	switch data {
	case stuffRPC.FeedBackReaction_Bad:
		return additionalFeedback.FeedBackReactionBad
	case stuffRPC.FeedBackReaction_Okey:
		return additionalFeedback.FeedBackReactionOkey
	case stuffRPC.FeedBackReaction_Good:
		return additionalFeedback.FeedBackReactionVeryGood
	case stuffRPC.FeedBackReaction_Great:
		return additionalFeedback.FeedBackReactionGreat
	case stuffRPC.FeedBackReaction_Excellent:
		return additionalFeedback.FeedBackReactionExcellent

	}

	return additionalFeedback.FeedBackReactionVeryBad
}

func rpcToFeedbackCompliment(data *stuffRPC.FeedBackCompliment) additionalFeedback.FeedBackCompliment {

	if data == nil {
		return additionalFeedback.FeedBackCompliment{}
	}

	return additionalFeedback.FeedBackCompliment{
		FavoriteFeatures:  data.GetFavoriteFeatures(),
		ImproveExperience: data.GetImproveExperience(),
		ServicesToHave:    data.GetServicesToHave(),
	}
}

func rpcToFeedbackComplaint(data *stuffRPC.FeedBackComplaint) additionalFeedback.FeedBackComplaint {

	if data == nil {
		return additionalFeedback.FeedBackComplaint{}
	}

	return additionalFeedback.FeedBackComplaint{
		ImproveExperience: data.GetImproveExperience(),
		MissingOrWrong:    data.GetMissingOrWrong(),
		TellUsMore:        data.GetTellUsMore(),
	}
}

func rpcToFeedbackBug(data string) additionalFeedback.FeedBackBugs {

	return additionalFeedback.FeedBackBugs{
		Description: data,
	}
}

func rpcToFeedbackSuggestion(data *stuffRPC.FeedBackSuggestion) additionalFeedback.FeedBackSuggestion {

	if data == nil {
		return additionalFeedback.FeedBackSuggestion{}
	}

	return additionalFeedback.FeedBackSuggestion{
		Idea:     data.GetIdea(),
		Proposal: data.GetProposal(),
	}
}

func rpcToFeedbackOther(data *stuffRPC.FeedBackOther) additionalFeedback.FeedBackOther {

	if data == nil {
		return additionalFeedback.FeedBackOther{}
	}

	return additionalFeedback.FeedBackOther{
		Subject:     data.GetSubject(),
		Description: data.GetDescription(),
	}
}

func fileRPCToProfileFile(data *stuffRPC.File) *additionalFeedback.File {
	if data == nil {
		return nil
	}

	f := additionalFeedback.File{
		MimeType: data.GetMimeType(),
		Name:     data.GetName(),
		URL:      data.GetURL(),
	}

	_ = f.SetID(data.GetID())

	return &f
}
