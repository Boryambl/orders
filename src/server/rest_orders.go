package server

import (
	"log"
	"orders/src/mapping"
	"orders/src/model"

	"github.com/gofiber/fiber/v2"
)

func getOrders(ctx *fiber.Ctx) error {
	var err error
	id := ctx.Query("user_id")
	var result *model.QueryResult
	if len(id) == 0 {
		result, err = mapping.GetOrders(ctx.Context())
	} else {
		result, err = mapping.GetUsersOrders(ctx.Context(), id)
	}
	if err != nil {
		log.Print(err)
		return err
	}
	return ctx.JSON(result)
}

func addOrder(ctx *fiber.Ctx) error {
	order := &model.Order{}
	if err := ctx.BodyParser(&order); err != nil {
		log.Print(err)
		return err
	}
	result, err := mapping.AddOrder(ctx.Context(), *order)
	if err != nil {
		log.Print(err)
		return err
	}
	return ctx.JSON(result)
}

func deleteOrder(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")
	result, err := mapping.DeleteOrder(ctx.Context(), id)
	if err != nil {
		log.Print(err)
		return err
	}
	return ctx.JSON(result)
}

func updateOrder(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")
	order := &model.Order{}
	if err := ctx.BodyParser(&order); err != nil {
		log.Print(err)
		return err
	}
	result, err := mapping.UpdateOrder(ctx.Context(), id, *order)
	if err != nil {
		return err
	}
	return ctx.JSON(result)
}
