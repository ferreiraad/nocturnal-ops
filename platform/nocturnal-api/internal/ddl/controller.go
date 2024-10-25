package internal

import (
	"context"
	"strconv"

	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/gofiber/fiber/v2"
)

// Controller holds the datastore client and other dependencies
type Controller struct {
	client *datastore.Client
}

// NewController initializes a new controller with the Datastore client
func NewController(client *datastore.Client) *Controller {
	return &Controller{client: client}
}

type NamespacePlaceholder struct {
	ID      int64  `datastore:"-"`
	Message string `datastore:"message"`
}

// CreateNamespace "creates" a namespace by adding a placeholder entity
func (ctr *Controller) CreateNamespace(c *fiber.Ctx) error {
	ctx := context.Background()

	// Parse the namespace from the request body
	type NamespaceRequest struct {
		Namespace string `json:"namespace" validate:"required"`
	}
	var req NamespaceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if req.Namespace == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Namespace is required"})
	}

	// Check if the namespace already exists using FilterField
	query := datastore.NewQuery("__namespace__").FilterField("__key__", "=", datastore.NameKey("__namespace__", req.Namespace, nil)).KeysOnly()
	keys, err := ctr.client.GetAll(ctx, query, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check namespace existence"})
	}

	// If keys are returned, the namespace already exists
	if len(keys) > 0 {
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("Namespace '%s' already exists", req.Namespace)})
	}

	// Create a placeholder entity in the specified namespace
	placeholder := NamespacePlaceholder{
		Message: "This is a placeholder entity for namespace creation",
	}
	key := datastore.IncompleteKey("NamespacePlaceholder", nil)
	key.Namespace = req.Namespace

	_, err = ctr.client.Put(ctx, key, &placeholder)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create namespace"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":   "Namespace created successfully",
		"namespace": req.Namespace,
	})
}

// ListNamespaces lists all namespaces in Datastore
func (ctr *Controller) ListNamespaces(c *fiber.Ctx) error {
	ctx := context.Background()
	var namespaces []string

	// Query for all namespaces using the special __namespace__ kind
	query := datastore.NewQuery("__namespace__").KeysOnly()
	keys, err := ctr.client.GetAll(ctx, query, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve namespaces"})
	}

	// Extract namespace IDs as strings
	for _, key := range keys {
		if key.ID > 0 {
			namespaces = append(namespaces, fmt.Sprintf("%d", key.ID))
		} else {
			namespaces = append(namespaces, key.Name)
		}
	}

	return c.JSON(fiber.Map{"namespaces": namespaces})
}

// EntityDefinition represents the structure of an entity
type EntityDefinition struct {
	Namespace string            `json:"namespace"`
	Kind      string            `json:"kind"`
	Fields    map[string]string `json:"fields"`
}

// ListEntities lists all entity kinds within a specific namespace
func (ctr *Controller) ListEntities(c *fiber.Ctx) error {
	ctx := context.Background()
	var kinds []string

	// Get the namespace from the path parameters
	namespace := c.Params("namespace")
	if namespace == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Namespace is required"})
	}

	// Query for all kinds in the specified namespace using __kind__ kind
	query := datastore.NewQuery("__kind__").Namespace(namespace).KeysOnly()
	keys, err := ctr.client.GetAll(ctx, query, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve entity kinds"})
	}

	// Extract kind names from the keys
	for _, key := range keys {
		kinds = append(kinds, key.Name)
	}

	return c.JSON(fiber.Map{"kinds": kinds})
}

// CreateEntity dynamically defines an entity in Datastore
func (ctr *Controller) CreateEntity(c *fiber.Ctx) error {
	ctx := context.Background()

	// Parse the incoming JSON to define the entity
	var entityDef EntityDefinition
	if err := c.BodyParser(&entityDef); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if entityDef.Kind == "" || len(entityDef.Fields) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid entity definition"})
	}

	// Prepare entity data for storage in Datastore based on dynamic fields
	entityData := make(map[string]interface{})
	for field, fieldType := range entityDef.Fields {
		// Default values based on field type (for demonstration)
		switch fieldType {
		case "string":
			entityData[field] = ""
		case "int":
			entityData[field] = 0
		case "float":
			entityData[field] = 0.0
		default:
			return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("Unsupported field type: %s", fieldType)})
		}
	}

	// Create a new Datastore key with the specified kind and namespace
	key := datastore.IncompleteKey(entityDef.Kind, nil)
	if entityDef.Namespace != "" {
		key.Namespace = entityDef.Namespace
	}

	// Insert the entity structure into Datastore
	_, err := ctr.client.Put(ctx, key, &entityData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create entity in Datastore"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Entity created successfully",
		"entity":  entityDef,
	})
}

// CreateEntityWithData creates an entity in Datastore with a specified kind and data
func (ctr *Controller) CreateEntityWithData(c *fiber.Ctx) error {
	ctx := context.Background()

	// Define the request structure
	type EntityRequest struct {
		Namespace string                 `json:"namespace"`
		Kind      string                 `json:"kind" validate:"required"`
		Data      map[string]interface{} `json:"data" validate:"required"`
	}
	var req EntityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Validate request parameters
	if req.Kind == "" || len(req.Data) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Kind and data are required"})
	}

	// Convert Data map to datastore.PropertyList for dynamic fields with compatible types
	var properties datastore.PropertyList
	for key, value := range req.Data {
		switch v := value.(type) {
		case string:
			properties = append(properties, datastore.Property{Name: key, Value: v})
		case int, int64, float64, bool:
			properties = append(properties, datastore.Property{Name: key, Value: v})
		default:
			return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("Unsupported field type for key '%s'", key)})
		}
	}

	// Use an incomplete key for Datastore to auto-generate an ID
	key := datastore.IncompleteKey(req.Kind, nil)
	if req.Namespace != "" {
		key.Namespace = req.Namespace
	}

	// Insert the entity into Datastore and let it generate a unique key
	createdKey, err := ctr.client.Put(ctx, key, &properties)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create entity: %v", err)})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Entity created successfully",
		"kind":    req.Kind,
		"data":    req.Data,
		"id":      createdKey.ID, // Return the generated ID
	})
}

// DeleteKind deletes all entities within a specific namespace and kind
func (ctr *Controller) DeleteKind(c *fiber.Ctx) error {
	ctx := context.Background()

	// Get the namespace and kind from the path parameters
	namespace := c.Params("namespace")
	kind := c.Params("kind")

	// Validate that both parameters are provided
	if namespace == "" || kind == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Namespace and kind are required"})
	}

	// Query for all entities in the specified kind and namespace
	query := datastore.NewQuery(kind).Namespace(namespace).KeysOnly()
	keys, err := ctr.client.GetAll(ctx, query, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve entities for deletion"})
	}

	// If no entities are found, return a 404 response
	if len(keys) == 0 {
		return c.Status(404).JSON(fiber.Map{"message": fmt.Sprintf("No entities found in kind '%s' in namespace '%s'", kind, namespace)})
	}

	// Delete all entities in the specified kind
	if err := ctr.client.DeleteMulti(ctx, keys); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete all entities in kind"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": fmt.Sprintf("All entities in kind '%s' in namespace '%s' deleted successfully", kind, namespace),
	})
}

// // DeleteEntity deletes an entity within a specific namespace and kind
// func (ctr *Controller) DeleteEntity(c *fiber.Ctx) error {
// 	ctx := context.Background()

// 	// Get the namespace, kind, and entity from the path parameters
// 	namespace := c.Params("namespace")
// 	kind := c.Params("kind")
// 	entity := c.Params("entity")

// 	// Validate that all parameters are provided
// 	if namespace == "" || kind == "" || entity == "" {
// 		return c.Status(400).JSON(fiber.Map{"error": "Namespace, kind, and entity are required"})
// 	}

// 	// Construct a Datastore key with the specified namespace, kind, and entity ID
// 	key := datastore.NameKey(kind, entity, nil)
// 	key.Namespace = namespace

// 	// Attempt to delete the entity
// 	if err := ctr.client.Delete(ctx, key); err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete entity"})
// 	}

// 	return c.Status(200).JSON(fiber.Map{
// 		"message": fmt.Sprintf("Entity '%s' of kind '%s' deleted successfully from namespace '%s'", entity, kind, namespace),
// 	})
// }

// ListEntitiesWithLimit lists entities of a specific kind within a namespace with an optional limit
func (ctr *Controller) ListEntitiesWithLimit(c *fiber.Ctx) error {
	ctx := context.Background()

	// Get the namespace and kind from the path parameters
	namespace := c.Params("namespace")
	kind := c.Params("kind")
	if namespace == "" || kind == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Namespace and kind are required"})
	}

	// Get the limit from the query parameter (optional, default to 10)
	limit := c.QueryInt("limit", 10)
	if limit <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Limit must be a positive integer"})
	}

	// Query for the specified kind within the namespace and apply the limit
	query := datastore.NewQuery(kind).Namespace(namespace).Limit(limit)
	var entities []datastore.PropertyList
	keys, err := ctr.client.GetAll(ctx, query, &entities)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve entities", "details": err.Error()})
	}

	// Format the response to include the key ID or name along with the entity properties
	var results []fiber.Map
	for i, key := range keys {
		entityData := make(fiber.Map)
		for _, property := range entities[i] {
			entityData[property.Name] = property.Value
		}
		results = append(results, fiber.Map{
			"key":  key.ID, // Use key.Name if key.ID is zero (for named keys)
			"data": entityData,
		})
	}

	return c.JSON(fiber.Map{"entities": results})
}

// GetEntityByKey retrieves a single entity by its key within a specific namespace and kind
func (ctr *Controller) GetEntityByKey(c *fiber.Ctx) error {
	ctx := context.Background()

	// Get the namespace, kind, and key from the path parameters
	namespace := c.Params("namespace")
	kind := c.Params("kind")
	keyParam := c.Params("key")

	// Validate parameters
	if namespace == "" || kind == "" || keyParam == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Namespace, kind, and key are required"})
	}

	// Parse the key as either an ID or name
	var datastoreKey *datastore.Key
	if id, err := strconv.ParseInt(keyParam, 10, 64); err == nil {
		// Key is an integer ID
		datastoreKey = datastore.IDKey(kind, id, nil)
	} else {
		// Key is a string name
		datastoreKey = datastore.NameKey(kind, keyParam, nil)
	}
	datastoreKey.Namespace = namespace

	// Retrieve the entity from Datastore
	var entity datastore.PropertyList
	if err := ctr.client.Get(ctx, datastoreKey, &entity); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Entity not found", "details": err.Error()})
	}

	// Convert the PropertyList to a map for JSON response
	entityData := make(fiber.Map)
	for _, property := range entity {
		entityData[property.Name] = property.Value
	}

	return c.JSON(fiber.Map{
		"key":  datastoreKey.ID, // Use datastoreKey.Name if datastoreKey.ID is zero (for named keys)
		"data": entityData,
	})
}

// UpdateEntityByID updates an entity by its ID within a specific namespace and kind
func (ctr *Controller) UpdateEntityByID(c *fiber.Ctx) error {
	ctx := context.Background()

	// Get the namespace, kind, and id from the path parameters
	namespace := c.Params("namespace")
	kind := c.Params("kind")
	idParam := c.Params("id")

	// Validate that all parameters are provided
	if namespace == "" || kind == "" || idParam == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Namespace, kind, and ID are required"})
	}

	// Parse the ID as an integer
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID must be a valid integer"})
	}

	// Parse the request body to get the new data
	var updatedData map[string]interface{}
	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON", "details": err.Error()})
	}

	// Construct the Datastore key with the specified namespace, kind, and ID
	key := datastore.IDKey(kind, id, nil)
	key.Namespace = namespace

	// Retrieve the existing entity
	var entity datastore.PropertyList
	if err := ctr.client.Get(ctx, key, &entity); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Entity not found", "details": err.Error()})
	}

	// Update the entity's properties with the new data
	for _, property := range entity {
		if newValue, exists := updatedData[property.Name]; exists {
			property.Value = newValue
		}
	}
	// Convert updated data back to PropertyList for saving
	var newProperties datastore.PropertyList
	for name, value := range updatedData {
		newProperties = append(newProperties, datastore.Property{
			Name:  name,
			Value: value,
		})
	}

	// Save the updated entity back to Datastore
	if _, err := ctr.client.Put(ctx, key, &newProperties); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update entity", "details": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":     fmt.Sprintf("Entity with ID '%d' of kind '%s' updated successfully in namespace '%s'", id, kind, namespace),
		"updatedData": updatedData,
	})
}

// DeleteEntityByID deletes an entity by its ID within a specific namespace and kind
func (ctr *Controller) DeleteEntityByID(c *fiber.Ctx) error {
	ctx := context.Background()

	// Get the namespace, kind, and id from the path parameters
	namespace := c.Params("namespace")
	kind := c.Params("kind")
	idParam := c.Params("id")

	// Validate that all parameters are provided
	if namespace == "" || kind == "" || idParam == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Namespace, kind, and ID are required"})
	}

	// Parse the ID as an integer
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID must be a valid integer"})
	}

	// Construct the Datastore key with the specified namespace, kind, and ID
	key := datastore.IDKey(kind, id, nil)
	key.Namespace = namespace

	// Attempt to delete the entity
	if err := ctr.client.Delete(ctx, key); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete entity", "details": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": fmt.Sprintf("Entity with ID '%d' of kind '%s' deleted successfully from namespace '%s'", id, kind, namespace),
	})
}
