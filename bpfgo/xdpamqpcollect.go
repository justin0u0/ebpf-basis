package bpfgo

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/cilium/ebpf/link"
	"github.com/spf13/cobra"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -Wall" xdpAmqpCollect bpf/xdp_amqp_collect.c

func attachXdpAmqpCollect(cmd *cobra.Command, args []string) {
	ifaceName := args[0]

	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Fatalf("lookup network iface %q: %s", ifaceName, err)
	}

	// Load pre-compiled programs into the kernel.
	objs := xdpAmqpCollectObjects{}
	if err := loadXdpAmqpCollectObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %s", err)
	}
	defer objs.Close()

	// Attach the program.
	l, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.XdpAmqpCollect,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatalf("could not attach XDP program: %s", err)
	}
	defer l.Close()

	log.Printf("Attached XDP program to iface %q (index %d)", iface.Name, iface.Index)

	ctx := cmd.Context()
	<-ctx.Done()

	var packetCount uint32
	if err := objs.xdpAmqpCollectMaps.PacketsCountMap.Lookup(uint32(0), &packetCount); err != nil {
		log.Fatalf("could not get packet count: %s", err)
	}
	log.Printf("packet count: %d", packetCount)

	iter := objs.xdpAmqpCollectMaps.AmqpFramesMap.Iterate()
	for i := uint32(0); i < packetCount; i++ {
		var (
			key      uint32
			val      []byte
			classID  uint16
			methodID uint16
		)

		if iter.Next(&key, &val) {
			classID = binary.BigEndian.Uint16(val[0:2])
			methodID = binary.BigEndian.Uint16(val[2:4])

			log.Printf("classID: %d, methodID: %d", classID, methodID)
		}
	}

	if err := iter.Err(); err != nil {
		log.Println("error iterating map:", err)
	}

	log.Println("detaching XDP program")
}
