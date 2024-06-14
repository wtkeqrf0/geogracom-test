package route

import (
	"geogracom-test/pkg/rest"
	"github.com/gofiber/fiber/v2"
)

func (m *method) MapHandlers(app fiber.Router) {

	app.Post("register", rest.API[Upsert](m.Create).Handler(func(c *fiber.Ctx, out any) error {
		return c.BodyParser(out)
	}))
	app.Get(":id", rest.API[Get](m.Get).Handler(func(c *fiber.Ctx, out any) error {
		return c.ParamsParser(out)
	}))
	app.Delete("", rest.API[Delete](m.Delete).Handler(func(c *fiber.Ctx, out any) error {
		return c.BodyParser(out)
	}))
}
