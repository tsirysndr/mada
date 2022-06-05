package mada

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/rakyll/statik/fs"
	"github.com/tsirysndr/mada/graph"
	"github.com/tsirysndr/mada/graph/generated"
	_ "github.com/tsirysndr/mada/statik"
)

const PORT = 8010

func StartHttpServer(db *sql.DB) {
	statikFS, _ := fs.New()

	router := chi.NewRouter()
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	corsMiddleware := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	index, err := InitializeBleve(db)
	if err != nil {
		panic(err)
	}

	r := &graph.Resolver{
		CommuneService:   NewCommuneService(db, index),
		DistrictService:  NewDistrictService(db, index),
		FokontanyService: NewFokontanyService(db, index),
		RegionService:    NewRegionService(db, index),
		SearchService:    NewSearchService(db, index),
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: r}))

	log.Println("🚀 Starting server on port", PORT)
	log.Printf("🚀 Connect to http://localhost:%d/playground for GraphQL playground", PORT)

	router.Use(corsMiddleware.Handler)
	router.Handle("/", http.FileServer(statikFS))
	router.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)
	router.Handle("/*", http.FileServer(statikFS))
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", PORT), router)
	if err != nil {
		log.Fatal(err)
	}

}
