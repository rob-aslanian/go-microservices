package network

import (
	"crypto/tls"
	"log"
	"net"
	httpgo "net/http"
	"strings"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type NetworkRepo struct {
	connection               driver.Connection
	client                   driver.Client
	db                       driver.Database
	graph                    driver.Graph
	users                    driver.Collection
	friendships              driver.Collection
	follows                  driver.Collection
	categories               driver.Collection
	categoriesForFollowings  driver.Collection
	category_tree            driver.Collection
	category_tree_followings driver.Collection
	companies                driver.Collection
	owns_company             driver.Collection
	admins                   driver.Collection
	works_at                 driver.Collection
	blocks                   driver.Collection
	recommendationRequest    driver.Collection
	recommendations          driver.Collection
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

	a.category_tree, err = a.createVertexCollection(CategoryTreeName)
	if err != nil {
		return err
	}

	a.category_tree_followings, err = a.createVertexCollection(CategoryTreeForFollowingsName)
	if err != nil {
		return err
	}

	// *********************
	// I think we don't need index on email in arangodb because we already have it in mongodb and we will not insert document in arangodb which is not already in mongodb
	// *********************
	//_,_, err = a.users.EnsureHashIndex(nil, []string{"Email"}, &driver.EnsureHashIndexOptions{Unique: true})
	//if err != nil{
	//	return err
	//}

	a.friendships, err = a.createEdgeCollection(FriendshipName, driver.VertexConstraints{From: []string{UsersName}, To: []string{UsersName}})
	if err != nil {
		return err
	}
	a.friendships.EnsureHashIndex(nil, []string{"_from", "_to"}, &driver.EnsureHashIndexOptions{Unique: true})

	a.follows, err = a.createEdgeCollection(FollowName, driver.VertexConstraints{From: []string{UsersName, CompaniesName}, To: []string{UsersName, CompaniesName}})
	if err != nil {
		return err
	}
	a.follows.EnsureHashIndex(nil, []string{"_from", "_to"}, &driver.EnsureHashIndexOptions{Unique: true})

	a.categories, err = a.createEdgeCollection(CategoriesName, driver.VertexConstraints{From: []string{UsersName}, To: []string{UsersName, CompaniesName}})
	if err != nil {
		return err
	}
	a.categories.EnsureHashIndex(nil, []string{"_from", "_to", "name"}, &driver.EnsureHashIndexOptions{Unique: true})

	a.categoriesForFollowings, err = a.createEdgeCollection(CategoriesForFollowingsName, driver.VertexConstraints{From: []string{UsersName, CompaniesName}, To: []string{CompaniesName, UsersName}})
	if err != nil {
		return err
	}
	a.categoriesForFollowings.EnsureHashIndex(nil, []string{"_from", "_to", "name"}, &driver.EnsureHashIndexOptions{Unique: true})

	a.owns_company, err = a.createEdgeCollection(OwnsCompanyName, driver.VertexConstraints{From: []string{UsersName}, To: []string{CompaniesName}})
	if err != nil {
		return err
	}
	a.owns_company.EnsureHashIndex(nil, []string{"_from", "_to"}, &driver.EnsureHashIndexOptions{Unique: true})

	a.admins, err = a.createEdgeCollection(AdminsName, driver.VertexConstraints{From: []string{UsersName}, To: []string{CompaniesName}})
	if err != nil {
		return err
	}
	a.admins.EnsureHashIndex(nil, []string{"_from", "_to"}, &driver.EnsureHashIndexOptions{Unique: true})

	a.works_at, err = a.createEdgeCollection(WorksAtName, driver.VertexConstraints{From: []string{UsersName}, To: []string{CompaniesName}})
	if err != nil {
		return err
	}

	a.blocks, err = a.createEdgeCollection(BlocksName, driver.VertexConstraints{From: []string{UsersName}, To: []string{UsersName, CompaniesName}})
	if err != nil {
		return err
	}
	a.blocks.EnsureHashIndex(nil, []string{"_from", "_to"}, &driver.EnsureHashIndexOptions{Unique: true})

	a.recommendationRequest, err = a.createEdgeCollection(RecommendationRequestsName, driver.VertexConstraints{From: []string{UsersName}, To: []string{UsersName}})
	if err != nil {
		return err
	}

	a.recommendations, err = a.createEdgeCollection(RecommendationsName, driver.VertexConstraints{From: []string{UsersName}, To: []string{UsersName}})
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

// Custom HTTP tranport

type transport struct {
	current *httpgo.Request
	tr      *httpgo.Transport
}

func (t transport) RoundTrip(request *httpgo.Request) (response *httpgo.Response, err error) {

	t.current = request

	str := strings.Builder{}
	str.WriteString("URL: ")
	str.WriteString(request.URL.String())
	str.WriteString("\n")

	if request.Method == "POST" {
		str.WriteString("Request:\n")
		request.ParseForm()
		str.WriteString(request.Form.Encode())
		str.WriteString("\n")
	}

	// bodyBuffer, _ := ioutil.ReadAll(request.Body)
	// t.current.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuffer))
	//
	// str.WriteString("Body:\n")
	// str.Write(bodyBuffer)

	log.Println(str.String())

	return t.tr.RoundTrip(request)
}

func createCustomTransport() *transport {
	trans := httpgo.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Proxy: httpgo.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &transport{
		tr: &trans,
	}
}
