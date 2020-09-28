package arangorepo

import (
	"crypto/tls"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

const (
	DbName        = "hypercube"
	GraphName     = "network"
	UsersName     = "users"
	CompaniesName = "companies"
)

type NetworkRepo struct {
	connection driver.Connection
	client     driver.Client
	db         driver.Database
	graph      driver.Graph
	users      driver.Collection
	companies  driver.Collection
}

func NewNetworkRepo(username, password string, addresses []string) (*NetworkRepo, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: addresses,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		// Transport: createCustomTransport(),
	})
	if err != nil {
		return nil, err
	}
	conn, err = conn.SetAuthentication(driver.BasicAuthentication(username, password))
	if err != nil {
		return nil, err
	}

	c, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
	})

	repo := NetworkRepo{
		connection: conn,
		client:     c,
	}
	err = repo.setUpGraph()
	return &repo, err
}

func (a *NetworkRepo) setUpGraph() error {
	exists, err := a.client.DatabaseExists(nil, DbName)
	if err != nil {
		return err
	}
	if !exists {
		a.db, err = a.client.CreateDatabase(nil, DbName, &driver.CreateDatabaseOptions{})
		if err != nil {
			return err
		}
	} else {
		a.db, err = a.client.Database(nil, DbName)
		if err != nil {
			return err
		}
	}

	exists, err = a.db.GraphExists(nil, GraphName)
	if err != nil {
		return err
	}
	if !exists {
		a.graph, err = a.db.CreateGraph(nil, GraphName, nil)
		if err != nil {
			return err
		}
	} else {
		a.graph, err = a.db.Graph(nil, GraphName)
		if err != nil {
			return err
		}
	}

	a.users, err = a.createVertexCollection(UsersName)
	if err != nil {
		return err
	}

	a.companies, err = a.createVertexCollection(CompaniesName)
	if err != nil {
		return err
	}

	return nil
}

func (a *NetworkRepo) createVertexCollection(name string) (driver.Collection, error) {
	var collection driver.Collection
	exists, err := a.graph.VertexCollectionExists(nil, name)
	if err != nil {
		return nil, err
	}
	if !exists {
		collection, err = a.graph.CreateVertexCollection(nil, name)
		if err != nil {
			return nil, err
		}
	} else {
		collection, err = a.graph.VertexCollection(nil, name)
		if err != nil {
			return nil, err
		}
	}
	return collection, nil
}

func (a *NetworkRepo) createEdgeCollection(name string, constraints driver.VertexConstraints) (driver.Collection, error) {
	var collection driver.Collection
	exists, err := a.graph.EdgeCollectionExists(nil, name)
	if err != nil {
		return nil, err
	}
	if !exists {
		collection, err = a.graph.CreateEdgeCollection(nil, name, constraints)
		if err != nil {
			return nil, err
		}
	} else {
		collection, _, err = a.graph.EdgeCollection(nil, name)
		if err != nil {
			return nil, err
		}
	}
	return collection, nil
}
