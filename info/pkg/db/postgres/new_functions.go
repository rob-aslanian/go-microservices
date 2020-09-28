package postgres

import (
	"context"
	"errors"
	"log"
	"strconv"

	"gitlab.lan/Rightnao-site/microservices/info/pkg/models"
)

const (
	maxResultAmount = 10
)

// GetListOfAllCitiesNew Getting list of cities by first letters
func (p *PostgresDB) GetListOfAllCitiesNew(ctx context.Context, lang string, find string, first int, after int) (*[]models.City, error) {
	span := p.tracer.MakeSpan(ctx, "GetListOfAllCitiesNew")
	defer span.Finish()

	// set default language
	if lang == "" {
		lang = "en"
	}

	first = 10

	if first > maxResultAmount {
		first = maxResultAmount
	}

	if after < 0 {
		after = 0
	}

	if find == "" {
		c := make([]models.City, 0)
		return &c, nil
	}

	log.Printf("format: %q, %q, %q, %q", lang, find, first, after)

	rows, err := p.getListOfAllCitiesStmt.Query(lang, find, first, after)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cities := make([]models.City, 0)

	for rows.Next() {
		c := models.City{}

		err = rows.Scan(&c.Id, &c.Subdivision, &c.Country, new(int), &c.City) // TODO: refactor
		if err != nil {
			return nil, err
		}

		cities = append(cities, c)
	}

	log.Println("cities:", cities)

	return &cities, nil
}

// GetListOfCitiesNew Getting list of cities by first letters
func (p *PostgresDB) GetListOfCitiesNew(ctx context.Context, lang string, countryIso string, find string, first int, after int) (*[]models.City, error) {
	span := p.tracer.MakeSpan(ctx, "GetListOfCitiesNew")
	defer span.Finish()

	// set default language
	if lang == "" {
		lang = "en"
	}

	if first > maxResultAmount {
		first = maxResultAmount
	}

	if after < 0 {
		after = 0
	}

	if find == "" {
		c := make([]models.City, 0)
		return &c, nil
	}

	rows, err := p.getListOfCitiesStmt.Query(lang, find, countryIso, first, after)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cities := make([]models.City, 0)

	for rows.Next() {
		c := models.City{}

		err = rows.Scan(&c.Id, &c.Subdivision, &c.Country, new(int), &c.City) // TODO: refactor
		if err != nil {
			return nil, err
		}

		cities = append(cities, c)
	}

	return &cities, nil
}

// GetCityInfoByIDNew get city by id
func (p *PostgresDB) GetCityInfoByIDNew(ctx context.Context, lang string, cityID string) (*models.City, error) {
	span := p.tracer.MakeSpan(ctx, "GetCityInfoByIDNew")
	defer span.Finish()

	// set default language
	if lang == "" {
		lang = "en"
	}

	cityIDInt, err := strconv.Atoi(cityID)
	if err != nil {
		return nil, errors.New("wrong_id")
	}

	rows, err := p.getCityInfoByIDStmt.Query(lang, cityIDInt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	c := models.City{}
	if rows.Next() {

		err = rows.Scan(&c.Id, &c.Subdivision, &c.Country, &c.City)
		if err != nil {
			return nil, err
		}

	}

	return &c, nil
}

// GetListOfCountryCodesNew ...
func (p *PostgresDB) GetListOfCountryCodesNew(ctx context.Context) (*[]models.CountryCode, error) {
	span := p.tracer.MakeSpan(ctx, "GetListOfCountryCodesNew")
	defer span.Finish()

	rows, err := p.getListOfCountryCodesStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	countryCodes := make([]models.CountryCode, 0)
	var c models.CountryCode

	for rows.Next() {
		if err := rows.Scan(&c.Id, &c.CountryCode, &c.Country); err != nil {
			p.tracer.LogError(span, err)
			return nil, err
		}

		countryCodes = append(countryCodes, c)
	}

	return &countryCodes, nil
}

// GetCountryCodeByIDNew ...
func (p *PostgresDB) GetCountryCodeByIDNew(ctx context.Context, id uint) (*models.CountryCode, error) {
	span := p.tracer.MakeSpan(ctx, "GetCountryCodeByIDNew")
	defer span.Finish()

	rows, err := p.getCountryCodeByIDStmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var c models.CountryCode
	if rows.Next() {
		if err := rows.Scan(&c.CountryCode, &c.Country); err != nil {
			p.tracer.LogError(span, err)
			return nil, err
		}
	}

	return &c, nil
}
