package config

const (
	HostAddress = "0.0.0.0:5000"
	User        = "forum"
	Password    = "forum"
	DBName      = "forum"
	SSLMode     = "disable"
	// MaxConn     = 1000
	MaxConn  = 100000
	DBSchema = "db_create.sql"
	DBPath   = "postgresql://forum:forum@localhost:5432/forum"
)
