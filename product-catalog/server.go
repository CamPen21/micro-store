package main

import (
	"log"
	"micro-store/product-catalog/generated"
	"micro-store/product-catalog/resolver"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{}}))

    // http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/", srv)

	// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
