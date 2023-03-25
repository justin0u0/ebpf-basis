package bpfgo

import (
	"log"
	"net"

	"github.com/cilium/ebpf/link"
	"github.com/spf13/cobra"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -Wall" xdpDropTcp bpf/xdp_drop_tcp.c

// attachXdpDropTcp attaches the XDP program to the given interface.
// The function blocks until the program is interrupted.
func attachXdpDropTcp(cmd *cobra.Command, args []string) {
	ifaceName := args[0]

	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Fatalf("lookup network iface %q: %s", ifaceName, err)
	}

	// Load pre-compiled programs into the kernel.
	objs := xdpDropTcpObjects{}
	if err := loadXdpDropTcpObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %s", err)
	}
	defer objs.Close()

	// Attach the program.
	l, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.XdpDropTcp,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatalf("could not attach XDP program: %s", err)
	}
	defer l.Close()

	log.Printf("attached XDP program to iface %q (index %d)", iface.Name, iface.Index)

	ctx := cmd.Context()
	<-ctx.Done()

	log.Println("detaching XDP program")
}
