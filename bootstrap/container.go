package bootstrap

import "github.com/jackc/pgx/v5"

type Container struct {
	db *pgx.Conn
}

func InitContainer() *Container {
	return &Container{}
}
