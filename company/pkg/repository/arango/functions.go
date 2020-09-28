package arangorepo

import (
	"context"
	"time"
)

type Industry struct {
	Main string   `json:"main"`
	Sub  []string `json:"sub"`
}

type Company struct {
	ID        string    `json:"_key"`
	Industry  Industry  `json:"industry"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

func (n NetworkRepo) SaveCompany(ctx context.Context, c *Company) error {
	_, err := n.companies.CreateDocument(ctx, c)
	if err != nil {
		return err
	}

	return nil
}
