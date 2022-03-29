package application

import (
	"github.com/bdemirpolat/kubecd/models"
	"github.com/bdemirpolat/kubecd/validate"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
)

// HandlerInterface contains handler functions
type HandlerInterface interface {
	CreateHandler(c *fiber.Ctx) error
	UpdateHandler(c *fiber.Ctx) error
	GetHandler(c *fiber.Ctx) error
	ListHandler(c *fiber.Ctx) error
	DeleteHandler(c *fiber.Ctx) error
	RegisterHandlers(app *fiber.App)
}

// HandlerStruct
type HandlerStruct struct {
	service ServiceInterface
}

// NewHandler returns new HandlerStruct / HandlerInterface with service
func NewHandler(service ServiceInterface) HandlerInterface {
	return &HandlerStruct{service: service}
}

// compile time proof
var _ HandlerInterface = &HandlerStruct{}

// CreateHandler validates create request and sends to service
func (h *HandlerStruct) CreateHandler(c *fiber.Ctx) error {
	req := models.ApplicationCreateReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPResponseJSON{
			Error: err.Error(),
		})
	}

	err = validate.Validate.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPResponseJSON{
			Error: err.Error(),
		})
	}

	res, err := h.service.CreateService(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPResponseJSON{
			Error: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(models.HTTPResponseJSON{
		Data: res,
	})
}

// UpdateHandler validates update request and sends to service
func (h *HandlerStruct) UpdateHandler(c *fiber.Ctx) error {
	req := models.ApplicationUpdateReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPResponseJSON{
			Error: err.Error(),
		})
	}
	req.ID = cast.ToUint(c.Params("id"))
	err = validate.Validate.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPResponseJSON{
			Error: err.Error(),
		})
	}
	res, err := h.service.UpdateService(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPResponseJSON{
			Error: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(models.HTTPResponseJSON{
		Data: res,
	})
}

// GetHandler validates get request and sends to service
func (h *HandlerStruct) GetHandler(c *fiber.Ctx) error {
	req := models.ApplicationGetReq{}
	req.ID = cast.ToUint(c.Params("id"))

	err := validate.Validate.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPResponseJSON{
			Error: err.Error(),
		})
	}

	res, err := h.service.GetService(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPResponseJSON{
			Error: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(models.HTTPResponseJSON{
		Data: res.Data,
	})
}

// ListHandler validates list request and sends to service
func (h *HandlerStruct) ListHandler(c *fiber.Ctx) error {
	req := models.ApplicationListReq{}
	err := c.QueryParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPResponseJSON{
			Error: err.Error(),
		})
	}
	res, err := h.service.ListService(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPResponseJSON{
			Error: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(models.HTTPResponseJSON{
		Data: res.Data,
	})
}

// DeleteHandler validates get request and sends to service
func (h *HandlerStruct) DeleteHandler(c *fiber.Ctx) error {
	req := models.ApplicationDeleteReq{}
	req.ID = cast.ToUint(c.Params("id"))

	err := validate.Validate.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPResponseJSON{
			Error: err.Error(),
		})
	}

	_, err = h.service.DeleteService(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HTTPResponseJSON{
			Error: err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// RegisterHandlers registers handlers to fiber app
func (h *HandlerStruct) RegisterHandlers(app *fiber.App) {
	app.Post("/applications", h.CreateHandler)
	app.Put("/applications/:id", h.UpdateHandler)
	app.Get("/applications/:id", h.GetHandler)
	app.Get("/applications", h.ListHandler)
	app.Delete("/applications/:id", h.DeleteHandler)
}
