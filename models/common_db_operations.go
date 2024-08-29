package models

import (
	"context"
	"fmt"
	"time"

	"github.com/deepenpatel19/prismatic-be/core"
	"github.com/deepenpatel19/prismatic-be/logger"
	pgx "github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type QueryStructToExecute struct {
	Query     string
	QueryList []string
}

func (query QueryStructToExecute) InsertOrUpdateOperations(uuidString string, args ...any) (int64, error) {
	logger.Logger.Info("MODELS :: Will do insert operations", zap.String("requestId", uuidString), zap.String("query", query.Query), zap.Any("args", args))

	var id int64
	dbConnection := DbPool()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.Config.DBQueryTimeout)*time.Second)
	defer cancel()

	tx, err := dbConnection.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadWrite})
	if err != nil {
		logger.Logger.Error("MODELS :: Error while begin transaction", zap.Error(err), zap.String("requestId", uuidString))
		return id, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	err = tx.QueryRow(ctx, query.Query, args...).Scan(&id)
	if err != nil {
		logger.Logger.Error("MODELS :: Error while executing query.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		return id, err
	}

	return id, nil
}

func (query QueryStructToExecute) DeleteOperation(uuidString string) (bool, error) {
	logger.Logger.Info("MODELS :: Will do delete operation.", zap.String("requestId", uuidString), zap.String("query", query.Query))
	dbConnection := DbPool()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.Config.DBQueryTimeout)*time.Second)
	defer cancel()

	tx, err := dbConnection.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadWrite})
	if err != nil {
		logger.Logger.Error("MODELS :: Error while begin transaction", zap.Error(err), zap.String("requestId", uuidString))
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	tx.Exec(ctx, query.Query)
	return true, nil
}

func (query QueryStructToExecute) InsertOrUpdateMultipleQueries(uuidString string) (int64, error) {
	logger.Logger.Info("MODELS :: Will do insert/update multiple operations", zap.String("requestId", uuidString))

	var id int64
	dbConnection := DbPool()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.Config.DBQueryTimeout)*time.Second)
	defer cancel()

	tx, err := dbConnection.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadWrite})
	if err != nil {
		logger.Logger.Error("MODELS :: Error while begin transaction", zap.Error(err), zap.String("requestId", uuidString))
		return id, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	for _, singleQuery := range query.QueryList {
		err = tx.QueryRow(ctx, singleQuery).Scan(&id)
		if err != nil {
			logger.Logger.Error("MODELS :: Error while executing query.",
				zap.String("requestId", uuidString),
				zap.Error(err),
			)
			return id, err
		}
	}

	return id, nil
}

func (query QueryStructToExecute) FetchRows(uuidString string) ([]map[string]any, int64, error) {

	var data map[string]any
	var queryData []map[string]any
	var count int64
	// var intConversion bool
	dbConnection := DbPool()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(core.Config.DBQueryTimeout)*time.Second)
	defer cancel()

	tx, err := dbConnection.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadWrite})
	if err != nil {
		logger.Logger.Error("MODELS :: Error while begin transaction", zap.Error(err), zap.String("requestId", uuidString))
		return nil, count, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	rows, err := tx.Query(ctx, query.Query)
	if err != nil {
		logger.Logger.Error("MODELS :: Error while executing query.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		return queryData, count, err
	}

	defer rows.Close()
	logger.Logger.Info("MODELS :: Rows fetched ", zap.Any("rows", rows), zap.Error(err))

	for rows.Next() {
		// values, err := rows.Values()
		// logger.Logger.Info("Row values ", zap.Any("rows", values), zap.Error(err))

		data, err = pgx.RowToMap(rows)

		if x, found := data["count"]; found {
			fmt.Println("before conversion ", x, found)
			count = x.(int64)
			// count1, intConversion := x.(int)
			// fmt.Println("count ", count1, intConversion)
			fmt.Printf("count type %T => %d", x, count)
		}

		logger.Logger.Info("MODELS :: map conversion ", zap.Any("map", data), zap.Error(err))
		queryData = append(queryData, data)
	}

	return queryData, count, nil
}
