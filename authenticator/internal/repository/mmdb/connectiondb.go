package repo

import (
	"net"

	geoip2 "github.com/oschwald/geoip2-golang"
	"gitlab.lan/Rightnao-site/microservices/authenticator/internal/location"
)

// IPsRepo serves as struct to connect to mongodb
type IPsRepo struct {
	Db *geoip2.Reader
}

// NewIPsRepo recieved ID and returns ...
func NewIPsRepo(address string) (*IPsRepo, error) {
	db, err := geoip2.Open(address)
	if err != nil {
		return nil, err
	}

	return &IPsRepo{
		Db: db,
	}, nil
}

// Close ...
func (d *IPsRepo) Close() {
	d.Db.Close()
}

// GetCityByIP recieves the ip as sstring and returns location type.
func (d IPsRepo) GetCityByIP(ip string) (location.Location, error) {
	c, err := d.Db.City(net.ParseIP((ip)))
	if err != nil {
		return location.Location{}, err
	}

	city := c.City.GeoNameID
	country := c.Country.IsoCode

	return location.Location{
		City:    &city,
		Country: &country,
	}, nil
}
