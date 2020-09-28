package clientRPC

//
// import (
// 	"context"
// 	"log"
// 	"strconv"
// 	"time"
//
// 	"gitlab.lan/Rightnao-site/microservices/grpc-proto/authRPC"
// 	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
// 	"gitlab.lan/Rightnao-site/microservices/grpc-proto/infoRPC"
// 	"gitlab.lan/Rightnao-site/microservices/grpc-proto/mailRPC"
// 	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
// 	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
// 	"gitlab.lan/Rightnao-site/microservices/user/pkg/profile"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )
//
// // GetUserID returns user id
// func (a Auth) GetUserID(ctx context.Context, token string) (string, error) {
// 	u, err := a.authClient.GetUser(ctx, &authRPC.Session{
// 		Token: token,
// 	})
//
// 	handleError(err)
//
// 	// ---------------
//
// 	return u.GetId(), nil
// }
//
// // LoginUser creates session for user. Returns token.
// func (a Auth) LoginUser(ctx context.Context, userID string) (string, error) {
// 	result, err := a.authClient.LoginUser(ctx, &authRPC.User{
// 		Id: userID,
// 	})
//
// 	handleError(err)
//
// 	return result.GetToken(), nil
// }
//
// // SignOut ...
// func (a Auth) SignOut(ctx context.Context, token string) error {
// 	_, err := a.authClient.LogoutSession(ctx, &authRPC.Session{
// 		Token: token,
// 	})
// 	handleError(err)
// 	return err
// }
//
// // GetTimeOfLastActivity ...
// func (a Auth) GetTimeOfLastActivity(ctx context.Context, id string) (time.Time, error) {
// 	tm, err := a.authClient.GetTimeOfLastActivity(ctx, &authRPC.User{
// 		Id: id,
// 	})
//
// 	handleError(err)
// 	return time.Unix(0, tm.GetTime()), err
// }
//
// // ----------------------------------------------------------------------------
//
// // SendEmail sends email
// func (m Mail) SendEmail(ctx context.Context, email string, message string) error {
// 	_, err := m.mailClient.SendMail(ctx, &mailRPC.SendMailRequest{
// 		Receiver: email,
// 		Data:     message,
// 	})
//
// 	handleError(err)
//
// 	return err
// }
//
// // ----------------------------------------------------------------------------
//
// // GetCountryIDAndCountryCode ...
// func (i Info) GetCountryIDAndCountryCode(ctx context.Context, countryCodeID int32) (countryCode string, countryID string, err error) {
// 	result, err := i.infoClient.GetCountryCodeByID(ctx, &infoRPC.CountryCode{
// 		Id: countryCodeID,
// 	})
//
// 	err = handleError(err)
// 	if err != nil {
// 		return "", "", err
// 	}
//
// 	return result.GetCountryCode(), result.GetCountry(), err
// }
//
// // GetCityInformationByID ...
// func (i Info) GetCityInformationByID(ctx context.Context, cityID int32, lang *string) (cityName, subdivision, countryID string, err error) {
// 	var l string
//
// 	if lang != nil {
// 		l = *lang
// 	}
//
// 	result, err := i.infoClient.GetCityInfoByID(ctx, &infoRPC.IDWithLang{
// 		ID:   strconv.Itoa(int(cityID)),
// 		Lang: l,
// 	})
// 	err = handleError(err)
// 	if err != nil {
// 		return
// 	}
//
// 	return result.GetTitle(), result.GetSubdivision(), result.GetCountry(), nil
// }
//
// // ----------------------------------------------------------------------------
//
// // IsFriend checks if user with userID is friend for user who called this procedure.
// func (n Network) IsFriend(ctx context.Context, userID string) (bool, error) {
// 	value, err := n.networkClient.IsFriend(ctx, &networkRPC.User{
// 		Id: userID,
// 	})
//
// 	err = handleError(err)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	return value.GetValue(), nil
// }
//
// // IsBlocked checks if user with userID is blocked for user who called this procedure.
// func (n Network) IsBlocked(ctx context.Context, userID string) (bool, error) {
// 	value, err := n.networkClient.IsBlocked(ctx, &networkRPC.User{
// 		Id: userID,
// 	})
//
// 	err = handleError(err)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	return value.GetValue(), nil
// }
//
// // IsFavourite checks if user with userID is favorite for user who called this procedure.
// func (n Network) IsFavourite(ctx context.Context, userID string) (bool, error) {
// 	value, err := n.networkClient.IsFavourite(ctx, &networkRPC.User{
// 		Id: userID,
// 	})
//
// 	err = handleError(err)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	return value.GetValue(), nil
// }
//
// // IsFollowing checks if user with userID is following by user who called this procedure.
// func (n Network) IsFollowing(ctx context.Context, userID string) (bool, error) {
// 	value, err := n.networkClient.IsFollowing(ctx, &networkRPC.User{
// 		Id: userID,
// 	})
//
// 	err = handleError(err)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	return value.GetValue(), nil
// }
//
// // GetUserCompanies returns list of companies for user.
// func (n Network) GetUserCompanies(ctx context.Context, userID string) ([]string, error) {
// 	companies, err := n.networkClient.GetUserCompanies(ctx, &networkRPC.User{
// 		Id: userID,
// 	})
// 	err = handleError(err)
// 	if err != nil {
// 		return []string{}, err
// 	}
//
// 	ids := make([]string, 0, len(companies.GetList()))
//
// 	for i := range companies.GetList() {
// 		if companies.GetList() != nil {
// 			if companies.GetList()[i].GetCompany() != nil {
// 				ids = append(ids, companies.GetList()[i].GetCompany().GetId())
// 			}
// 		}
// 	}
//
// 	return ids, nil
// }
//
// // GetReceivedRecommendationByID ...
// func (n Network) GetReceivedRecommendationByID(ctx context.Context, userID string, first int32, after int32) ([]*profile.Recommendation, error) {
// 	recommendations, err := n.networkClient.GetReceivedRecommendationById(ctx, &networkRPC.IdWithPagination{
// 		Id: userID,
// 		Pagination: &networkRPC.Pagination{
// 			After:  first,
// 			Amount: after,
// 		},
// 	})
// 	err = handleError(err)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	recs := make([]*profile.Recommendation, 0, len(recommendations.GetRecommendations()))
// 	for i := range recommendations.GetRecommendations() {
// 		rec := profile.Recommendation{
// 			ID:            recommendations.GetRecommendations()[i].GetId(),
// 			CreatedAt:     recommendations.GetRecommendations()[i].GetCreatedAt(),
// 			IsHidden:      recommendations.GetRecommendations()[i].GetHidden(),
// 			Text:          recommendations.GetRecommendations()[i].GetText(),
// 			Receiver:      &profile.Profile{},
// 			Recommendator: &profile.Profile{},
// 		}
//
// 		err := rec.Receiver.SetID(recommendations.GetRecommendations()[i].GetReceiver().GetId())
// 		if err != nil {
// 			// n.tracer.LogError(span, error)
// 			continue
// 		}
//
// 		err = rec.Recommendator.SetID(recommendations.GetRecommendations()[i].GetRecommendator().GetId())
// 		if err != nil {
// 			// n.tracer.LogError(span, error)
// 			continue
// 		}
//
// 		recs = append(recs, &rec)
// 	}
//
// 	return recs, nil
// }
//
// // GetGivenRecommendationByID ...
// func (n Network) GetGivenRecommendationByID(ctx context.Context, userID string, first int32, after int32) ([]*profile.Recommendation, error) {
// 	recommendations, err := n.networkClient.GetGivenRecommendationsById(ctx, &networkRPC.IdWithPagination{
// 		Id: userID,
// 		Pagination: &networkRPC.Pagination{
// 			After:  first,
// 			Amount: after,
// 		},
// 	})
// 	err = handleError(err)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	recs := make([]*profile.Recommendation, 0, len(recommendations.GetRecommendations()))
// 	for i := range recommendations.GetRecommendations() {
// 		rec := profile.Recommendation{
// 			ID:            recommendations.GetRecommendations()[i].GetId(),
// 			CreatedAt:     recommendations.GetRecommendations()[i].GetCreatedAt(),
// 			IsHidden:      recommendations.GetRecommendations()[i].GetHidden(),
// 			Text:          recommendations.GetRecommendations()[i].GetText(),
// 			Receiver:      &profile.Profile{},
// 			Recommendator: &profile.Profile{},
// 		}
//
// 		err := rec.Receiver.SetID(recommendations.GetRecommendations()[i].GetReceiver().GetId())
// 		if err != nil {
// 			// n.tracer.LogError(span, error)
// 			continue
// 		}
//
// 		err = rec.Recommendator.SetID(recommendations.GetRecommendations()[i].GetRecommendator().GetId())
// 		if err != nil {
// 			// n.tracer.LogError(span, error)
// 			continue
// 		}
//
// 		recs = append(recs, &rec)
// 	}
//
// 	return recs, nil
// }
//
// // GetReceivedRecommendationRequests ...
// func (n Network) GetReceivedRecommendationRequests(ctx context.Context, userID string, first int32, after int32) ([]*profile.RecommendationRequest, error) {
// 	recommendations, err := n.networkClient.GetReceivedRecommendationRequests(ctx, &networkRPC.Pagination{
// 		After:  first,
// 		Amount: after,
// 	})
// 	err = handleError(err)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	recs := make([]*profile.RecommendationRequest, 0, len(recommendations.GetRequests()))
// 	for i := range recommendations.GetRequests() {
// 		rec := profile.RecommendationRequest{
// 			ID:        recommendations.GetRequests()[i].GetId(),
// 			CreatedAt: recommendations.GetRequests()[i].GetCreatedAt(),
// 			Text:      recommendations.GetRequests()[i].GetText(),
// 			Requested: &profile.Profile{},
// 			Requestor: &profile.Profile{},
// 		}
//
// 		err := rec.Requested.SetID(recommendations.GetRequests()[i].GetRequested().GetId())
// 		if err != nil {
// 			// n.tracer.LogError(span, error)
// 			continue
// 		}
//
// 		err = rec.Requestor.SetID(recommendations.GetRequests()[i].GetRequestor().GetId())
// 		if err != nil {
// 			// n.tracer.LogError(span, error)
// 			continue
// 		}
//
// 		recs = append(recs, &rec)
// 	}
//
// 	return recs, nil
// }
//
// // GetRequestedRecommendationRequests ...
// func (n Network) GetRequestedRecommendationRequests(ctx context.Context, userID string, first int32, after int32) ([]*profile.RecommendationRequest, error) {
// 	recommendations, err := n.networkClient.GetRequestedRecommendations(ctx, &networkRPC.Pagination{
// 		After:  first,
// 		Amount: after,
// 	})
// 	err = handleError(err)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	recs := make([]*profile.RecommendationRequest, 0, len(recommendations.GetRequests()))
// 	for i := range recommendations.GetRequests() {
// 		rec := profile.RecommendationRequest{
// 			ID:        recommendations.GetRequests()[i].GetId(),
// 			CreatedAt: recommendations.GetRequests()[i].GetCreatedAt(),
// 			Text:      recommendations.GetRequests()[i].GetText(),
// 			Requested: &profile.Profile{},
// 			Requestor: &profile.Profile{},
// 		}
//
// 		err := rec.Requested.SetID(recommendations.GetRequests()[i].GetRequested().GetId())
// 		if err != nil {
// 			// n.tracer.LogError(span, error)
// 			continue
// 		}
//
// 		err = rec.Requestor.SetID(recommendations.GetRequests()[i].GetRequestor().GetId())
// 		if err != nil {
// 			// n.tracer.LogError(span, error)
// 			continue
// 		}
//
// 		recs = append(recs, &rec)
// 	}
//
// 	return recs, nil
// }
//
// // IsFriendRequestSend ...
// func (n Network) IsFriendRequestSend(ctx context.Context, targetUserID string) (bool, error) {
// 	v, err := n.networkClient.IsFriendRequestSend(ctx, &networkRPC.User{Id: targetUserID})
// 	if err != nil {
// 		log.Println("Error: IsFriendRequestSend:", err)
// 		return false, err
// 	}
// 	return v.GetValue(), nil
// }
//
// // IsFriendRequestRecieved ...
// func (n Network) IsFriendRequestRecieved(ctx context.Context, targetUserID string) (bool, string, error) {
// 	v, err := n.networkClient.IsFriendRequestRecieved(ctx, &networkRPC.User{Id: targetUserID})
// 	if err != nil {
// 		log.Println("Error: IsFriendRequestRecieved:", err)
// 		return false, "", err
// 	}
// 	return v.GetRecivied(), v.GetFriendshipID(), nil
// }
//
// // ----------------------------------------------------------------------------
//
// // GetCompanies ...
// func (c Company) GetCompanies(ctx context.Context, ids []string) (interface{}, error) {
// 	result, err := c.companyClient.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
// 		Ids: ids,
// 	})
// 	err = handleError(err)
// 	if err != nil {
// 		return []string{}, err
// 	}
//
// 	return result.GetProfiles(), nil
// }
//
// // ----------------------------------------------------------------------------
//
// func handleError(err error) error {
// 	// gRPC status
// 	if st, ok := status.FromError(err); ok {
// 		switch st.Code() {
// 		case codes.OK:
// 			return nil
// 		case codes.AlreadyExists:
// 			return usersErrors.ErrAlreadyExists
// 		case codes.NotFound:
// 			return usersErrors.ErrNotFound
// 		case codes.InvalidArgument:
// 			return usersErrors.ErrWrongArgument
//
// 			// codes.OutOfRange
// 			// codes.DataLoss
// 			// codes.Aborted
// 			// codes.FailedPrecondition
//
// 		default:
// 			// codes below are generated by gRPC
//
// 			// codes.Canceled
// 			// codes.DeadlineExceeded
// 			// codes.Internal
// 			// codes.PermissionDenied
// 			// codes.ResourceExhausted
// 			// codes.Unavailable
// 			// codes.Unimplemented
// 			// codes.Unknown
// 			// codes.Unauthenticated
// 			return usersErrors.ErrInternalError
// 		}
// 	}
//
// 	// golang error
// 	return err
// }
