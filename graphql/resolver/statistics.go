package resolver

import (
	"context"
	"errors"
	"fmt"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/statisticsRPC"
)

func (_ *Resolver) PersistEvent(ctx context.Context, input PersistEventRequest) (*bool, error) {
	request := statisticsRPC.PersistEventRequest{
		Event:     input.Event,
		ActorId:   input.ActorId,
		ActorType: input.ActorType,
	}
	if input.TargetId != nil && input.TargetType != nil {
		request.TargetId = *input.TargetId
		request.TargetType = *input.TargetType
	}
	statsData := make(map[string]string)
	for key, val := range input.Data {
		var ok bool
		statsData[key], ok = val.(string)
		if !ok {
			return nil, errors.New(fmt.Sprintf("Expected string value for %s, got %s", key, val))
		}
	}
	request.Data = statsData
	_, err := statistics.PersistEvent(ctx, &request)
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) IncrementCounter(ctx context.Context, input IncrementCounterRequest) (*bool, error) {
	request := statisticsRPC.IncrementCounterRequest{}
	if input.TargetId != nil {
		request.TargetId = *input.TargetId
	}
	counters := make(map[string]int32)
	for key, val := range input.Increments {
		num, ok := val.(float64)
		if !ok {
			return nil, errors.New(fmt.Sprintf("Expected int value for %s, got %s", key, val))
		}
		counters[key] = int32(num)

	}
	request.Increments = counters
	_, err := statistics.IncrementCounter(ctx, &request)
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}
