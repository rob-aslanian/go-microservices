package main

import (
	"testing"

	graphql "github.com/graph-gophers/graphql-go"
	"gitlab.lan/Rightnao-site/microservices/graphql/resolver"
	"gitlab.lan/Rightnao-site/microservices/graphql/schema"
)

// check resolvers according to graphQL schema
func TestGraphQL(t *testing.T) {
	res := resolver.Resolver{}
	res.Init()

	graphql.MustParseSchema(
		schema.GetRootSchema(),
		&res,
	)
}
