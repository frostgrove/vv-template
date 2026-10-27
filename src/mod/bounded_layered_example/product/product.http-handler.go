package product

import (
	"github.com/frostgrove/vv/app/http/appfiber"
	"github.com/frostgrove/vv/crud/decorators/specs"
	"github.com/frostgrove/vv/crud/http/crudfiber"
	"github.com/frostgrove/vv/crud/query"
	"github.com/google/uuid"
)

// NewHttpHandler exposes Product as a CRUD resource: list with the query DSL,
// query, count, get, create, patch, replace, delete and bulk delete.
//
// Every operation goes through the registrar, so one appfiber.Policy is both
// the access declaration the boot gate compares with Fiber's routing table and
// the check that runs in front of the handler (D-073, D-100). A route mounted
// past it, or a declaration with no route, fails the start.
func NewHttpHandler(repository *ProductRepo) (appfiber.Route, error) {
	resource := crudfiber.New(specs.Executor(repository),
		crudfiber.WithQuery[Product, uuid.UUID, ProductUpdate](&query.Config{
			Filterable: []string{"Name", "Description", "CreatedAt", "UpdatedAt"},
			Sortable:   []string{"Name", "CreatedAt", "UpdatedAt"},
			Searchable: []string{"Name", "Description"},
		}),
	)

	read := appfiber.Requires(PermRead)
	write := appfiber.Requires(PermWrite)
	remove := appfiber.Requires(PermDelete)
	return appfiber.Routes("/products").
		GET("", read, resource.List).
		POST("/query", read, resource.Query).
		GET("/count", read, resource.CountGet).
		POST("/count", read, resource.CountPost).
		GET("/:id", read, resource.GetByID).
		POST("", write, resource.Create).
		PATCH("/:id", write, resource.Update).
		PUT("/:id", write, resource.Replace).
		DELETE("/:id", remove, resource.Delete).
		POST("/bulk-delete", remove, resource.BulkDelete).
		Route()
}
