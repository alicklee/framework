package utils

import (
	uuid "github.com/satori/go.uuid"
)

func NewUuid() string {
	uuid := uuid.Must(uuid.NewV4(), nil)
	return uuid.String()
}
