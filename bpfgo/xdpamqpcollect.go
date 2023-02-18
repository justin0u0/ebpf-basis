package bpfgo

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/spf13/cobra"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -Wall" xdpAmqpCollect bpf/xdp_amqp_collect.c

type amqpFrameHeader struct {
	Type    uint8
	Channel uint16
	Size    uint32
}

func attachXdpAmqpCollect(cmd *cobra.Command, args []string) {
	ifaceName := args[0]

	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Fatalf("lookup network iface %q: %s", ifaceName, err)
	}

	// Load pre-compiled programs into the kernel.
	objs := xdpAmqpCollectObjects{}
	if err := loadXdpAmqpCollectObjects(&objs, nil); err != nil {
		var ve *ebpf.VerifierError
		if errors.As(err, &ve) {
			for _, l := range ve.Log {
				fmt.Println(l)
			}
		}
		log.Fatalf("could not load XDP program: %s", err)
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
			key    uint32
			val    []byte
			header amqpFrameHeader
		)

		if iter.Next(&key, &val) {
			header.Type = uint8(val[0])
			header.Channel = binary.BigEndian.Uint16(val[1:3])
			header.Size = binary.BigEndian.Uint32(val[3:7])

			log.Printf("packet %d: type=%d, channel=%d, size=%d\n", i, header.Type, header.Channel, header.Size)
		}
	}

	if err := iter.Err(); err != nil {
		log.Println("error iterating map:", err)
	}

	log.Println("detaching XDP program")
}
