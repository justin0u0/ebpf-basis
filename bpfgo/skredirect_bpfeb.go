// Code generated by bpf2go; DO NOT EDIT.
//go:build arm64be || armbe || mips || mips64 || mips64p32 || ppc64 || s390 || s390x || sparc || sparc64
// +build arm64be armbe mips mips64 mips64p32 ppc64 s390 s390x sparc sparc64

package bpfgo

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

// loadSkRedirect returns the embedded CollectionSpec for skRedirect.
func loadSkRedirect() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_SkRedirectBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load skRedirect: %w", err)
	}

	return spec, err
}

// loadSkRedirectObjects loads skRedirect and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*skRedirectObjects
//	*skRedirectPrograms
//	*skRedirectMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadSkRedirectObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadSkRedirect()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// skRedirectSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type skRedirectSpecs struct {
	skRedirectProgramSpecs
	skRedirectMapSpecs
}

// skRedirectSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type skRedirectProgramSpecs struct {
	SkLookup     *ebpf.ProgramSpec `ebpf:"sk_lookup"`
	SkSkb        *ebpf.ProgramSpec `ebpf:"sk_skb"`
	SkSkbVerdict *ebpf.ProgramSpec `ebpf:"sk_skb_verdict"`
	Sockops      *ebpf.ProgramSpec `ebpf:"sockops"`
}

// skRedirectMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type skRedirectMapSpecs struct {
	Sockmap *ebpf.MapSpec `ebpf:"sockmap"`
}

// skRedirectObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadSkRedirectObjects or ebpf.CollectionSpec.LoadAndAssign.
type skRedirectObjects struct {
	skRedirectPrograms
	skRedirectMaps
}

func (o *skRedirectObjects) Close() error {
	return _SkRedirectClose(
		&o.skRedirectPrograms,
		&o.skRedirectMaps,
	)
}

// skRedirectMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadSkRedirectObjects or ebpf.CollectionSpec.LoadAndAssign.
type skRedirectMaps struct {
	Sockmap *ebpf.Map `ebpf:"sockmap"`
}

func (m *skRedirectMaps) Close() error {
	return _SkRedirectClose(
		m.Sockmap,
	)
}

// skRedirectPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadSkRedirectObjects or ebpf.CollectionSpec.LoadAndAssign.
type skRedirectPrograms struct {
	SkLookup     *ebpf.Program `ebpf:"sk_lookup"`
	SkSkb        *ebpf.Program `ebpf:"sk_skb"`
	SkSkbVerdict *ebpf.Program `ebpf:"sk_skb_verdict"`
	Sockops      *ebpf.Program `ebpf:"sockops"`
}

func (p *skRedirectPrograms) Close() error {
	return _SkRedirectClose(
		p.SkLookup,
		p.SkSkb,
		p.SkSkbVerdict,
		p.Sockops,
	)
}

func _SkRedirectClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed skredirect_bpfeb.o
var _SkRedirectBytes []byte
