package resolver

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/stuffRPC"
)

func (_ *Resolver) SaveFeedback(ctx context.Context, input SaveFeedbackRequest) (*SuccessResolver, error) {
	_, err := stuff.SaveFeedback(ctx, &stuffRPC.FeedbackForm{
		Name:    input.Feedback.Name,
		Email:   input.Feedback.Email,
		Message: input.Feedback.Message,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) SubmitFeedback(ctx context.Context, input SubmitFeedbackRequest) (*SuccessResolver, error) {

	var fb AdditionalFeedbackInput = input.Feedback

	id, err := stuff.SubmitFeedBack(ctx, &stuffRPC.FeedBack{
		CompanyId:          pointerToString(fb.Company_id),
		FeedBackReaction:   reactionToRPC(fb.Reaction),
		FeedBackCompliment: complimentToRPC(fb.Compliment),
		FeedBackComplaint:  complaintToRPC(fb.Complaint),
		FeedBackBugs:       pointerToString(fb.Bugs),
		CouldNotFind:       pointerToString(fb.Could_not_find),
		FeedBackSuggestion: suggestionToRPC(fb.Suggestion),
		FeedBackOther:      otherToRPC(fb.Other),
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      id.GetID(),
			Success: true,
		},
	}, nil
}

func (_ *Resolver) GetAllFeedBack(ctx context.Context, input GetAllFeedBackRequest) (*FeedBacksResolver, error) {

	feedbacks, err := stuff.GetAllFeedBack(ctx, &stuffRPC.FeedBackRequest{
		After: NullToString(input.Pagination.After),
		First: Nullint32ToUint32(input.Pagination.First),
	})

	if err != nil {
		return nil, err
	}

	data := feedBacksRPCToFeedBacks(feedbacks)

	userIDs := make([]string, 0, len(feedbacks.GetFeedBack()))

	for _, c := range feedbacks.GetFeedBack() {
		userIDs = append(userIDs, c.GetUserID())

	}

	userResp, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: userIDs,
	})
	if err != nil {
		return nil, err
	}

	for _, p := range feedbacks.GetFeedBack() {
		if data != nil {
			for i := range data.Feedbacks {
				// user profile
				profile := ToProfile(ctx, userResp.GetProfiles()[p.GetUserID()])
				data.Feedbacks[i].Profile = profile
			}
		}
	}

	return &FeedBacksResolver{
		R: data,
	}, nil
}

func (_ *Resolver) GetAccoutWalletAmount(ctx context.Context, input GetAccoutWalletAmountRequest) (*WalletAmountResolver, error) {

	amount, err := stuff.GetAccoutWalletAmount(ctx, &stuffRPC.WalletAmountRequest{
		UserID: input.User_id,
	})

	if err != nil {
		return nil, err
	}

	return &WalletAmountResolver{
		R: &WalletAmount{
			Gold_coins:     amount.GetGoldCoins(),
			Silver_coins:   amount.GetSilverCoins(),
			Pending_amount: amount.GetPendingAmount(),
		},
	}, nil
}

func (_ *Resolver) EarnCoinsForWallet(ctx context.Context, input EarnCoinsForWalletRequest) (*WalletResponseResolver, error) {

	res, err := stuff.EarnCoinsForWallet(ctx, &stuffRPC.WalletRequest{
		UserID:     input.User_id,
		Amount:     stuffWalletAmountToRPC(input.Wallet_input.Amount),
		ActionType: stuffWalletActionTypeToRPC(input.Wallet_input.Action_type),
	})

	if err != nil {
		return nil, err
	}

	return &WalletResponseResolver{
		R: stuffWalletRPCToWallet(res),
	}, nil

}

func (_ *Resolver) GetWalletTransactions(ctx context.Context, input GetWalletTransactionsRequest) (*WalletTransactionsResolver, error) {

	res, err := stuff.GetWalletTransactions(ctx, &stuffRPC.WalletTransactionRequest{
		UserID:          input.User_id,
		First:           Nullint32ToUint32(input.Pagination.First),
		After:           NullToString(input.Pagination.After),
		TransactionType: walletTransactionTypeToRPC(input.Type),
	})

	if err != nil {
		return nil, err
	}

	return &WalletTransactionsResolver{
		R: &WalletTransactions{
			Transition_amount: res.GetTransactionAmount(),
			Transitions:       walletTranactionsRPCToTransactions(res.GetTransactions()),
		},
	}, nil
}

func (_ *Resolver) VoteForComingSoon(ctx context.Context, input VoteForComingSoonRequest) (*SuccessResolver, error) {

	_, err := stuff.VoteForComingSoon(ctx, &stuffRPC.VoteForComingSoonRequest{
		Email: input.Email,
		Type:  input.Type,
	})

	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

func walletTransactionRPCToType(data stuffRPC.TransactionType) string {

	switch data {
	case stuffRPC.TransactionType_GOLD:
		return "gold"
	case stuffRPC.TransactionType_SILVER:
		return "silver"
	}

	return "all"
}

func walletTransactionTypeToRPC(data string) stuffRPC.TransactionType {

	switch data {
	case "gold":
		return stuffRPC.TransactionType_GOLD
	case "silver":
		return stuffRPC.TransactionType_SILVER
	}

	return stuffRPC.TransactionType_ALL
}

func walletTranactionsRPCToTransactions(data []*stuffRPC.WalletTransaction) []WalletTransaction {

	if data == nil {
		return nil
	}

	transactions := make([]WalletTransaction, 0, len(data))

	for _, t := range data {
		transactions = append(transactions, walletTranactionRPCToTransaction(t))
	}

	return transactions

}

func walletTranactionRPCToTransaction(data *stuffRPC.WalletTransaction) WalletTransaction {

	if data == nil {
		return WalletTransaction{}
	}

	return WalletTransaction{
		Status:           stuffWalletStatusRPCToWalletStatus(data.GetTransactionStatus()),
		Wallet_amount:    stuffWalletAmountRPCToWalletAmount(data.GetWalletAmount()),
		Transaction_type: stuffWalletActionTypeRPCToActionType(data.GetTransactionType()),
		Transition_at:    data.GetTransactionAt(),
		Coin_type:        walletTransactionRPCToType(data.GetCoinType()),
	}
}

func stuffWalletRPCToWallet(data *stuffRPC.WalletResponse) *WalletResponse {
	if data == nil {
		return nil
	}

	return &WalletResponse{
		Amount: stuffWalletAmountRPCToWalletAmount(data.GetAmount()),
		Status: stuffWalletStatusRPCToWalletStatus(data.GetStatus()),
	}
}

func stuffWalletAmountRPCToWalletAmount(data *stuffRPC.WalletAmountResponse) WalletAmount {
	if data == nil {
		return WalletAmount{}
	}

	return WalletAmount{
		Gold_coins:     data.GetGoldCoins(),
		Silver_coins:   data.GetSilverCoins(),
		Pending_amount: data.GetPendingAmount(),
	}
}

func stuffWalletAmountToRPC(data WalletInputAmount) *stuffRPC.WalletAmountResponse {
	return &stuffRPC.WalletAmountResponse{
		GoldCoins:     NullToInt32(data.Gold_coins),
		SilverCoins:   NullToInt32(data.Silver_coins),
		PendingAmount: NullToInt32(data.Pending_amount),
	}
}

func stuffWalletActionTypeToRPC(data string) stuffRPC.WalletActionEnum {

	switch data {
	case "share":
		return stuffRPC.WalletActionEnum_SHARE
	case "apply_job":
		return stuffRPC.WalletActionEnum_APPLY_JOB
	case "invitation":
		return stuffRPC.WalletActionEnum_INVITATION
	case "user_registration":
		return stuffRPC.WalletActionEnum_USER_REGISTATION
	case "company_registation":
		return stuffRPC.WalletActionEnum_COMPANY_REGISTATION
	case "become_candidate":
		return stuffRPC.WalletActionEnum_BECOME_CANDIDATE
	case "create_post":
		return stuffRPC.WalletActionEnum_CREATE_POST
	case "job_share":
		return stuffRPC.WalletActionEnum_JOB_SHARE
	}

	return stuffRPC.WalletActionEnum_UNKNOWN_ACTION
}

func stuffWalletActionTypeRPCToActionType(data stuffRPC.WalletActionEnum) string {

	switch data {
	case stuffRPC.WalletActionEnum_SHARE:
		return "share"
	case stuffRPC.WalletActionEnum_APPLY_JOB:
		return "apply_job"
	case stuffRPC.WalletActionEnum_INVITATION:
		return "invitation"
	case stuffRPC.WalletActionEnum_USER_REGISTATION:
		return "user_registration"
	case stuffRPC.WalletActionEnum_COMPANY_REGISTATION:
		return "company_registation"
	case stuffRPC.WalletActionEnum_BECOME_CANDIDATE:
		return "become_candidate"
	case stuffRPC.WalletActionEnum_CREATE_POST:
		return "create_post"
	case stuffRPC.WalletActionEnum_JOB_SHARE:
		return "job_share"

	}

	return "unknown"
}

func stuffWalletStatusRPCToWalletStatus(data stuffRPC.WalletStatusEnum) string {
	switch data {
	case stuffRPC.WalletStatusEnum_DONE:
		return "done"
	case stuffRPC.WalletStatusEnum_PENDING:
		return "pending"
	}

	return "rejected"
}

func feedBacksRPCToFeedBacks(data *stuffRPC.FeedBacks) *FeedBacks {
	if data == nil {
		return nil
	}

	feedbacks := make([]FeedBack, 0, len(data.GetFeedBack()))

	for _, f := range data.GetFeedBack() {
		feedbacks = append(feedbacks, feedBackRPCToFeedBack(f))
	}

	return &FeedBacks{
		Feedbacks: feedbacks,
	}
}

func feedBackRPCToFeedBack(data *stuffRPC.FeedBackResponse) FeedBack {

	if data == nil {
		return FeedBack{}
	}

	return FeedBack{
		Could_not_find: data.GetCouldNotFind(),
		Created_at:     data.GetCreatedAt(),
		Compliment:     feedBackComplimentRPCToCompliment(data.GetFeedBackCompliment()),
		Complaint:      feedBackComplaintRPCToComplaint(data.GetFeedBackComplaint()),
		Other:          feedBackOtherRPCToOther(data.GetFeedBackOther()),
		Reaction:       feedBackReactionRPCToString(data.GetFeedBackReaction()),
		Bugs:           feedBackBugRPCToBug(data.GetFeedBackBugs()),
		Suggestion:     feedBackSuggestionRPCToSuggestion(data.GetFeedBackSuggestion()),
	}
}

func feedBackSuggestionRPCToSuggestion(data *stuffRPC.FeedBackSuggestion) FeedbackSuggestion {

	if data == nil {
		return FeedbackSuggestion{}
	}

	return FeedbackSuggestion{
		Idea:     data.GetIdea(),
		Proposal: data.GetProposal(),
		Files:    stuffFilesRPCToFiles(data.GetFile()),
	}
}

func feedBackReactionRPCToString(data stuffRPC.FeedBackReaction) string {

	switch data {
	case stuffRPC.FeedBackReaction_Bad:
		return "bad"
	case stuffRPC.FeedBackReaction_Okey:
		return "okey"
	case stuffRPC.FeedBackReaction_Good:
		return "good"
	case stuffRPC.FeedBackReaction_Great:
		return "great"
	case stuffRPC.FeedBackReaction_Excellent:
		return "excellent"

	}

	return "very_bad"
}

func feedBackOtherRPCToOther(data *stuffRPC.FeedBackOther) FeedbackOther {

	if data == nil {
		return FeedbackOther{}
	}

	return FeedbackOther{
		Description: data.GetDescription(),
		Subject:     data.GetSubject(),
	}
}

func feedBackComplimentRPCToCompliment(data *stuffRPC.FeedBackCompliment) FeedbackComplimet {

	if data == nil {
		return FeedbackComplimet{}
	}

	return FeedbackComplimet{
		Favorite_features:  data.GetFavoriteFeatures(),
		Improve_experience: data.GetImproveExperience(),
		Services_to_have:   data.GetServicesToHave(),
	}
}

func feedBackComplaintRPCToComplaint(data *stuffRPC.FeedBackComplaint) FeedbackComplaint {

	if data == nil {
		return FeedbackComplaint{}
	}

	return FeedbackComplaint{
		Improve_experience: data.GetImproveExperience(),
		Missing_or_wrong:   data.GetMissingOrWrong(),
		Tell_us_more:       data.GetTellUsMore(),
	}
}

func feedBackBugRPCToBug(data *stuffRPC.FeedBackBugs) FeedBackBug {

	if data == nil {
		return FeedBackBug{}
	}

	return FeedBackBug{
		Description: data.GetDescription(),
		Files:       stuffFilesRPCToFiles(data.GetFile()),
	}
}

func stuffFilesRPCToFiles(data []*stuffRPC.File) []File {

	if data == nil {
		return nil
	}

	files := make([]File, 0, len(data))

	for _, f := range data {
		files = append(files, stuffFileRPCToFile(f))
	}

	return files
}

func stuffFileRPCToFile(data *stuffRPC.File) File {
	if data == nil {
		return File{}
	}

	return File{
		ID:        data.GetID(),
		Address:   data.GetURL(),
		Mime_type: data.GetMimeType(),
		Name:      data.GetName(),
	}
}

func pointerToString(in *string) string {
	if in == nil {
		return ""
	}
	return *in
}

func reactionToRPC(reaction string) stuffRPC.FeedBackReaction {
	switch reaction {
	case "bad":
		return stuffRPC.FeedBackReaction_Bad
	case "okey":
		return stuffRPC.FeedBackReaction_Okey
	case "good":
		return stuffRPC.FeedBackReaction_Good
	case "great":
		return stuffRPC.FeedBackReaction_Great
	case "excellent":
		return stuffRPC.FeedBackReaction_Excellent
	}

	return stuffRPC.FeedBackReaction_Very_Bad
}

func complimentToRPC(compliment FeedbackComplimetInput) *stuffRPC.FeedBackCompliment {
	return &stuffRPC.FeedBackCompliment{
		FavoriteFeatures:  pointerToString(compliment.Favorite_features),
		ImproveExperience: pointerToString(compliment.Improve_experience),
		ServicesToHave:    pointerToString(compliment.Services_to_have),
	}
}

func complaintToRPC(complaint FeedbackComplaintInput) *stuffRPC.FeedBackComplaint {
	return &stuffRPC.FeedBackComplaint{
		MissingOrWrong:    pointerToString(complaint.Missing_or_wrong),
		ImproveExperience: pointerToString(complaint.Improve_experience),
		TellUsMore:        pointerToString(complaint.Tell_us_more),
	}
}

func suggestionToRPC(suggestion FeedbackSuggestionInput) *stuffRPC.FeedBackSuggestion {
	return &stuffRPC.FeedBackSuggestion{
		Idea:     pointerToString(suggestion.Idea),
		Proposal: pointerToString(suggestion.Proposal),
	}
}

func otherToRPC(other FeedbackOtherInput) *stuffRPC.FeedBackOther {
	return &stuffRPC.FeedBackOther{
		Subject:     pointerToString(other.Subject),
		Description: pointerToString(other.Description),
	}
}
