package server

import (
	"log"
	"orders/src/mapping"
	"orders/src/model"

	"github.com/gofiber/fiber/v2"
)

func getUsers(ctx *fiber.Ctx) error {
	var err error
	id := ctx.Query("id")
	var result *model.QueryResult
	if len(id) == 0 {
		result, err = mapping.GetUsers(ctx.Context())
	} else {
		result, err = mapping.GetUser(ctx.Context(), id)
	}
	if err != nil {
		log.Print(err)
		return err
	}
	return ctx.JSON(result)
}

func addUser(ctx *fiber.Ctx) error {
	user := new(model.User)
	if err := ctx.BodyParser(user); err != nil {
		log.Print(err)
		return err
	}
	result, err := mapping.AddUser(ctx.Context(), *user)
	if err != nil {
		log.Print(err)
		return err
	}
	return ctx.JSON(result)
}

func deleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")
	result, err := mapping.DeleteUser(ctx.Context(), id)
	if err != nil {
		log.Print(err)
		return err
	}
	return ctx.JSON(result)
}

func updateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")
	user := new(model.User)
	if err := ctx.BodyParser(user); err != nil {
		log.Print(err)
		return err
	}
	result, err := mapping.UpdateUser(ctx.Context(), id, *user)
	if err != nil {
		return err
	}
	return ctx.JSON(result)
}
