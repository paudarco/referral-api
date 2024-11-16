package storage

import (
	"time"
)

type Token interface {
	Set(key int, value interface{}, ttl time.Time)
	Get(key int) (interface{}, bool)
	StartCleanup()
	Cleanup()
	DeleteUserCode(id int)
}

type Storage struct {
	Token
}

func NewStorage() *Storage {
	return &Storage{
		Token: NewCache(),
	}
}
