package mapping

import (
	"context"
	"fmt"
	"log"
	"orders/src/model"
	"orders/src/storage"

	"github.com/jackc/pgx/v4"
)

func GetProducts(ctx context.Context) (*model.QueryResult, error) {
	repo := storage.SQLRepo()
	rows, err := repo.Pool.Query(ctx, `SELECT UUID, DESCRIPTION, PRICE, LEFT_IN_STOCK,CREATE_TS FROM PRODUCTS`)
	if err != nil {
		log.Printf("Failed to query products %v", err)
		return nil, err
	}
	defer rows.Close()
	result := []*model.Product{}
	for rows.Next() {
		product := model.Product{}
		err := rows.Scan(&product.UUID, &product.Description, &product.Price, &product.LeftInStock, &product.CreatedAt)
		if err != nil {
			log.Printf("Failed to query product details %v", err)
			return nil, err
		}
		result = append(result, &product)
	}
	rows.Close()
	return model.NewOKResult(result), nil
}

func GetProduct(ctx context.Context, productID string) (*model.QueryResult, error) {
	repo := storage.SQLRepo()
	rows, err := repo.Pool.Query(ctx,
		`SELECT UUID, DESCRIPTION, PRICE, LEFT_IN_STOCK, CREATE_TS FROM PRODUCTS WHERE UUID = $1`, productID)
	if err != nil {
		log.Printf("Failed to query product %v", err)
		return nil, err
	}
	defer rows.Close()
	product := model.Product{
		UUID: productID,
	}
	for rows.Next() {
		err := rows.Scan(&product.UUID, &product.Description, &product.Price, &product.LeftInStock, &product.CreatedAt)
		if err != nil {
			log.Printf("Failed to query product details %v", err)
			return nil, err
		}
	}
	return model.NewOKResult(&product), nil
}

func UpdateProduct(ctx context.Context, productID string, count int) (*model.QueryResult, error) {
	repo := storage.SQLRepo()
	tx, err := repo.Pool.Begin(ctx)
	if err != nil {
		log.Printf("Failed to start transaction %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)
	cmd, err := tx.Exec(ctx, `UPDATE PRODUCTS SET
		LEFT_IN_STOCK = LEFT_IN_STOCK + $2,
		UPDATE_TS = NOW()
		WHERE UUID = $1
	`, productID, count)
	if err != nil {
		log.Printf("Failed to update product %v", err)
		return nil, err
	}
	if cmd.RowsAffected() == 0 {
		return nil, fmt.Errorf("product with id %v not found", productID)
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("Failed to commit transaction %v", err)
		return nil, err
	}
	return model.EmptyOKResult(), nil
}

func AddProduct(ctx context.Context, p model.Product) (*model.QueryResult, error) {
	repo := storage.SQLRepo()
	tx, err := repo.Pool.Begin(ctx)
	if err != nil {
		log.Printf("Failed to start transaction %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(ctx, `INSERT INTO PRODUCTS(
				DESCRIPTION, PRICE, LEFT_IN_STOCK, CREATE_TS) 
		VALUES($1, $2, $3, NOW()) RETURNING UUID`,
		p.Description, p.Price, p.LeftInStock)
	if err != nil {
		log.Printf("Failed to add product %v", err)
		return nil, err
	}
	defer rows.Close()
	var id *string
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Printf("Failed to return product id %v", err)
			return nil, err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("Failed to commit transaction %v", err)
		return nil, err
	}
	return model.NewOKResult(*id), nil
}

func DeleteProduct(ctx context.Context, productID string) (*model.QueryResult, error) {
	repo := storage.SQLRepo()
	tx, err := repo.Pool.Begin(ctx)
	if err != nil {
		log.Printf("Failed to start transaction %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)
	cmd, err := tx.Exec(ctx, `DELETE FROM PRODUCTS WHERE UUID = $1`, productID)
	if err != nil {
		log.Printf("Failed to delete product %v", err)
		return nil, err
	}
	if cmd.RowsAffected() == 0 {
		return nil, fmt.Errorf("product with id %v not found", productID)
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("Failed to commit transaction %v", err)
		return nil, err
	}
	return model.EmptyOKResult(), nil
}

func HasProductInStock(tx pgx.Tx, ctx context.Context, productID string, count int) (bool, error) {
	rows, err := tx.Query(ctx, `SELECT LEFT_IN_STOCK FROM PRODUCTS WHERE UUID = $1`, productID)
	if err != nil {
		log.Printf("Failed to delete product %v", err)
		return false, err
	}
	defer rows.Close()
	var countInStock int
	for rows.Next() {
		err := rows.Scan(&countInStock)
		if err != nil {
			log.Printf("Failed to query count of product %v", err)
			return false, err
		}
	}
	if count > countInStock {
		return false, nil
	} else {
		return true, nil
	}
}
