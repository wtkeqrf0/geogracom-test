package api

import (
	"context"
	"geogracom-test/pkg/kernel"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"net"
	"time"
)

type Server struct {
	App *fiber.App
	ls  *net.ListenConfig
}

func New() *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: handleError,
	})

	app.Use(
		// query logger
		logger.New(logger.Config{
			TimeFormat: time.DateTime,
			Format:     "${time} -  ${method} ${path} - ${status} - ${ip}\n",
			TimeZone:   "Europe/Moscow",
		}),

		// panic recover
		recover.New(recover.Config{
			EnableStackTrace: true,
		}),
	)

	return &Server{
		App: app,
		ls:  new(net.ListenConfig),
	}
}

func (s *Server) Start(ctx context.Context, address string) {
	lis, err := s.ls.Listen(ctx, "tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server started on %s", address)

	context.AfterFunc(ctx, s.shutdown(time.Second*3))
	if err = s.App.Listener(lis); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func handleError(c *fiber.Ctx, err error) error {
	_ = c.Status(fiber.StatusInternalServerError).JSON(kernel.Response{Description: err.Error()})
	return nil
}

func (s *Server) shutdown(timeout time.Duration) func() {
	return func() {
		if err := s.App.ShutdownWithTimeout(timeout); err != nil {
			log.Printf("failed to shutdown server: %v", err)
		}
	}
}
