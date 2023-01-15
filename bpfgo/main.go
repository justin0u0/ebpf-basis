package bpfgo

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cmd := &cobra.Command{Use: "bpfgo"}

	cmd.AddCommand(
		&cobra.Command{
			Use:  "xdp-amqp-collect [iface]",
			Run:  attachXdpAmqpCollect,
			Args: cobra.ExactArgs(1),
		},
		&cobra.Command{
			Use: "xdp-count-drop-tcp [iface]",
			Run: attachXdpCountDropTcp,
		},
		&cobra.Command{
			Use:  "xdp-drop-tcp [iface]",
			Run:  attachXdpDropTcp,
			Args: cobra.ExactArgs(1),
		},
	)

	return cmd
}
