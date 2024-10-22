package main

import (
	"github.com/LissaGreense/discogs_record_label/backend/api"
	"github.com/LissaGreense/discogs_record_label/backend/graphQL"
	"github.com/LissaGreense/discogs_record_label/backend/storage"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	db, err := storage.InitDatabase()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	if err := storage.CreateSchema(db); err != nil {
		log.Fatalf("Error creating schema: %v", err)
	}

	labelId := getLabelID(err)

	if err := api.FetchAndStoreReleases(db, labelId); err != nil {
		log.Fatalf("Error fetching and storing releases: %v", err)
	}

	log.Println("Finished fetching and storing releases.")

	schemaConfig := graphql.SchemaConfig{
		Query: graphQL.NewQueryType(db),
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create new schema, error: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", enableCors(h))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsOrigin := os.Getenv("CORS_ORIGIN")
		if corsOrigin == "" {
			corsOrigin = "http://localhost:5173"
		}
		w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getLabelID(err error) int {
	labelIdStr := os.Getenv("SELECTED_LABEL")

	if labelIdStr == "" {
		log.Fatal("SELECTED_LABEL environment variable is not set")
	}

	labelId, err := strconv.Atoi(labelIdStr)
	if err != nil {
		log.Fatalf("SELECTED_LABEL must be an integer, got: %v", labelIdStr)
	}
	return labelId
}
