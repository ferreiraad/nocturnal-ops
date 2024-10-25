// database/database.go
package database

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/datastore"
)

var Client *datastore.Client

func Connect() *datastore.Client {
	ctx := context.Background()

	// Set the Datastore emulator host environment variable
	os.Setenv("DATASTORE_EMULATOR_HOST", "localhost:8000")
	os.Setenv("DATASTORE_PROJECT_ID", "datastore-wrapper")

	// Initialize the Datastore client
	client, err := datastore.NewClient(ctx, "datastore-wrapper") // Replace with your local project ID
	if err != nil {
		log.Fatalf("Failed to connect to Datastore: %v", err)
	}

	log.Println("Connected to the Datastore emulator at localhost:8000")
	return client
}
