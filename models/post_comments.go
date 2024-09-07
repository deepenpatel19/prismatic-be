package models

import (
	"context"
	"fmt"
	"time"

	"github.com/deepenpatel19/prismatic-be/logger"
	pgx "github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type PostComment struct {
	PostId  int64  `json:"post_id" form:"post_id"`
	UserId  int64  `json:"user_id" form:"user_id"`
	Comment string `json:"comment" form:"comment" binding:"required"`
}

type PostCommentResponse struct {
	Id        int64  `json:"id" form:"id"`
	UserId    int64  `json:"user_id" form:"user_id"`
	PostId    int64  `json:"post_id" form:"post_id"`
	Comment   string `json:"comment" form:"comment" binding:"required"`
	CreatedAt string `json:"created_at" form:"created_at"`
}

func (data PostComment) Insert(uuidString string) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO
							post_comments
								(user_id, post_id, comment, created_at)
							VALUES
								(%d, %d, '%s', '%s')
							RETURNING id`,
		data.UserId,
		data.PostId,
		data.Comment,
		time.Now().UTC().Format(time.RFC3339),
	)

	queryToExecute := QueryStructToExecute{Query: query}
	id, err := queryToExecute.InsertOrUpdateOperations(uuidString)
	return id, err
}

func (data PostComment) Update(uuidString string, Id int64) (int64, error) {
	query := fmt.Sprintf(`UPDATE 
								post_comments
							SET comment='%s'
							WHERE id=%d AND user_id=%d
							RETURNING id`,
		data.Comment,
		Id,
		data.UserId,
	)
	queryToExecute := QueryStructToExecute{Query: query}
	id, err := queryToExecute.InsertOrUpdateOperations(uuidString)
	return id, err
}

func DeletePostComment(uuidString string, Id int64, userId int64, postId int64) (bool, error) {
	query := fmt.Sprintf(`DELETE 
							FROM post_comments
						  WHERE id=%d AND user_id=%d AND post_id=%d`, Id, userId, postId)
	queryToExecute := QueryStructToExecute{Query: query}
	status, err := queryToExecute.DeleteOperation(uuidString)
	return status, err
}

func FetchPostComments(uuidString string, limit int, offset int) ([]PostCommentResponse, int64, error) {
	var count int64
	var err error
	var postCommentData []PostCommentResponse

	dbConnection := DbPool()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tx, err := dbConnection.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
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

	query := fmt.Sprintf(`SELECT 
							id,
							user_id,
							post_id,
							comment,
							created_at,
							COUNT(*) OVER() as count
						  FROM
						  	post_comments
						  ORDER BY id DESC
						  LIMIT %d OFFSET %d`, limit, offset)

	rows, err := tx.Query(ctx, query)
	if err != nil {
		logger.Logger.Error("MODELS :: Error while executing query.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		return postCommentData, count, err
	}

	defer rows.Close()
	logger.Logger.Info("MODELS :: Rows fetched ", zap.Any("rows", rows), zap.Error(err))

	for rows.Next() {

		var singleData PostCommentResponse

		err = rows.Scan(
			&singleData.Id,
			&singleData.UserId,
			&singleData.PostId,
			&singleData.Comment,
			&singleData.CreatedAt,
			&count,
		)
		if err != nil {
			logger.Logger.Error("MODELS :: Error while scanning values", zap.String("requestId", uuidString), zap.Error(err))
			return postCommentData, count, err
		}

		postCommentData = append(postCommentData, singleData)
	}

	return postCommentData, count, err
}
