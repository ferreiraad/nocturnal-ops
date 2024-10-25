package database

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/datastore"
)

// Global datastore client
var Client *datastore.Client

// Connect initializes and returns a Datastore client
func Connect() *datastore.Client {
	ctx := context.Background()

	// Default Datastore emulator host for local development
	datastoreHost := "localhost:8000"

	// Set environment variables
	os.Setenv("DATASTORE_EMULATOR_HOST", datastoreHost)
	os.Setenv("DATASTORE_PROJECT_ID", "datastore-wrapper")

	// Initialize the Datastore client
	client, err := datastore.NewClient(ctx, "datastore-wrapper")
	if err != nil {
		log.Fatalf("Failed to connect to Datastore: %v", err)
	}

	log.Printf("Connected to the Datastore emulator at %s", datastoreHost)
	return client
}
