package main

import (
	"datastore-fiber-crud/internal/database"
	internal "datastore-fiber-crud/internal/ddl"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	println("START")

	// Initialize Datastore client
	client := database.Connect()
	defer client.Close() // Ensure the client is closed when the application shuts down

	// Use CORS middleware to allow requests from any origin
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allows requests from any origin (useful for local development)
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}))

	// Initialize the controller with the Datastore client
	controller := internal.NewController(client)

	// Register routes
	app.Post("/namespaces", controller.CreateNamespace)
	app.Get("/namespaces", controller.ListNamespaces)

	app.Get("/kinds", controller.ListKinds)
	app.Get("/kinds/:namespace", controller.ListKindByNamespace)
	app.Delete("/kinds/:namespace/:kind", controller.DeleteKind)

	app.Get("/entities/:namespace/:kind", controller.ListEntitiesWithLimit)
	app.Get("/entities/:namespace/:kind/:key", controller.GetEntityByKey)
	app.Post("/entities/filter/:namespace/:kind", controller.FilterEntitiesByFields)

	app.Post("/entities", controller.CreateEntityWithData)
	app.Put("/entities/:namespace/:kind/:id", controller.UpdateEntityByID)
	app.Delete("/entities/:namespace/:kind/:id", controller.DeleteEntityByID)

	// Start the server
	if err := app.Listen(":9000"); err != nil {
		println("Error starting the server:", err.Error())
	}
}
