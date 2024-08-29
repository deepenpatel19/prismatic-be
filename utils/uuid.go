package utils

import (
	"github.com/google/uuid"
)

func GetUUID() string {
	uuidId := uuid.New()
	uuidString := uuidId.String()
	return uuidString
}
