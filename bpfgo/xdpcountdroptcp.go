package bpfgo

import (
	"log"
	"net"

	"github.com/cilium/ebpf/link"
	"github.com/spf13/cobra"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -Wall" xdpCountDropTcp bpf/xdp_count_drop_tcp.c

func attachXdpCountDropTcp(cmd *cobra.Command, args []string) {
	ifaceName := args[0]

	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Fatalf("lookup network iface %q: %s", ifaceName, err)
	}

	// Load pre-compiled programs into the kernel.
	objs := xdpCountDropTcpObjects{}
	if err := loadXdpCountDropTcpObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %s", err)
	}
	defer objs.Close()

	// Attach the program.
	l, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.XdpCountDropTcp,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatalf("could not attach XDP program: %s", err)
	}
	defer l.Close()

	log.Printf("Attached XDP program to iface %q (index %d)", iface.Name, iface.Index)

	ctx := cmd.Context()
	<-ctx.Done()

	var (
		key   uint32 = 0
		count uint32
	)
	if err := objs.DropPacketsCountMap.Lookup(&key, &count); err != nil {
		log.Fatalf("could not read drop count: %s", err)
	}
	log.Printf("Dropped %d packets", count)

	log.Println("detaching XDP program")
}
