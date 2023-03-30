package bpfgo

import (
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -Wall" skRedirect bpf/sk_redirect.c

func loadSkRedirectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sk-redirect [program]",
		Short: "Load the socket redirect program",
		Args:  cobra.ExactArgs(1),
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "sklookup",
			Short: "Load the socket lookup program",
			Run:   runLoadSkLookup,
		},
		&cobra.Command{
			Use:   "sockops",
			Short: "Load the socket operations program",
			Run:   runLoadSockOps,
		},
	)

	return cmd
}

func runLoadSkLookup(cmd *cobra.Command, args []string) {
	// Load pre-compiled programs into the kernel.
	objs := skRedirectObjects{}
	if err := loadSkRedirectObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %s", err)
	}
	defer objs.Close()

	netns, err := os.Open("/proc/self/ns/net")
	if err != nil {
		panic(err)
	}
	defer netns.Close()

	// Attach the programs.
	l, err := link.AttachNetNs(int(netns.Fd()), objs.SkLookup)
	if err != nil {
		log.Fatalf("could not attach the SK_LOOKUP program: %s", err)
	}
	defer l.Close()

	log.Printf("attached SK_LOOKUP program to netns %d\n", netns.Fd())

	ctx := cmd.Context()
	<-ctx.Done()

	log.Println("detaching SK_LOOKUP program")
}

func runLoadSockOps(cmd *cobra.Command, args []string) {
	cgroupPath, err := findCgroupPath()
	if err != nil {
		log.Fatalf("finding cgroup path: %s", err)
	}

	// Load pre-compiled programs into the kernel.
	objs := skRedirectObjects{}
	if err := loadSkRedirectObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %s", err)
	}
	defer objs.Close()

	l, err := link.AttachCgroup(link.CgroupOptions{
		Path:    cgroupPath,
		Program: objs.Sockops,
		Attach:  ebpf.AttachCGroupSockOps,
	})
	if err != nil {
		log.Fatalf("could not attach the SOCK_OPS program: %s", err)
	}
	defer l.Close()

	log.Printf("attached SOCK_OPS program to cgroup %s\n", cgroupPath)

	ctx := cmd.Context()
	<-ctx.Done()

	log.Println("detaching SOCK_OPS program")
}

func findCgroupPath() (string, error) {
	cgroupPath := "/sys/fs/cgroup"

	var st syscall.Statfs_t
	err := syscall.Statfs(cgroupPath, &st)
	if err != nil {
		return "", err
	}
	isCgroupV2Enabled := st.Type == unix.CGROUP2_SUPER_MAGIC
	if !isCgroupV2Enabled {
		cgroupPath = filepath.Join(cgroupPath, "unified")
	}
	return cgroupPath, nil
}
