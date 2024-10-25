package main

import (
	"datastore-fiber-crud/internal/database"
	internal "datastore-fiber-crud/internal/ddl"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	println(" START")

	// Initialize Datastore client
	client := database.Connect()
	defer client.Close() // Ensure the client is closed when the application shuts down

	// Initialize the controller with the Datastore client
	controller := internal.NewController(client)

	// Register routes
	app.Post("/namespaces", controller.CreateNamespace)
	app.Get("/namespaces", controller.ListNamespaces)

	app.Get("/kinds/:namespace", controller.ListEntities)
	app.Delete("/kinds/:namespace/:kind", controller.DeleteKind)

	// app.Delete("/entities/:namespace/:kind/:entity", controller.DeleteEntity)

	app.Get("/entities/:namespace/:kind", controller.ListEntitiesWithLimit)
	app.Get("/entities/:namespace/:kind/:key", controller.GetEntityByKey)
	app.Post("/entities", controller.CreateEntityWithData)
	app.Put("/entities/:namespace/:kind/:id", controller.UpdateEntityByID)
	app.Delete("/entities/:namespace/:kind/:id", controller.DeleteEntityByID)

	if err := app.Listen(":9000"); err != nil {
		println("Error starting the server:", err.Error())
	}

}
