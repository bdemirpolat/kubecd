package application

import (
	"github.com/bdemirpolat/kubecd/pkg/models"
	"github.com/bdemirpolat/kubecd/pkg/validate"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
)

type HandlerInterface interface {
	CreateHandler(c *fiber.Ctx) error
	UpdateHandler(c *fiber.Ctx) error
	GetHandler(c *fiber.Ctx) error
	ListHandler(c *fiber.Ctx) error
	RegisterHandlers(app *fiber.App)
}

type HandlerStruct struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) HandlerInterface {
	return &HandlerStruct{service: service}
}

var _ HandlerInterface = &HandlerStruct{}

func (h *HandlerStruct) CreateHandler(c *fiber.Ctx) error {
	req := models.ApplicationCreateReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HttpResponseJSON{
			Error: err.Error(),
		})
	}

	err = validate.Validate.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HttpResponseJSON{
			Error: err.Error(),
		})
	}

	res, err := h.service.CreateService(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HttpResponseJSON{
			Error: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(models.HttpResponseJSON{
		Data: res,
	})
}

func (h *HandlerStruct) UpdateHandler(c *fiber.Ctx) error {
	req := models.ApplicationUpdateReq{}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HttpResponseJSON{
			Error: err.Error(),
		})
	}
	req.ID = cast.ToUint(c.Params("id"))
	err = validate.Validate.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HttpResponseJSON{
			Error: err.Error(),
		})
	}
	res, err := h.service.UpdateService(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HttpResponseJSON{
			Error: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(models.HttpResponseJSON{
		Data: res,
	})
}

func (h *HandlerStruct) GetHandler(c *fiber.Ctx) error {
	req := models.ApplicationGetReq{}
	req.ID = cast.ToUint(c.Params("id"))

	err := validate.Validate.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HttpResponseJSON{
			Error: err.Error(),
		})
	}

	res, err := h.service.GetService(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HttpResponseJSON{
			Error: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(models.HttpResponseJSON{
		Data: res.Data,
	})
}

func (h *HandlerStruct) ListHandler(c *fiber.Ctx) error {
	req := models.ApplicationListReq{}
	err := c.QueryParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HttpResponseJSON{
			Error: err.Error(),
		})
	}
	res, err := h.service.ListService(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.HttpResponseJSON{
			Error: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(models.HttpResponseJSON{
		Data: res.Data,
	})
}

func (h *HandlerStruct) RegisterHandlers(app *fiber.App) {
	app.Post("/applications", h.CreateHandler)
	app.Put("/applications/:id", h.UpdateHandler)
	app.Get("/applications/:id", h.GetHandler)
	app.Get("/applications", h.ListHandler)
}
