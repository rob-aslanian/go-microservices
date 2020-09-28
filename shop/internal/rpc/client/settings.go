package clientRPC

import (
	"context"
	"errors"

	"google.golang.org/grpc/metadata"
)

// Settings for gRPC client
type Settings struct {
	Address string
}

func passThroughContext(ctx *context.Context) error {
	md, b := metadata.FromIncomingContext(*ctx)
	if b {
		*ctx = metadata.NewOutgoingContext(*ctx, md)
	} else {
		return errors.New("token is empty") // TODO: handle error
	}
	return nil
}
