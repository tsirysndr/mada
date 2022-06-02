package mada

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/rakyll/statik/fs"
	_ "github.com/tsirysndr/mada/statik"
)

func StartHttpServer() {
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
	router.Use(corsMiddleware.Handler)
	router.Handle("/", http.FileServer(statikFS))
	router.Handle("/*", http.FileServer(statikFS))
	err := http.ListenAndServe("0.0.0.0:8010", router)
	if err != nil {
		log.Fatal(err)
	}

}
