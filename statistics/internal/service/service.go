package service

import (
	"context"
)

// Service ...
type Service struct {
	repository Repository
}

// NewService ...
func NewService(repo Repository) Service {
	return Service{
		repository: repo,
	}
}

// SaveUserEvent ...
func (s Service) SaveUserEvent(ctx context.Context, request map[string]interface{}, typeEvent string) error {
	err := s.repository.SaveUserEvent(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

// SaveCompanyEvent ...
func (s Service) SaveCompanyEvent(ctx context.Context, request map[string]interface{}, typeEvent string) error {
	err := s.repository.SaveUserEvent(ctx, request)
	if err != nil {
		return err
	}

	return nil
}
