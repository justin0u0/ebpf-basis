all: build

BPF_TARGETS := xdp_amqp_collect xdp_count_drop_tcp xdp_drop xdp_pass xdp_drop_tcp xdp_echo
BPF_OBJECTS := $(addprefix bpfgo/bpf/,$(addsuffix .o,$(BPF_TARGETS)))

build: generate bin/cmd

generate: $(BPF_OBJECTS)
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
