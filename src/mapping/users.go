package mapping

import (
	"context"
	"fmt"
	"log"
	"orders/src/model"
	"orders/src/storage"
	"time"
)

func GetUsers(ctx context.Context) (*model.QueryResult, error) {
	repo := storage.SQLRepo()
	rows, err := repo.Pool.Query(ctx, `SELECT UUID, FIRSTNAME, SURNAME, MIDDLENAME, FIO, SEX, AGE, CREATE_TS FROM USERS `)
	if err != nil {
		log.Printf("Failed to query users %v", err)
		return nil, err
	}
	defer rows.Close()
	result := []*model.User{}
	for rows.Next() {
		user := model.User{}
		var fname, sname, mname *string
		var age *int
		var createdTS *time.Time
		err := rows.Scan(&user.UUID, &fname, &sname, &mname, &user.FIO, &user.Sex, &age, &createdTS)
		if err != nil {
			log.Printf("Failed to query user details %v", err)
			return nil, err
		}
		if age != nil {
			user.Age = *age
		}
		if fname != nil {
			user.Firstname = *fname
		}
		if sname != nil {
			user.Surname = *sname
		}
		if mname != nil {
			user.Middlename = *mname
		}
		if createdTS != nil {
			user.CreatedAt = *createdTS
		}
		result = append(result, &user)
	}
	rows.Close()
	return model.NewOKResult(result), nil
}

func GetUser(ctx context.Context, userID string) (*model.QueryResult, error) {
	repo := storage.SQLRepo()
	rows, err := repo.Pool.Query(ctx,
		`SELECT UUID, FIRSTNAME, SURNAME, MIDDLENAME, FIO, SEX, AGE, CREATE_TS FROM USERS WHERE UUID = $1`, userID)
	if err != nil {
		log.Printf("Failed to query user %v", err)
		return nil, err
	}
	defer rows.Close()
	user := model.User{
		UUID: userID,
	}
	for rows.Next() {
		var fname, sname, mname *string
		var age *int
		var createdTS *time.Time
		err := rows.Scan(&user.UUID, &fname, &sname, &mname, &user.FIO, &user.Sex, &age, &createdTS)
		if err != nil {
			log.Printf("Failed to query user details %v", err)
			return nil, err
		}
		if age != nil {
			user.Age = *age
		}
		if fname != nil {
			user.Firstname = *fname
		}
		if sname != nil {
			user.Surname = *sname
		}
		if mname != nil {
			user.Middlename = *mname
		}
		if createdTS != nil {
			user.CreatedAt = *createdTS
		}
	}
	return model.NewOKResult(&user), nil
}

func UpdateUser(ctx context.Context, userID string, user model.User) (*model.QueryResult, error) {
	repo := storage.SQLRepo()
	tx, err := repo.Pool.Begin(ctx)
	if err != nil {
		log.Printf("Failed to start transaction %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)
	cmd, err := tx.Exec(ctx, `UPDATE USERS SET
		FIRSTNAME = $2,
		SURNAME = $3,
		MIDDLENAME = $4,
		SEX = $5,
		AGE = $6,
		UPDATE_TS = NOW()
		WHERE UUID = $1
	`, userID, user.Firstname, user.Surname, user.Middlename, user.Sex, user.Age)
	if err != nil {
		log.Printf("Failed to update user %v", err)
		return nil, err
	}
	if cmd.RowsAffected() == 0 {
		return nil, fmt.Errorf("user with id %v not found", userID)
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("Failed to commit transaction %v", err)
		return nil, err
	}
	return model.EmptyOKResult(), nil
}

func AddUser(ctx context.Context, u model.User) (*model.QueryResult, error) {
	repo := storage.SQLRepo()
	tx, err := repo.Pool.Begin(ctx)
	if err != nil {
		log.Printf("Failed to start transaction %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(ctx, `INSERT INTO USERS(
				FIRSTNAME, SURNAME, MIDDLENAME, SEX, AGE, CREATE_TS) 
		VALUES($1, $2, $3, $4, $5, NOW()) RETURNING UUID`,
		u.Firstname,
		u.Surname,
		u.Middlename,
		u.Sex,
		u.Age)
	if err != nil {
		log.Printf("Failed to add user %v", err)
		return nil, err
	}
	defer rows.Close()
	var id *string
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Printf("Failed to return user id %v", err)
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

func DeleteUser(ctx context.Context, userID string) (*model.QueryResult, error) {
	repo := storage.SQLRepo()
	tx, err := repo.Pool.Begin(ctx)
	if err != nil {
		log.Printf("Failed to start transaction %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)
	cmd, err := tx.Exec(ctx, `DELETE FROM USERS WHERE UUID = $1`, userID)
	if err != nil {
		log.Printf("Failed to delete user %v", err)
		return nil, err
	}
	if cmd.RowsAffected() == 0 {
		return nil, fmt.Errorf("user with id %v not found", userID)
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("Failed to commit transaction %v", err)
		return nil, err
	}
	return model.EmptyOKResult(), nil
}
