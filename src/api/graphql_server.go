package api

import (
	"fmt"
	"github.com/rigglo/gql"
	"github.com/rigglo/gql/pkg/handler"
	"github.com/rigglo/gqlws"
	"net/http"
)

func StartGraphql(host string, port int) {
	schema := &gql.Schema{
		Query:        queryType,
		Subscription: subscriptionsType,
		Mutation:     mutationType,
	}

	executor := gql.DefaultExecutor(schema)

	graphQLHandler := handler.New(handler.Config{
		Executor:   gql.DefaultExecutor(schema),
		Playground: false,
	})
	wsQL := gqlws.New(
		gqlws.Config{
			Subscriber: executor.Subscribe,
		},
		graphQLHandler)

	graphIQLHandler := handler.New(handler.Config{
		Executor:   gql.DefaultExecutor(schema),
		Playground: true,
	})
	wsIQL := gqlws.New(
		gqlws.Config{
			Subscriber: executor.Subscribe,
		},
		graphIQLHandler)

	http.Handle("/graphql", wsQL)
	http.Handle("/graphiql", wsIQL)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil); err != nil {
		panic(err)
	}
}
