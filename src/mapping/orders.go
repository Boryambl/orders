package mapping

import (
	"context"
	"fmt"
	"log"
	"orders/src/model"
	"orders/src/storage"

	"github.com/jackc/pgx/v4"
)

func GetOrders(ctx context.Context) (*model.QueryResult, error) {
	repo := storage.SQLRepo()
	rows, err := repo.Pool.Query(ctx, `select number, description, fio, count  from orders 
	join product_order on orders.uuid = product_order.order_id 
	join products on products.uuid = product_id
	join user_order on orders.uuid = user_order.order_id
	join users on users.uuid = user_id
	order by orders.create_ts desc`)
	if err != nil {
		log.Printf("Failed to query orders %v", err)
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]map[int][]model.ProductToShow)
	for rows.Next() {
		var number, count int
		var fio, desc string
		err := rows.Scan(&number, &desc, &fio, &count)
		if err != nil {
			log.Printf("Failed to query order details %v", err)
			return nil, err
		}
		orders, ok := result[fio]
		if ok {
			_, okok := orders[number]
			if okok {
				orders[number] = append(orders[number], model.ProductToShow{Product: desc, Count: count})
				result[fio] = orders
			} else {
				result[fio][number] = append(result[fio][number], model.ProductToShow{Product: desc, Count: count})
			}
		} else {
			result[fio] = make(map[int][]model.ProductToShow)
			result[fio][number] = append(result[fio][number], model.ProductToShow{Product: desc, Count: count})
		}
	}
	rows.Close()
	return model.NewOKResult(result), nil
}

func GetUsersOrders(ctx context.Context, userID string) (*model.QueryResult, error) {
	repo := storage.SQLRepo()
	rows, err := repo.Pool.Query(ctx,
		`select number, description, fio, count  from orders 
		join product_order on orders.uuid = product_order.order_id 
		join products on products.uuid = product_id
		join user_order on orders.uuid = user_order.order_id
		join users on users.uuid = user_id
		where users.uuid = $1
		order by orders.create_ts desc`, userID)
	if err != nil {
		log.Printf("Failed to query product %v", err)
		return nil, err
	}
	defer rows.Close()
	result := make(map[int][]model.ProductToShow)
	var fio string
	for rows.Next() {
		var number, count int
		var desc string
		err := rows.Scan(&number, &desc, &fio, &count)
		if err != nil {
			log.Printf("Failed to query product details %v", err)
			return nil, err
		}
		result[number] = append(result[number], model.ProductToShow{Product: desc, Count: count})
	}
	rows.Close()
	return model.NewOKResult(result), nil
}

func UpdateOrder(ctx context.Context, orderID string, o model.Order) (*model.QueryResult, error) {
	repo := storage.SQLRepo()
	tx, err := repo.Pool.Begin(ctx)
	if err != nil {
		log.Printf("Failed to start transaction %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)
	cmd, err := tx.Exec(ctx, `UPDATE ORDERS SET
		UPDATE_TS = NOW()
		WHERE UUID = $1
	`, orderID)
	if err != nil {
		log.Printf("Failed to update order %v", err)
		return nil, err
	}
	if cmd.RowsAffected() == 0 {
		return nil, fmt.Errorf("order with id %v not found", orderID)
	}
	batch := &pgx.Batch{}
	for _, v := range o.Products {
		has, err := HasProductInStock(tx, ctx, v.ProductID, v.Count)
		if err != nil {
			log.Printf("Failed to check availability of products in stock %v", err)
			return nil, err
		}
		if !has {
			return nil, fmt.Errorf("product with id %v not found in stock in the right quantity", v.ProductID)
		}
		batch.Queue(`INSERT INTO PRODUCT_ORDER (ORDER_ID, PRODUCT_ID, COUNT) VALUES ($1, $2, $3) ON CONFLICT (ORDER_ID, PRODUCT_ID) DO UPDATE SET COUNT = PRODUCT_ORDER.COUNT + $3`, orderID, v.ProductID, v.Count)
		batch.Queue(`UPDATE PRODUCTS SET LEFT_IN_STOCK = LEFT_IN_STOCK - $2 WHERE UUID = $1`, v.ProductID, v.Count)
	}
	rs := tx.SendBatch(ctx, batch)
	defer rs.Close()
	for i := 0; i < batch.Len(); i++ {
		_, err := rs.Exec()
		if err != nil {
			return nil, err
		}
	}
	rs.Close()
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("Failed to commit transaction %v", err)
		return nil, err
	}
	return model.EmptyOKResult(), nil
}

func AddOrder(ctx context.Context, o model.Order) (*model.QueryResult, error) {
	repo := storage.SQLRepo()
	tx, err := repo.Pool.Begin(ctx)
	if err != nil {
		log.Printf("Failed to start transaction %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(ctx, `INSERT INTO ORDERS(CREATE_TS) VALUES(NOW()) RETURNING UUID`)
	if err != nil {
		log.Printf("Failed to add order %v", err)
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
	batch := &pgx.Batch{}
	for _, v := range o.Products {
		has, err := HasProductInStock(tx, ctx, v.ProductID, v.Count)
		if err != nil {
			log.Printf("Failed to check availability of products in stock %v", err)
			return nil, err
		}
		if !has {
			return nil, fmt.Errorf("product with id %v not found in stock in the right quantity", v.ProductID)
		}
		batch.Queue(`INSERT INTO PRODUCT_ORDER (ORDER_ID, PRODUCT_ID, COUNT) VALUES ($1, $2, $3)`, id, v.ProductID, v.Count)
		batch.Queue(`UPDATE PRODUCTS SET LEFT_IN_STOCK = LEFT_IN_STOCK - $2 WHERE UUID = $1`, v.ProductID, v.Count)
	}
	batch.Queue(`INSERT INTO USER_ORDER (ORDER_ID, USER_ID) VALUES ($1, $2) ON CONFLICT DO NOTHING`, id, o.UserID)
	rs := tx.SendBatch(ctx, batch)
	defer rs.Close()
	for i := 0; i < batch.Len(); i++ {
		_, err := rs.Exec()
		if err != nil {
			return nil, err
		}
	}
	rs.Close()
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("Failed to commit transaction %v", err)
		return nil, err
	}
	return model.NewOKResult(*id), nil
}

func DeleteOrder(ctx context.Context, orderID string) (*model.QueryResult, error) {
	repo := storage.SQLRepo()
	tx, err := repo.Pool.Begin(ctx)
	if err != nil {
		log.Printf("Failed to start transaction %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)
	cmd, err := tx.Exec(ctx, `DELETE FROM ORDERS WHERE UUID = $1`, orderID)
	if err != nil {
		log.Printf("Failed to delete order %v", err)
		return nil, err
	}
	if cmd.RowsAffected() == 0 {
		return nil, fmt.Errorf("order with id %v not found", orderID)
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("Failed to commit transaction %v", err)
		return nil, err
	}
	return model.EmptyOKResult(), nil
}
