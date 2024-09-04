package models

import (
	"fmt"
	"time"
)

type UserConnection struct {
	UserId    int64     `json:"user_id" binding:"required"`
	FriendId  int64     `json:"friend_id" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func (data UserConnection) Insert(uuidString string) (int64, error) {
	query := `INSERT INTO 
				user_connections
					(user_id, friend_at, created_at)
			VALUES
				($1, $2, $3)
			RETURNING id`
	queryToExecute := QueryStructToExecute{Query: query}
	id, err := queryToExecute.InsertOrUpdateOperations(uuidString, data.UserId, data.FriendId, time.Now().UTC())
	return id, err
}

func (data UserConnection) Remove(uuidString string) (bool, error) {
	query := fmt.Sprintf(`DELETE FROM 
				user_connections
			  WHERE user_id=%d AND friend_id=%d`, data.UserId, data.FriendId)
	queryToExecute := QueryStructToExecute{Query: query}
	status, err := queryToExecute.DeleteOperation(uuidString)
	return status, err
}
