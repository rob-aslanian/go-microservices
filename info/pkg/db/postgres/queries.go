package postgres

import (
	"strconv"
)

const MAX_RECORDS = 50

// TODO: use cursor if possible and Kallax or gorm

func getQueryGetListOfCountries(language string) string {
	return `SELECT country_iso_code FROM ` + countriesTable(language) + ` WHERE country_iso_code IS NOT NULL;`
}

// retrun query string with $1 for `WHERE country_iso_code = $1` as placeholder,
// $2 for `WHERE geoname_id = $2` as placeholder,
// $3 for `AND city_name LIKE $3 || '%' as placeholder
func getQueryGetListOfCities(language string, findCity string, first int, after int) string {
	if first > MAX_RECORDS {
		first = MAX_RECORDS
	}
	if first < 0 {
		first = 0
	}

	var findNameOfCity, q string
	if after != 0 {
		// findNameOfCity = ` AND city_name > (SELECT city_name FROM ` + citiesTable(language) + ` WHERE geoname_id = $2 LIMIT 1) `
		findNameOfCity = ` AND city_name > (SELECT city_name FROM ` + citiesTable(language) + `  LIMIT 1) `
	}

	if findCity != "" {
		q = `SELECT geoname_id, city_name, subdivision_1_name FROM ` + citiesTable(language)
		q += ` WHERE country_iso_code = $1`
		q += findNameOfCity
		q += ` AND city_name LIKE $2 || '%' `
		q += `ORDER BY city_name ASC LIMIT ` + strconv.Itoa(first) + `;`
		return q
	}

	q = `SELECT geoname_id, city_name, subdivision_1_name FROM ` + citiesTable(language)
	q += ` WHERE country_iso_code = $1`
	q += findNameOfCity
	q += ` ORDER BY city_name ASC LIMIT ` + strconv.Itoa(first) + `;`
	return q
}

// retrun query string with $1 for `WHERE  city_name LIKE $1 || '%' as placeholder
func getQueryGetListOfAllCities(language string) string {
	q := `SELECT geoname_id, city_name, subdivision_1_name, country_iso_code FROM ` + citiesTable(language) + ` `
	q += `WHERE  city_name LIKE $1 || '%' `
	q += `ORDER BY city_name ASC;`
	return q
}

func getQueryGetListOfCountryCodes(language string) string {
	q := `SELECT code.id, code.country_code, country.country_iso_code FROM country_codes `
	q += `AS code INNER JOIN ` + countriesTable(language) + ` `
	q += `AS country ON code.geoname_id = country.geoname_id AND country_iso_code IS NOT NULL ORDER BY country.country_iso_code ASC;`
	return q
}

// retrun query string with $1 for `WHERE code.id=$1` as placeholder
func getQueryGetCountryCodeByID() string {
	q := `SELECT code.country_code, country.country_iso_code FROM country_codes `
	q += `AS code INNER JOIN ` + countriesTable("") + ` `
	q += `AS country ON code.geoname_id = country.geoname_id AND country_iso_code IS NOT NULL WHERE code.id = $1 `
	q += `ORDER BY country.country_name ASC, code.country_code ASC;`
	return q
}

// retrun query string with $1 for `WHERE geoname_id=$1` as placeholder
func getQueryGetCityInfoByID() string {
	return `SELECT geoname_id, country_iso_code, subdivision_1_iso_code, city_name FROM cities_en WHERE geoname_id=$1;`
}
