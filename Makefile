all: build

BPF_TARGETS := xdp_amqp_collect xdp_drop xdp_pass xdp_drop_tcp
BPF_OBJECTS := $(addprefix bpfgo/bpf/,$(addsuffix .o,$(BPF_TARGETS)))
BPF_SOURCES := $(addprefix bpfgo/bpf/,$(addsuffix .c,$(BPF_TARGETS)))

build: generate bin/cmd $(BPF_OBJECTS) dc-build

generate: $(BPF_SOURCES)
	go generate ./bpfgo/...

bpfgo/bpf/%.o: bpfgo/bpf/%.c
	clang -O2 -Wall -target bpf -c $< -o $@

bin/cmd: $(shell find . -name '*.go')
	@mkdir -p bin/
	go build -o $@ ./cmd/...

dc-build:
	docker-compose build

clean:
	rm -f bpfgo/bpf/*.o
	rm -f bin/cmd
