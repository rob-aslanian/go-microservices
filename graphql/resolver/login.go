package resolver

import (
	"context"
	"log"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
	hc_errors "gitlab.lan/Rightnao-site/microservices/shared/hc-errors"
)

// TODO: finish
func (_ *Resolver) Login(ctx context.Context, input LoginRequest) (*LoginResponseResolver, error) {
	a := userRPC.Credentials{
		Login:    input.Input.Login,
		Password: input.Input.Password,
		TwoFACode: func(s *string) string {
			if s != nil {
				return *s
			}
			return ""
		}(input.Input.Two_fa_code),
	}

	result, err := user.Login(ctx, &a)
	if err != nil {
		// TODO: handle error
		log.Println(err)
		e, b := hc_errors.UnwrapJsonErrorFromRPCError(err)
		if !b {
			return nil, err
		}
		return nil, e
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
			Token:            result.GetToken(),
			Is_2fa_requeried: result.GetTwoFARequired(),
			Url:              result.GetURL(),
			Avatar:           result.GetAvatar(),
			First_name:       result.GetFirstName(),
			Last_name:        result.GetLastName(),
			Gender:           result.GetGender(),
			Email:			  result.GetEmail(),
		},
	}, nil
}

func (_ *Resolver) SignOut(ctx context.Context) (*SuccessResolver, error) {
	_, err := user.SignOut(ctx, &userRPC.Empty{})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (_ *Resolver) SignOutSession(ctx context.Context, sesID SignOutSessionRequest) (*SuccessResolver, error) {
	_, err := user.SignOutSession(ctx, &userRPC.SessionID{
		ID: sesID.SessionID,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (_ *Resolver) SignOutFromAll(ctx context.Context) (*SuccessResolver, error) {
	_, err := user.SignOutFromAll(ctx, &userRPC.Empty{})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (_ *Resolver) CheckToken(ctx context.Context) (bool, error) {
	response, err := user.CheckToken(ctx, &userRPC.Empty{})
	if err != nil {
		return false, nil
	}

	return response.GetValue(), nil
}
