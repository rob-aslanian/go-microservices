package geoip

import (
	"net"

	geoip2 "github.com/oschwald/geoip2-golang"
)

// Repository represent GeoIP repository
type Repository struct {
	db *geoip2.Reader
}

// NewRepository creates  new instance of repository
func NewRepository(path string) (*Repository, error) {
	db, err := geoip2.Open(path)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

// Close ...
func (d *Repository) Close() {
	d.db.Close()
}

// GetCountryISOCode ...
func (d *Repository) GetCountryISOCode(ipAddress net.IP) (string, error) {

	if len([]byte(ipAddress)) == 0 {
		ipAddress = net.ParseIP("94.43.151.246") // TODO: delete later
	}

	country, err := d.db.Country(ipAddress)
	if err != nil {
		return "", err
	}
	return country.Country.IsoCode, nil
}
