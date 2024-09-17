package rest

import "github.com/gofiber/fiber/v2"

type ErrorResponse struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id"`
}

func HandleError(c *fiber.Ctx, statusCode int, err error) error {
	return c.Status(statusCode).JSON(ErrorResponse{
		Error:     err.Error(),
		RequestID: c.GetRespHeader("x-request-id"),
	})
}
