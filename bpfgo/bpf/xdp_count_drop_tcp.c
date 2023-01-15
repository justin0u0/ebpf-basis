#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <arpa/inet.h>

char _license[] SEC("license") = "GPL";

struct {  
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, 1);
	__type(key, uint32_t);
	__type(value, uint32_t);
} drop_packets_count_map SEC(".maps");

SEC("xdp/count_drop_tcp")
int xdp_count_drop_tcp(struct xdp_md *ctx) {
	void *data = (void *)(long)ctx->data;
	void *data_end = (void *)(long)ctx->data_end;

	struct ethhdr *eth = data;
	if ((void*)(eth + 1) > data_end) {
		return XDP_PASS;
	}

	// Check if the packet is an IP packet
	if (eth->h_proto != htons(ETH_P_IP)) {
		return XDP_PASS;
	}

	struct iphdr *ip = (void*)(eth + 1);
	if ((void*)(ip + 1) > data_end) {
		return XDP_PASS;
	}

	// Check if the packet is a TCP packet
	if (ip->protocol != IPPROTO_TCP) {
		return XDP_PASS;
	}

	struct tcphdr *tcp = (void*)(ip + 1);
	if ((void*)(tcp + 1) > data_end) {
		return XDP_PASS;
	}

	if (tcp->psh != 1) {
		return XDP_PASS;
	}

	// Drop data packets only
	uint32_t key = 0;
	uint32_t *value = bpf_map_lookup_elem(&drop_packets_count_map, &key);
	if (!value) {
		uint32_t zero = 0;
		bpf_map_update_elem(&drop_packets_count_map, &key, &zero, BPF_ANY);
		value = bpf_map_lookup_elem(&drop_packets_count_map, &key);
	}
	if (value) {
		__sync_fetch_and_add(value, 1);
	}

	return XDP_DROP;
}
