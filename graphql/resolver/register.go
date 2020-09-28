package resolver

import (
	"context"
	"errors"
	"log"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
	"google.golang.org/grpc/status"
)

func (_ *Resolver) CheckUsername(ctx context.Context, input CheckUsernameRequest) (bool, error) {
	response, err := user.IsUsernameBusy(ctx, &userRPC.Username{
		Username: input.Username,
	})
	if err != nil {
		return false, nil
	}

	return response.GetValue(), nil
}

func (_ *Resolver) Register(ctx context.Context, input RegisterRequest) (*LoginResponseResolver, error) {
	result, err := user.RegisterUser(ctx, &userRPC.RegisterRequest{
		FirstName: input.Input.Firstname,
		LastName:  input.Input.Lastname,
		Email:     input.Input.Email,
		Username:  input.Input.Username,
		// CountryPrefixCode: input.Input.Country_code,
		// PhoneNumber:       input.Input.Number,
		Password:   input.Input.Password,
		Birthday:   input.Input.Birthday,
		CountryId:  NullToString(input.Input.Country),
		LanguageId: NullToString(input.Input.Language),
		Gender: func(g string) userRPC.GenderValue {
			if g == "male" {
				return userRPC.GenderValue_MALE
			}
			return userRPC.GenderValue_FEMALE
		}(input.Input.Gender),
		InvitedBy: NullToString(input.Input.Invited_by),
	})
	if err != nil {
		// // TODO: handle error
		// e, b := hc_errors.UnwrapJsonErrorFromRPCError(err)
		// if !b {
		// 	return nil, err
		// }
		// return nil, e

		errStatus, _ := status.FromError(err)
		log.Println(errStatus.Code())
		return nil, errors.New(errStatus.Message())
	}

	return &LoginResponseResolver{
		R: &LoginResponse{
			ID: result.GetUserId(),
			Status: func(s userRPC.Status) string {
				switch s {
				case userRPC.Status_ACTIVATED:
					return "activated"
				case userRPC.Status_DISABLED:
					return "disabled"
				case userRPC.Status_BLOCKED:
					return "blocked"
				}
				return "not_activated"
			}(result.GetStatus()),
			Token: result.GetToken(),
			Url:   result.GetURL(),
		},
	}, nil
}

func (_ *Resolver) ChooseActivationMethod(ctx context.Context, input ChooseActivationMethodRequest) (*LoginResponseResolver, error) {
	// _, err := user.SendActivation(ctx, &userRPC.SendActivationRequest{
	// 	UserId:  input.ID,
	// 	ByEmail: input.Methods.By_email,
	// 	ByPhone: input.Methods.By_SMS,
	// })
	// if err != nil {
	// 	// TODO: handle error
	// 	log.Println(err)
	// 	e, b := hc_errors.UnwrapJsonErrorFromRPCError(err)
	// 	if !b {
	// 		return nil, err
	// 	}
	// 	return nil, e
	// }

	return nil, nil
	// return &LoginResponseResolver{
	// 	R: &LoginResponse{
	// 		ID: result.GetUserId(),
	// 		Status: func(s userRPC.Status) string {
	// 			switch s {
	// 			case userRPC.Status_ACTIVATED:
	// 				return "activated"
	// 			case userRPC.Status_DISABLED:
	// 				return "disabled"
	// 			case userRPC.Status_BLOCKED:
	// 				return "blocked"
	// 			}
	// 			return "not_activated"
	// 		}(result.GetStatus()),
	// 		Token: result.GetToken(),
	// 	},
	// }, nil
}

func (_ *Resolver) ActivateUser(ctx context.Context, input ActivateUserRequest) (*ActivateResponseResolver, error) {
	res, err := user.ActivateUser(ctx, &userRPC.ActivateUserRequest{
		Code:   input.Code,
		UserID: input.User_id,
	})
	if err != nil {
		return nil, err
	}

	return &ActivateResponseResolver{
		R: &ActivateResponse{
			ID:         res.GetUserId(),
			Avatar:     res.GetAvatar(),
			First_name: res.GetFirstName(),
			Last_name:  res.GetLastName(),
			Token:      res.GetToken(),
			Url:        res.GetURL(),
		},
	}, nil
}
