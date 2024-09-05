package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/deepenpatel19/prismatic-be/logger"
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
	var rows []map[string]any
	var count int64
	var err error
	var postCommentData []PostCommentResponse

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
	queryToExecute := QueryStructToExecute{Query: query}
	rows, count, err = queryToExecute.FetchRows(uuidString)
	if err != nil {
		return postCommentData, count, err
	}

	for _, data := range rows {
		var singleData PostCommentResponse
		jsonbody, err := json.Marshal(data)
		if err != nil {
			logger.Logger.Error("MODELS :: Post comment => Error while json marshalling", zap.Error(err), zap.String("requestId", uuidString))
			return postCommentData, count, err
		}

		if err := json.Unmarshal(jsonbody, &singleData); err != nil {
			logger.Logger.Error("MODELS :: Post comment => Error while json unmarshalling", zap.Error(err), zap.String("requestId", uuidString))
			return postCommentData, count, err
		}

		postCommentData = append(postCommentData, singleData)
	}

	return postCommentData, count, err
}
