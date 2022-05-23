package server

import (
	"log"
	"orders/src/mapping"
	"orders/src/model"

	"github.com/gofiber/fiber/v2"
)

func getProducts(ctx *fiber.Ctx) error {
	var err error
	id := ctx.Query("id")
	var result *model.QueryResult
	if len(id) == 0 {
		result, err = mapping.GetProducts(ctx.Context())
	} else {
		result, err = mapping.GetProduct(ctx.Context(), id)
	}
	if err != nil {
		log.Print(err)
		return err
	}
	return ctx.JSON(result)
}

func addProduct(ctx *fiber.Ctx) error {
	product := new(model.Product)
	if err := ctx.BodyParser(product); err != nil {
		log.Print(err)
		return err
	}
	result, err := mapping.AddProduct(ctx.Context(), *product)
	if err != nil {
		log.Print(err)
		return err
	}
	return ctx.JSON(result)
}

func deleteProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")
	result, err := mapping.DeleteProduct(ctx.Context(), id)
	if err != nil {
		log.Print(err)
		return err
	}
	return ctx.JSON(result)
}

func updateProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id", "")
	var count struct {
		Count int `json:"count"`
	}
	if err := ctx.BodyParser(&count); err != nil {
		log.Print(err)
		return err
	}
	result, err := mapping.UpdateProduct(ctx.Context(), id, count.Count)
	if err != nil {
		return err
	}
	return ctx.JSON(result)
}
