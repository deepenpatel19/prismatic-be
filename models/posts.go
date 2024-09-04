package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/deepenpatel19/prismatic-be/logger"
	"go.uber.org/zap"
)

type Post struct {
	UserId      int64  `json:"user_id" form:"user_id" binding:"required"`
	Title       string `json:"title" form:"title" binding:"title"`
	Description string `json:"description" form:"description"`
	File        string `json:"file" form:"file"`
}

func (data Post) Insert(uuidString string) (int64, error) {
	query := fmt.Sprintf(`INSERT INTO
				posts
					(user_id, title, description, file, created_at)
				VALUES
					(%d, '%s', '%s', '%s', '%s')
				RETURNING id`,
		data.UserId,
		data.Title,
		data.Description,
		data.File,
		time.Now().UTC().Format(time.RFC3339),
	)

	queryToExecute := QueryStructToExecute{Query: query}
	id, err := queryToExecute.InsertOrUpdateOperations(uuidString)
	return id, err

}

func (data Post) Update(uuidString string, Id int64) (int64, error) {
	query := fmt.Sprintf(`UPDATE 
								posts
							SET title='%s', description='%s', file='%s'
							WHERE id=%d
							RETURNING id`,
		data.Title,
		data.Description,
		data.File,
		Id,
	)
	queryToExecute := QueryStructToExecute{Query: query}
	id, err := queryToExecute.InsertOrUpdateOperations(uuidString)
	return id, err
}

func DeletePost(uuidString string, Id int64) (bool, error) {
	query := fmt.Sprintf(`DELETE 
							FROM posts
							WHERE id=%d`, Id)
	queryToExecute := QueryStructToExecute{Query: query}
	status, err := queryToExecute.DeleteOperation(uuidString)
	return status, err
}

func FetchPosts(uuidString string, limit int, offset int) ([]Post, int64, error) {
	var rows []map[string]any
	var count int64
	var err error
	var postData []Post

	query := fmt.Sprintf(`SELECT 
							id,
							title,
							description,
							file,
							created_at
						  FROM
						  	posts
						  ORDER BY id DESC
						  LIMIT %d OFFSET %d`, limit, offset)
	queryToExecute := QueryStructToExecute{Query: query}
	rows, count, err = queryToExecute.FetchRows(uuidString)
	if err != nil {
		return postData, count, err
	}

	for _, data := range rows {
		var singleData Post
		jsonbody, err := json.Marshal(data)
		if err != nil {
			// do error check
			// fmt.Println(err)
			logger.Logger.Error("MODELS :: Post => Error while json marshalling", zap.Error(err), zap.String("requestId", uuidString))
			return postData, count, err
		}

		if err := json.Unmarshal(jsonbody, &singleData); err != nil {
			// do error check
			// fmt.Println(err)
			logger.Logger.Error("MODELS :: Post => Error while json unmarshalling", zap.Error(err), zap.String("requestId", uuidString))
			return postData, count, err
		}

		postData = append(postData, singleData)
	}

	return postData, count, err
}
