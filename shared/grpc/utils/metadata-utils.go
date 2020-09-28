package utils

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func ExtractMetadata(ctx context.Context, keys ...string) map[string]string {
	res := make(map[string]string, len(keys))
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for _, key := range keys {
			arr := md.Get(key)
			if len(arr) > 0 {
				res[key] = arr[0]
			}
		}
	}
	return res
}

func ExtractToken(ctx context.Context) (string, bool) {
	meta := ExtractMetadata(ctx, "token")
	token, ok := meta["token"]
	return token, ok
}

func ExtractUserId(ctx context.Context) (string, bool) {
	meta := ExtractMetadata(ctx, "user-id")
	id, ok := meta["user-id"]
	return id, ok
}

func AddToIncomingMetadata(ctx context.Context, key, val string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		md.Set(key, val)
	}
}

func ToOutContext(ctx context.Context) context.Context {
	md, b := metadata.FromIncomingContext(ctx)
	var outCtx context.Context
	if b {
		outCtx = metadata.NewOutgoingContext(ctx, md)
	} else {
		outCtx = ctx
	}
	return outCtx
}
