package postgres

import (
	"context"
	"log"
	"strings"

	"gitlab.lan/Rightnao-site/microservices/info/pkg/models"
)

// TODO: finish
func (p *PostgresDB) GetListOfCountries(ctx context.Context, lang string) (*[]string, error) {
	span := p.tracer.MakeSpan(ctx, "")
	defer span.Finish()

	rows, err := p.queryContext(ctx, getQueryGetListOfCountries(lang))
	if err != nil {
		p.tracer.LogError(span, err)
		return nil, err
	}
	defer rows.Close()
	countries := make([]string, 0) // TODO: count amount of records
	var c string
	for rows.Next() {
		if err := rows.Scan(&c); err != nil {
			p.tracer.LogError(span, err)
			return nil, err
		} else {
			countries = append(countries, c)
		}
	}

	return &countries, nil
}

// TODO: finish
func (p *PostgresDB) GetListOfCities(ctx context.Context, lang string, countryIso string, find string, first int, after int) (*[]models.City, error) {
	span := p.tracer.MakeSpan(ctx, "")
	defer span.Finish()

	rows, err := p.queryContext(ctx, getQueryGetListOfCities(lang, strings.Title(find), first, after), strings.ToUpper(countryIso) /*after,*/, strings.Title(find))
	if err != nil {
		p.tracer.LogError(span, err)
		return nil, err
	}
	defer rows.Close()
	cities := make([]models.City, 0) // TODO: count amount of records
	var c models.City
	for rows.Next() {
		if err := rows.Scan(&c.Id, &c.City, &c.Subdivision); err != nil {
			p.tracer.LogError(span, err)
			return nil, err
		} else {
			cities = append(cities, c)
		}
	}

	return &cities, nil
}

func (p *PostgresDB) GetListOfAllCities(ctx context.Context, lang string, find string) (*[]models.City, error) {
	span := p.tracer.MakeSpan(ctx, "")
	defer span.Finish()

	rows, err := p.queryContext(ctx, getQueryGetListOfAllCities(lang), strings.Title(find))
	if err != nil {
		p.tracer.LogError(span, err)
		return nil, err
	}
	defer rows.Close()
	cities := make([]models.City, 0) // TODO: count amount of records
	var c models.City
	for rows.Next() {
		if err := rows.Scan(&c.Id, &c.City, &c.Subdivision, &c.Country); err != nil {
			p.tracer.LogError(span, err)
			return nil, err
		} else {
			cities = append(cities, c)
		}
	}

	return &cities, nil
}

// getQueryGetListOfCountryCodes
func (p *PostgresDB) GetListOfCountryCodes(ctx context.Context, lang string) (*[]models.CountryCode, error) {
	span := p.tracer.MakeSpan(ctx, "")
	defer span.Finish()

	rows, err := p.queryContext(ctx, getQueryGetListOfCountryCodes(lang))
	if err != nil {
		p.tracer.LogError(span, err)
		return nil, err
	}
	defer rows.Close()
	countryCodes := make([]models.CountryCode, 0) // TODO: count amount of records
	var c models.CountryCode
	for rows.Next() {
		if err := rows.Scan(&c.Id, &c.CountryCode, &c.Country); err != nil {
			p.tracer.LogError(span, err)
			return nil, err
		} else {
			countryCodes = append(countryCodes, c)
		}
	}

	return &countryCodes, nil
}

func (p *PostgresDB) GetCountryCodeByID(ctx context.Context, id uint) (*models.CountryCode, error) {
	span := p.tracer.MakeSpan(ctx, "")
	defer span.Finish()

	rows, err := p.queryContext(ctx, getQueryGetCountryCodeByID(), id)
	if err != nil {
		p.tracer.LogError(span, err)
		return nil, err
	}
	defer rows.Close()
	var c models.CountryCode
	for rows.Next() {
		if err := rows.Scan(&c.CountryCode, &c.Country); err != nil {
			p.tracer.LogError(span, err)
			return nil, err
		}
	}
	log.Println(c)
	return &c, nil
}

func (p *PostgresDB) GetCityInfoByID(ctx context.Context, id uint, lang string) (*models.City, error) {
	span := p.tracer.MakeSpan(ctx, "")
	defer span.Finish()

	rows, err := p.queryContext(ctx, getQueryGetCityInfoByID(), id)
	if err != nil {
		p.tracer.LogError(span, err)
		return nil, err
	}
	defer rows.Close()

	var c models.City

	for rows.Next() {
		if err := rows.Scan(&c.Id, &c.Country, &c.Subdivision, &c.City); err != nil {
			p.tracer.LogError(span, err)
			return nil, err
		}
	}

	return &c, nil
}
