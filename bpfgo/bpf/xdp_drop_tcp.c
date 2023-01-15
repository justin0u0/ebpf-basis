#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <arpa/inet.h>

char _license[] SEC("license") = "GPL";

SEC("xdp/drop_tcp")
int xdp_drop_tcp(struct xdp_md *ctx) {
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

	// Drop data packets only
	if (tcp->psh == 1) {
		return XDP_DROP;
	}

	return XDP_PASS;
}
