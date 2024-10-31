package bootstrap

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (c *Container) Db() *pgx.Conn {
	if c.db == nil {
		ctx := context.Background()

		conn, err := pgx.Connect(ctx, "postgres://sso_user:sso123@192.168.1.252:5432/db_user?search_path=sso")
		if err != nil {
			fmt.Printf("failed to connect with error %v", err)
			return nil
		}

		c.db = conn
	}

	return c.db
}

func (c *Container) StopDb(ctx context.Context) {
	if c.db != nil {
		c.db.Close(ctx)
	}
}
