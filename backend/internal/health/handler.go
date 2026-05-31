package health

import (
	"github.com/gofiber/fiber/v2"
	"github.com/narinsak-u/backend/pkg/utils"
)

// Check returns a simple health check response. Used by container
// orchestrators (Docker healthcheck, Kubernetes liveness probe) to
// verify the server is alive and serving requests.
func Check(c *fiber.Ctx) error {
	return utils.Success(c, fiber.Map{"status": "ok"})
}
