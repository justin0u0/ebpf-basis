#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <arpa/inet.h>

SEC("xdp/drop_amqp")
int xdp_drop_amqp(struct xdp_md *ctx) {
	void *data = (void *)(long)ctx->data;
	void *data_end = (void *)(long)ctx->data_end;

	struct ethhdr *eth = data;
	if (data + sizeof(*eth) > data_end) {
		return XDP_PASS;
	}

	// Check if the packet is an IP packet
	if (eth->h_proto != htons(ETH_P_IP)) {
		return XDP_PASS;
	}

	struct iphdr *ip = data + sizeof(*eth);
	if (data + sizeof(*eth) + sizeof(*ip) > data_end) {
		return XDP_PASS;
	}

	// Check if the packet is a TCP packet
	if (ip->protocol != IPPROTO_TCP) {
		return XDP_PASS;
	}

	struct tcphdr *tcp = data + sizeof(*eth) + sizeof(*ip);
	if (data + sizeof(*eth) + sizeof(*ip) + sizeof(*tcp) > data_end) {
		return XDP_PASS;
	}

	return XDP_PASS;
}
