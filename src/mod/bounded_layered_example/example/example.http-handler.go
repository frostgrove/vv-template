package example

import (
	"net/http"

	"github.com/frostgrove/vv/app/http/appfiber"
	"github.com/frostgrove/vv/auth"
	"github.com/frostgrove/vv/errs"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const (
	PermRead  auth.Permission = "example:read"
	PermWrite auth.Permission = "example:write"
)

// exampleHttpHandler decodes a request into a command, calls the use case and
// encodes the answer. Errors go back through Fiber's error handler, which
// renders the framework envelope (400/404/409/422 with a code and a field).
type exampleHttpHandler struct {
	service *service
}

func newHttpHandler(service *service) (appfiber.Route, error) {
	handler := &exampleHttpHandler{service: service}
	return appfiber.Routes("/examples").
		GET("", appfiber.Requires(PermRead), handler.listOpen).
		POST("", appfiber.Requires(PermWrite), handler.create).
		POST("/:id/complete", appfiber.Requires(PermWrite), handler.complete).
		Route()
}

func (this *exampleHttpHandler) listOpen(c fiber.Ctx) error {
	open, err := this.service.listOpen(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(open)
}

func (this *exampleHttpHandler) create(c fiber.Ctx) error {
	var command createExampleCommand
	if err := c.Bind().JSON(&command); err != nil {
		return errs.BadRequest().Code(errs.CodeMalformedBody).Message("the body is not a JSON object").Fault()
	}
	created, err := this.service.create(c.Context(), command)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(created)
}

func (this *exampleHttpHandler) complete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errs.BadRequest().Code(errs.CodeInvalidID).Field("id").Message("the id is not a UUID").Fault()
	}
	completed, err := this.service.complete(c.Context(), completeExampleCommand{Id: id})
	if err != nil {
		return err
	}
	return c.JSON(completed)
}
