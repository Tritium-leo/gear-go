package model

type PostgresConfig struct {
	UserName string
	Password string
	Address  string
	DbName   string
	SSLMode  string
	LogMode  bool
}
