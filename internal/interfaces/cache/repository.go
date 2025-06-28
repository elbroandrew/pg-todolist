package cache

import (
	"time"
)


type Repository interface {
    RevokeToken(token string, ttl time.Duration) error
    IsTokenRevoked(token string) (bool, error)
    SetLastRefresh(userID uint, t time.Time) error
    GetLastRefresh(userID uint) (time.Time, error)
    PushDelete(data interface{}) (int64, error)
    PushUpdate(data interface{}) (int64, error) 
    Close() error
}

