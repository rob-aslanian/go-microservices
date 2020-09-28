package statisticsrepo

import (
	"context"
)

// SaveUserEvent ...
func (r *Repository) SaveUserEvent(ctx context.Context, event map[string]interface{}) error {
	_, err := r.usersStatisticsCollection.InsertOne(
		ctx,
		event,
	)
	if err != nil {
		return err
	}

	return nil
}

// SaveCompanyEvent ...
func (r *Repository) SaveCompanyEvent(ctx context.Context, event map[string]interface{}) error {
	_, err := r.companyStatisticsCollection.InsertOne(
		ctx,
		event,
	)
	if err != nil {
		return err
	}

	return nil
}
