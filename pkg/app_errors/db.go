package app_errors

var (
    ErrDBError          = New(500, "database error")
    ErrDBConnection     = New(500, "database connection error")
    ErrRedisUnavailable = New(500, "redis unavailable")
    ErrDuplicateKey     = New(409, "duplicate key")
)