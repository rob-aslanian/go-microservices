package maxmind

import (
	"net"
	"strings"

	"gitlab.lan/Rightnao-site/microservices/info/pkg/models"
)

var defaultLanguage = "en"

func (db *DB) GetCityByIp(ip string, language string) (*models.City, error) {
	if language == "" {
		language = defaultLanguage
	}

	pIp := net.ParseIP(ip)

	var record struct {
		City struct {
			GeoNameID uint              `maxminddb:"geoname_id"`
			Names     map[string]string `maxminddb:"names"`
		} `maxminddb:"city"`
		Country struct {
			GeoNameID uint              `maxminddb:"geoname_id"`
			Names     map[string]string `maxminddb:"names"`
			IsoCode   string            `maxminddb:"iso_code"`
		} `maxminddb:"country"`
		Subdivisions []struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"subdivisions"`
	}

	err := db.Lookup(pIp, &record)
	if err != nil {
		return nil, err
	}

	cityName := record.City.Names[language]
	countryName := record.Country.Names[language]
	// get all subdivisions
	var subdivisions []string
	for k, _ := range record.Subdivisions {
		subdivisions = append(subdivisions, record.Subdivisions[k].Names[language])
	}
	subdivision := strings.Join(subdivisions, ", ")

	return &models.City{
		Id:          int32(record.City.GeoNameID),
		City:        &cityName,
		Subdivision: &subdivision,
		Country:     &countryName,
	}, nil
}

func (db *DB) GetCountryByIp(ip string, language string) (*models.CountryCode, error) {
	if language == "" {
		language = defaultLanguage
	}

	pIp := net.ParseIP(ip)

	var record struct {
		Country struct {
			GeoNameID uint              `maxminddb:"geoname_id"`
			Names     map[string]string `maxminddb:"names"`
			IsoCode   string            `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}

	err := db.Lookup(pIp, &record)
	if err != nil {
		return nil, err
	}

	return &models.CountryCode{
		Id:          int32(record.Country.GeoNameID),
		CountryCode: record.Country.IsoCode,
		Country:     record.Country.Names[language],
	}, nil
}
