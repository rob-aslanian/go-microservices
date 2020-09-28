package rpc

// "gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"

// "google.golang.org/grpc/metadata"

// func (s *GrpcServer) Test(ctx context.Context, in *TemplateRPC.Query) (*TemplateRPC.Answer, error) {
// 	mt, _ := metadata.FromIncomingContext(ctx)
// 	answer := TemplateRPC.Answer{
// 		Message: mt["key"][0],
// 	}
// 	return &answer, nil
// }

// func (s *GrpcServer) GetAccount(context.Context, *userRPC.Empty) (*userRPC.Account, error) {
// 	a := &userRPC.Account{
// 		// Status:    UserRPC.Status_activated,
// 		Status:    userRPC.Status_NOT_ACTIVATED,
// 		Firstname: "John",
// 		Lastname:  "Doe",
// 	}
// 	return a, nil
// }
