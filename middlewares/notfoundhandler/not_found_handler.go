package notfoundhandler

import "github.com/gofiber/fiber/v2"

func New() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		accept := c.Get("Accept")
		if accept == "application/json" {
			return c.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"success": false,
					"data": fiber.Map{
						"msg": "Route Not Found"},
				})
		}
		if accept != "application/json" {
			return c.Status(fiber.StatusNotFound).SendString(
				"Route Not Found")
		}
		return c.Next()
	}
}
