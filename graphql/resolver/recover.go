package resolver

import (
	"context"
	"log"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
	hc_errors "gitlab.lan/Rightnao-site/microservices/shared/hc-errors"
)

func (_ *Resolver) SendRecoverRequest(ctx context.Context, input SendRecoverRequestRequest) (*SuccessResolver, error) {
	_, err := user.SendRecover(ctx, &userRPC.SendRecoverRequest{
		Login:         input.Login,
		ByEmail:       input.Methods.By_email,
		ByPhone:       input.Methods.By_SMS,
		ResetPassword: input.Methods.Reset_password,
		SendUsername:  input.Methods.Send_username,
	})

	if err != nil {
		// TODO: handle error
		log.Println(err)
		e, b := hc_errors.UnwrapJsonErrorFromRPCError(err)
		if !b {
			return nil, err
		}
		return nil, e
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}

func (_ *Resolver) RecoverPassword(ctx context.Context, input RecoverPasswordRequest) (*SuccessResolver, error) {
	_, err := user.RecoverPassword(ctx, &userRPC.RecoverPasswordRequest{
		UserId:   input.RecoveryRequest.ID,
		Password: input.RecoveryRequest.Password,
		Code:     input.RecoveryRequest.Code,
	},
	)
	if err != nil {
		// TODO: handle error
		log.Println(err)
		e, b := hc_errors.UnwrapJsonErrorFromRPCError(err)
		if !b {
			return nil, err
		}
		return nil, e
	}

	return &SuccessResolver{R: &Success{Success: true}}, nil
}
