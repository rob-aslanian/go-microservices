package postgres

// retrun countries table name by language
func countriesTable(language string) string {
	var tableName string
	switch language {
	case "ar":
		//todo
	default:
		tableName = "countries_en"
	}
	return tableName
}

// retrun cities table name by language
func citiesTable(language string) string {
	var tableName string
	switch language {
	case "ar":
		//todo
	default:
		tableName = "cities_en"
	}
	return tableName
}
