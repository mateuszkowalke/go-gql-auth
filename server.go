package main

import (
	"log"
	"net/http"
	"os"

	"example.com/go-graphql-auth/auth"
	"example.com/go-graphql-auth/database"
	"example.com/go-graphql-auth/graph"
	"example.com/go-graphql-auth/graph/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	database.InitDb()
	database.Migrate()

	mux := http.NewServeMux()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	mux.Handle("/", auth.Middleware(playground.Handler("GraphQL playground", "/query")))
	mux.Handle("/query", auth.Middleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
