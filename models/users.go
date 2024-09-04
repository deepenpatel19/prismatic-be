package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deepenpatel19/prismatic-be/logger"
	pgx "github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type UserCreateSchema struct {
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
	Email     string `json:"email" form:"email"`
	Password  string `json:"password" form:"password"`
}

type UserSchema struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserResponseSchema struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (data UserCreateSchema) Insert(uuidString string) (int64, error) {
	query := `INSERT INTO 
				users
					(first_name, last_name, email, password)
				VALUES
					($1, $2, $3, $4)
				RETURNING id`
	queryToExecute := QueryStructToExecute{Query: query}
	id, err := queryToExecute.InsertOrUpdateOperations(uuidString, data.FirstName, data.LastName, data.Email, data.Password)
	return id, err
}

func (data UserCreateSchema) Update(uuidString string, id int64) (int64, error) {
	query := `UPDATE 
				users
					SET first_name=$1, last_name=$2
				WHERE id=$3
				RETURNING id`
	queryToExecute := QueryStructToExecute{Query: query}
	id, err := queryToExecute.InsertOrUpdateOperations(uuidString, data.FirstName, data.LastName, id)
	return id, err

}

func DeleteUserFromDB(uuidString string, userId int64) (bool, error) {
	query := fmt.Sprintf(`DELETE FROM users WHERE id=%d`, userId)
	queryToExecute := QueryStructToExecute{Query: query}
	status, err := queryToExecute.DeleteOperation(uuidString)
	return status, err
}

func FetchUserForAuth(email string) UserSchema {
	logger.Logger.Info("MODELS :: Will fetch user details for auth", zap.String("email", email))

	var userData UserSchema
	dbConnection := DbPool()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	// ctx := context.Background()

	tx, err := dbConnection.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		logger.Logger.Error("MODELS :: Error while begin transaction", zap.Error(err))
		return userData
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	query := fmt.Sprintf(`SELECT
							u.id,
							u.first_name, 
							u.last_name, 
							u.email, 
							u.password
							FROM users u
							WHERE u.email='%s' LIMIT 1`, email)
	logger.Logger.Info("MODELS :: Query", zap.String("query", query))
	err = tx.QueryRow(ctx, query).Scan(
		&userData.Id,
		&userData.FirstName,
		&userData.LastName,
		&userData.Email,
		&userData.Password,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Logger.Info("MODELS :: Query - No rows found. ", zap.String("query", query))
			return userData
		}
		logger.Logger.Error("MODELS :: Error while executing query.",
			zap.Error(err),
		)
		return userData
	}
	return userData
}

func FetchUserForAuthV1(email string) UserSchema {
	logger.Logger.Info("MODELS :: Will fetch user details for auth v1", zap.String("email", email))

	var userData UserSchema
	query := fmt.Sprintf(`SELECT
							u.id,
							u.first_name, 
							u.last_name, 
							u.email, 
							u.password
							FROM users u
							WHERE u.email='%s' LIMIT 1`, email)
	logger.Logger.Info("MODELS :: Query", zap.String("query", query))

	queryToExecute := QueryStructToExecute{Query: query}
	rows, _, err := queryToExecute.FetchRows("")
	logger.Logger.Info("user models : ", zap.Any("rows ", rows), zap.Error(err))
	if err != nil {
		return userData
	}

	if len(rows) == 0 {
		return userData
	}

	singleRow := rows[0]

	jsonbody, err := json.Marshal(singleRow)
	if err != nil {
		// do error check
		fmt.Println(err)
	}

	if err := json.Unmarshal(jsonbody, &userData); err != nil {
		// do error check
		fmt.Println(err)
	}

	// for rows.Next() {
	// 	err = rows.Scan(
	// 		&userData.Id,
	// 		&userData.FirstName,
	// 		&userData.LastName,
	// 		&userData.Email,
	// 		&userData.Password,
	// 		&userData.Type,
	// 	)
	// 	if err != nil {
	// 		logger.Logger.Error("MODELS :: Error while iterating rows", zap.String("requestId", ""), zap.Error(err))
	// 		return userData
	// 	}
	// }
	logger.Logger.Info("MODELS :: before sending dat a", zap.Any("data", userData), zap.Any("id", userData.Id))
	return userData
}

func FetchUserForMeV1(uuidString string, id int64) UserResponseSchema {
	logger.Logger.Info("MODELS :: Will fetch user details for me v1", zap.Int64("userId", id), zap.String("requestId", uuidString))

	var userData UserResponseSchema
	query := fmt.Sprintf(`SELECT
							u.id,
							u.first_name, 
							u.last_name, 
							u.email
							FROM users u
							WHERE u.id=%d LIMIT 1`, id)
	logger.Logger.Info("MODELS :: Query", zap.String("query", query), zap.String("requestId", uuidString))

	queryToExecute := QueryStructToExecute{Query: query}
	rows, _, err := queryToExecute.FetchRows("")
	if err != nil {
		return userData
	}

	if len(rows) == 0 {
		return userData
	}

	singleRow := rows[0]

	jsonbody, err := json.Marshal(singleRow)
	if err != nil {
		// do error check
		fmt.Println(err)
	}

	if err := json.Unmarshal(jsonbody, &userData); err != nil {
		// do error check
		fmt.Println(err)
	}

	// for k, v := range rows {
	// 	err := SetField(s, k, v)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// for rows.Next() {
	// 	err = rows.Scan(
	// 		&userData.Id,
	// 		&userData.FirstName,
	// 		&userData.LastName,
	// 		&userData.Email,
	// 		&userData.Password,
	// 		&userData.Type,
	// 	)
	// 	if err != nil {
	// 		logger.Logger.Error("MODELS :: Error while iterating rows", zap.String("requestId", uuidString), zap.Error(err))
	// 		return userData
	// 	}
	// }

	return userData
}

func FetchAllUsers(uuidString string) ([]UserResponseSchema, int64, error) {
	logger.Logger.Info("MODELS :: Will fetch all users", zap.String("requestId", uuidString))

	var userData []UserResponseSchema
	var count int64
	var err error
	var rows []map[string]any
	query := `SELECT 
				id,
				first_name,
				last_name,
				email,
				count(*) over()
			  FROM users`

	queryToExecute := QueryStructToExecute{Query: query}
	rows, count, err = queryToExecute.FetchRows(uuidString)
	if err != nil {
		return userData, count, err
	}

	for _, data := range rows {
		var singleUserData UserResponseSchema
		jsonbody, err := json.Marshal(data)
		if err != nil {
			// do error check
			fmt.Println(err)
		}

		if err := json.Unmarshal(jsonbody, &singleUserData); err != nil {
			// do error check
			fmt.Println(err)
		}

		userData = append(userData, singleUserData)
	}

	logger.Logger.Info("MODELS :: Rows ", zap.Any("data", rows))
	return userData, count, nil
}

func InsertUserLoginHistory(uuidString string, userId int64, ipaddress string) (int64, error) {
	query := `INSERT INTO 
				user_login_history
					(user_id, login_at, ip_address)
			VALUES
				($1, $2, $3)
			RETURNING id`
	queryToExecute := QueryStructToExecute{Query: query}
	id, err := queryToExecute.InsertOrUpdateOperations(uuidString, userId, time.Now().UTC(), ipaddress)
	return id, err
}
