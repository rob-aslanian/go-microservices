package models

type Country struct {
	Id      string
	Country *string
}

type CountryCode struct {
	Id          int32
	CountryCode string
	Country     string
}
