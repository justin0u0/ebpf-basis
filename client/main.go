package client

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: "client",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use: "amqp-publish",
			Run: runAmqpPublish,
		},
		&cobra.Command{
			Use: "amqp-consume",
			Run: runAmqpConsume,
		},
		runPostgresCommand(),
	)

	return cmd
}
