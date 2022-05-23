package server

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Start(port uint) error {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		StrictRouting:         true,
		CaseSensitive:         true,
		BodyLimit:             1000 * 1024 * 1024,
	})
	app.Get("/api/users", getUsers)
	app.Post("/api/users", addUser)
	app.Put("/api/users/:id", updateUser)
	app.Delete("/api/users/:id", deleteUser)
	app.Get("/api/products", getProducts)
	app.Post("/api/products", addProduct)
	app.Put("/api/products/:id", updateProduct)
	app.Delete("/api/products/:id", deleteProduct)
	app.Get("/api/orders", getOrders)
	app.Post("/api/orders", addOrder)
	app.Put("/api/orders/:id", updateOrder)
	app.Delete("/api/orders/:id", deleteOrder)
	log.Printf("REST API server running on %v", port)
	return app.Listen(fmt.Sprintf("localhost:%d", port))
}
