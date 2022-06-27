package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/hendraprawira/go-procurement/directives"
	"github.com/hendraprawira/go-procurement/graph"
	"github.com/hendraprawira/go-procurement/graph/generated"
	middlewares "github.com/hendraprawira/go-procurement/middleware"
)

const defaultPort = "3000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := mux.NewRouter()
	router.Use(middlewares.AuthMiddleware)

	schema := generated.Config{Resolvers: &graph.Resolver{}}
	schema.Directives.Auth = directives.Auth

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(schema))

	buildHandler := http.FileServer(http.Dir("./app-web/dist/"))
	staticHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir("./app-web/dist/assets")))

	router.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)
	router.Handle("/", buildHandler)
	router.Handle("/assets", staticHandler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
