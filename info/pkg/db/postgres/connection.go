package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"gitlab.lan/Rightnao-site/microservices/info/pkg/tracing"

	_ "github.com/lib/pq"
)

type ConnectionInfo struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
	SSLMode  bool
	Tracer   *tracing.Tracer
}

type PostgresDB struct {
	db     *sql.DB
	tracer *tracing.Tracer

	getListOfCitiesStmt       *sql.Stmt
	getListOfAllCitiesStmt    *sql.Stmt
	getCityInfoByIDStmt       *sql.Stmt
	getListOfCountryCodesStmt *sql.Stmt
	getCountryCodeByIDStmt    *sql.Stmt
}

func Connect(c *ConnectionInfo) (*PostgresDB, error) {
	sslMode := "disable"
	if c.SSLMode {
		sslMode = "enable"
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.Name, sslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
		return &PostgresDB{}, err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Coudn't ping db:", err)
		return &PostgresDB{}, err
	}

	instance := PostgresDB{
		db:     db,
		tracer: c.Tracer,
	}

	instance.prepareStatements()

	return &instance, nil
}

// Close ...
func (p *PostgresDB) Close() error {
	return p.db.Close()
}

func (p *PostgresDB) queryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	span := p.tracer.MakeSpan(ctx, "Postgress Query")
	defer span.Finish()

	span.SetTag("postgres.query", query)
	span.SetTag("postgres.args", args)

	return p.db.QueryContext(ctx, query, args...)
}

func (p *PostgresDB) prepareStatements() {
	var err error

	// $1 - language (en)
	// $2 - search query (tbili)
	// $3 - country (GE)
	// $4 - amount of results (10)
	// $5 - offset (0)
	p.getListOfCitiesStmt, err = p.db.Prepare(`
		SELECT DISTINCT(geo.geonameid), adm_codes.name AS subdivision, geo.country, geo.population,
		COALESCE(
		  (
		      SELECT alt.alternatename
		      FROM alternatename AS alt
		      WHERE alt.geonameid = geo.geonameid AND "isolanguage" = $1
		      ORDER BY alt.ispreferredname ASC
		      LIMIT 1
		  ),
		  geo.name
		) AS translation
		FROM geoname AS geo
		INNER JOIN alternatename AS alt ON geo.geonameid = alt.geonameid
		INNER JOIN admin1codes AS adm_codes ON CONCAT("country", '.' , geo.admin1) = adm_codes.admin1
		WHERE fclass = 'P' -- only cities
		  AND (
				alt.gin_index @@ to_tsquery('simple', CONCAT( $2::text , ':*'))
				AND
				geo.country = $3
			)
		ORDER BY geo.population DESC
		LIMIT $4 OFFSET $5;
		`)
	if err != nil {
		panic(err)
	}

	// $1 - language (en)
	// $2 - search query (tbili)
	// $3 - amount of results (10)
	// $4 - offset (0)
	p.getListOfAllCitiesStmt, err = p.db.Prepare(`
		SELECT DISTINCT(geo.geonameid), adm_codes.name AS subdivision, geo.country, geo.population,
		COALESCE(
		  (
		      SELECT alt.alternatename
		      FROM alternatename AS alt
		      WHERE alt.geonameid = geo.geonameid AND "isolanguage" = $1
		      ORDER BY alt.ispreferredname ASC
		      LIMIT 1
		  ),
		  geo.name
		) AS translation
		FROM geoname AS geo
		INNER JOIN alternatename AS alt ON geo.geonameid = alt.geonameid
		INNER JOIN admin1codes AS adm_codes ON CONCAT("country", '.' , geo.admin1) = adm_codes.admin1
		WHERE fclass = 'P'
		  AND (
				alt.gin_index @@ to_tsquery('simple', CONCAT( $2::text , ':*'))
			)
		ORDER BY geo.population DESC
		LIMIT $3 OFFSET $4;
	`)
	if err != nil {
		panic(err)
	}

	// $1 - language (en)
	// $2 - city id (611717)
	p.getCityInfoByIDStmt, err = p.db.Prepare(`
		SELECT DISTINCT(geo.geonameid), adm_codes.name AS subdivision, geo.country,
		COALESCE(
		  (
		      SELECT alt.alternatename
		      FROM alternatename AS alt
		      WHERE alt.geonameid = geo.geonameid AND "isolanguage" = $1
		      ORDER BY alt.ispreferredname ASC
		      LIMIT 1
		  ),
		  geo.name
		) AS translation
		FROM geoname AS geo
		INNER JOIN alternatename AS alt ON geo.geonameid = alt.geonameid
		INNER JOIN admin1codes AS adm_codes ON CONCAT("country", '.' , geo.admin1) = adm_codes.admin1
		WHERE fclass = 'P' -- only cities
		  AND geo.geonameid = $2
		LIMIT 1;
	`)
	if err != nil {
		panic(err)
	}

	p.getListOfCountryCodesStmt, err = p.db.Prepare(`
		SELECT code.id, code.country_code, geo.country FROM country_codes
		AS code INNER JOIN geoname AS geo ON code.geoname_id = geo.geonameid
		AND geo.country IS NOT NULL
		ORDER BY geo.country ASC;
	`)
	if err != nil {
		panic(err)
	}

	// $1 - country code id
	p.getCountryCodeByIDStmt, err = p.db.Prepare(`
		SELECT code.country_code, geo.country FROM country_codes
		AS code INNER JOIN geoname AS geo ON code.geoname_id = geo.geonameid
		AND geo.country IS NOT NULL
		WHERE code.id = $1
		ORDER BY geo.country ASC
		LIMIT 1;
		`)
	if err != nil {
		panic(err)
	}

}
