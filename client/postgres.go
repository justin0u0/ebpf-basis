package client

import (
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
)

var keepAlive bool

func runPostgresCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "postgres",
		Run: runPostgres,
	}

	cmd.Flags().BoolVarP(&keepAlive, "keep-alive", "k", keepAlive, "keep connection alive")

	return cmd
}

func runPostgres(cmd *cobra.Command, args []string) {
	// pgbouncer
	url := "postgres://postgres:postgres@localhost:6432/postgres?sslmode=disable"

	ctx := cmd.Context()

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, "SELECT 1")
	if err != nil {
		log.Fatalf("failed to execute a query: %v", err)
	}

	log.Println("executed a query")

	if keepAlive {
		<-ctx.Done()
	}
}
