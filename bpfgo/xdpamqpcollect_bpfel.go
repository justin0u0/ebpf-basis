// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64 || amd64p32 || arm || arm64 || mips64le || mips64p32le || mipsle || ppc64le || riscv64
// +build 386 amd64 amd64p32 arm arm64 mips64le mips64p32le mipsle ppc64le riscv64

package bpfgo

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

// loadXdpAmqpCollect returns the embedded CollectionSpec for xdpAmqpCollect.
func loadXdpAmqpCollect() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_XdpAmqpCollectBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load xdpAmqpCollect: %w", err)
	}

	return spec, err
}

// loadXdpAmqpCollectObjects loads xdpAmqpCollect and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*xdpAmqpCollectObjects
//	*xdpAmqpCollectPrograms
//	*xdpAmqpCollectMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadXdpAmqpCollectObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadXdpAmqpCollect()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// xdpAmqpCollectSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type xdpAmqpCollectSpecs struct {
	xdpAmqpCollectProgramSpecs
	xdpAmqpCollectMapSpecs
}

// xdpAmqpCollectSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type xdpAmqpCollectProgramSpecs struct {
	XdpAmqpCollect *ebpf.ProgramSpec `ebpf:"xdp_amqp_collect"`
}

// xdpAmqpCollectMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type xdpAmqpCollectMapSpecs struct {
	AmqpFramesMap   *ebpf.MapSpec `ebpf:"amqp_frames_map"`
	PacketsCountMap *ebpf.MapSpec `ebpf:"packets_count_map"`
}

// xdpAmqpCollectObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadXdpAmqpCollectObjects or ebpf.CollectionSpec.LoadAndAssign.
type xdpAmqpCollectObjects struct {
	xdpAmqpCollectPrograms
	xdpAmqpCollectMaps
}

func (o *xdpAmqpCollectObjects) Close() error {
	return _XdpAmqpCollectClose(
		&o.xdpAmqpCollectPrograms,
		&o.xdpAmqpCollectMaps,
	)
}

// xdpAmqpCollectMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadXdpAmqpCollectObjects or ebpf.CollectionSpec.LoadAndAssign.
type xdpAmqpCollectMaps struct {
	AmqpFramesMap   *ebpf.Map `ebpf:"amqp_frames_map"`
	PacketsCountMap *ebpf.Map `ebpf:"packets_count_map"`
}

func (m *xdpAmqpCollectMaps) Close() error {
	return _XdpAmqpCollectClose(
		m.AmqpFramesMap,
		m.PacketsCountMap,
	)
}

// xdpAmqpCollectPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadXdpAmqpCollectObjects or ebpf.CollectionSpec.LoadAndAssign.
type xdpAmqpCollectPrograms struct {
	XdpAmqpCollect *ebpf.Program `ebpf:"xdp_amqp_collect"`
}

func (p *xdpAmqpCollectPrograms) Close() error {
	return _XdpAmqpCollectClose(
		p.XdpAmqpCollect,
	)
}

func _XdpAmqpCollectClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed xdpamqpcollect_bpfel.o
var _XdpAmqpCollectBytes []byte
