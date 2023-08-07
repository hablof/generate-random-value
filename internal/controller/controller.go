package controller

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/hablof/generate-random-value/internal/models"
	"github.com/hablof/generate-random-value/internal/repository"
	"github.com/hablof/generate-random-value/internal/service"
)

const (
	reqIDHeader = "request-id"
)

const (
	typeField    = "type"
	charsetField = "charset"
	lengthField  = "length"
)

const (
	retrieveIDField = "id"
)

type Generator interface {
	Generate(opts models.GenerateOptions) (string, error)
}

type Repository interface {
	Create(unit models.RandomValue) (uint64, error)
	ReadByReqID(reqID string) (models.RandomValue, error)
	ReadByValID(valID int) (string, error)
}

type Handler struct {
	r Repository
	g Generator
}

func NewServer(r Repository, g Generator) *fiber.App {

	h := Handler{
		r: r,
		g: g,
	}
	app := fiber.New()

	app.Use(logger.New(logger.ConfigDefault))
	app.Use(recover.New(recover.ConfigDefault))

	app.Get("/api/retrieve/", h.retrieve)
	app.Post("/api/generate/", h.generate)

	return app
}

// type requestBody struct {
// 	Type    string `json:"type"`
// 	Charset string `json:"charset"`
// 	Length  int    `json:"length"`
// }

func (h *Handler) retrieve(c *fiber.Ctx) error {

	idStr := c.Query(retrieveIDField)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		return h.returnErr(c, fiber.StatusBadRequest, "invalid id")
	}

	value, err := h.r.ReadByValID(id)
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return h.returnErr(c, fiber.StatusBadRequest, "id not found")

	case err != nil:
		log.Println(err)
		return h.returnErr(c, fiber.StatusInternalServerError, "internal server error")
	}

	unit := models.RandomValue{
		ID:    uint64(id),
		Value: value,
	}
	b, _ := json.Marshal(unit)

	return c.Status(fiber.StatusOK).Send(b)
}

func (h *Handler) generate(c *fiber.Ctx) error {

	// try get idempotent cached(?) value
	// P.S. this code smells bad
	// if found -- request body ignored
	reqID := c.Get(reqIDHeader)
	if reqID != "" {
		unit, err := h.r.ReadByReqID(reqID)
		switch {
		case errors.Is(err, repository.ErrNotFound):
			break

		case err != nil:
			log.Println(err)
			return h.returnErr(c, fiber.StatusInternalServerError, "internal server error")

		default:
			b, _ := json.Marshal(unit)
			return c.Status(fiber.StatusOK).Send(b)
		}
	}

	// reading json

	req := make(map[string]json.RawMessage, 3)
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return h.returnErr(c, fiber.StatusBadRequest, "bad request body")
	}

	// configuration generation options
	opts := models.GenerateOptions{}

	if rawType, ok := req[typeField]; ok {
		var genType string
		if err := json.Unmarshal(rawType, &genType); err != nil {
			return h.returnErr(c, fiber.StatusBadRequest, "invalid type")
		}
		opts.GenerationType = genType
	}

	if rawCharset, ok := req[charsetField]; ok {
		var charset string
		if err := json.Unmarshal(rawCharset, &charset); err != nil {
			return h.returnErr(c, fiber.StatusBadRequest, "invalid charset")
		}
		opts.SpecifyCharset(charset)
	}

	if length, ok := req[lengthField]; ok {
		i, err := strconv.Atoi(string(length))
		if err != nil {
			return h.returnErr(c, fiber.StatusBadRequest, "invalid length")
		}

		opts.SpecifyLength(i)
	}

	// generate random value
	randomValue, err := h.g.Generate(opts)
	switch {
	case errors.Is(err, service.ErrInvalidCharset):
		return h.returnErr(c, fiber.StatusBadRequest, "invalid charset")

	case errors.Is(err, service.ErrInvalidLength):
		return h.returnErr(c, fiber.StatusBadRequest, "invalid length")

	case errors.Is(err, service.ErrInvalidType):
		return h.returnErr(c, fiber.StatusBadRequest, "invalid type")

	case err != nil: // unpredicted err
		log.Println(err)
		return h.returnErr(c, fiber.StatusInternalServerError, "internal server error")
	}

	// save generated val
	unit := models.RandomValue{Value: randomValue}
	if reqID != "" {
		unit.RequestID.IsValid = true
		unit.RequestID.S = reqID
	}

	id, err := h.r.Create(unit)
	if err != nil { // unpredicted err
		log.Println(err)
		return h.returnErr(c, fiber.StatusInternalServerError, "internal server error")
	}

	unit.ID = id
	b, _ := json.Marshal(unit)

	return c.Status(fiber.StatusOK).Send(b)
}

func (*Handler) returnErr(c *fiber.Ctx, status int, msg string) error {
	b, _ := json.Marshal(map[string]string{"error": msg})
	return c.Status(status).Send(b)
}
