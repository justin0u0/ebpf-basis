#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <arpa/inet.h>
#include "amqp091.h"

char _license[] SEC("license") = "GPL";

#define MAX_ENTRIES 1024

struct {  
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, 1);
	__type(key, __u32);
	__type(value, __u32);
} packets_count_map SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, MAX_ENTRIES);
	__type(key, __u32);
	__type(value, __u32);
} amqp_frames_map SEC(".maps");

SEC("xdp/amqp_collect")
int xdp_amqp_collect(struct xdp_md *ctx) {
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

	// Check if the packet is a TCP data packet
	if (tcp->psh != 1) {
		return XDP_PASS;
	}

	// Check if the packet is an AMQP packet
	if (tcp->dest != htons(AMQP_PORT)) {
		return XDP_PASS;
	}

	struct amqp_frame_header *amqphdr = (void*)(ip + 1) + tcp->doff * 4;
	if ((void*)(amqphdr + 1) > data_end) {
		return XDP_PASS;
	}
	bpf_printk("AMQP frame header: type=%u, channel=%u, size=%u", amqphdr->type, htons(amqphdr->channel), htonl(amqphdr->size));

	struct amqp_method_frame *method = (void*)(amqphdr + 1);
	if ((void*)(method + 1) > data_end) {
		return XDP_PASS;
	}
	bpf_printk("AMQP method frame: class_id=%u, method_id=%u", htons(method->class_id), htons(method->method_id));

	__u32 packet_map_index = 0;
	__u32 *packet_count = bpf_map_lookup_elem(&packets_count_map, &packet_map_index);
	if (!packet_count) {
		__u32 zero = 0;
		bpf_map_update_elem(&packets_count_map, &packet_map_index, &zero, BPF_ANY);
		packet_count = bpf_map_lookup_elem(&packets_count_map, &packet_map_index);
	}

	if (packet_count) {
		if (*packet_count < MAX_ENTRIES) {
			bpf_map_update_elem(&amqp_frames_map, packet_count, method, BPF_ANY);
		}
		__sync_fetch_and_add(packet_count, 1);
	}

	return XDP_PASS;
}
