// The file resolves the shortened url to actual url
package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/mihir-chhatre/go-short-url/database"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	//calls the func in dataabse.go
	r := database.CreateClient(0)
	defer r.Close()

	//get the value corresponding to to key - url
	value, err := r.Get(database.Ctx, url).Result()

	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot connect to DB"})
	}

	//If everything went well ->

	rInr := database.CreateClient(1) //incr counter to 1
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")

	return c.Redirect(value, 301) //redirecting to actual url found in db with status 301
}
