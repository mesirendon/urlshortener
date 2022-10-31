package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/mesirendon/urlshortener/database"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	r := database.CreateClient(0)
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	value, err := r.Get(database.Ctx, url).Result()
	switch {
	case err == redis.Nil:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Short not found in the database",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot connect to the DB",
		})
	case value == "":
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Value empty",
		})
	}

	rInr := database.CreateClient(1)
	defer func() {
		if err := rInr.Close(); err != nil {
			panic(err)
		}
	}()

	_ = rInr.Incr(database.Ctx, "counter")

	return c.Redirect(value, 301)
}
