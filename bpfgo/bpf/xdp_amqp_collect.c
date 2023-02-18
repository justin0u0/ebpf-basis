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
	__type(value, struct amqp_frame_header);
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
	if (tcp->dest != htons(AMQP_PORT) && tcp->source != htons(AMQP_PORT)) {
		return XDP_PASS;
	}

	// Parse all AMQP frames (may be multiple frames in a single TCP packet)
	void* amqp_frame_start = (void*)(ip + 1) + tcp->doff * 4;

	for (int frames = 0; frames < AMQP_MAX_NUM_FRAMES; ++frames) {
		struct amqp_frame_header *amqphdr = amqp_frame_start;
		if ((void*)(amqphdr + 1) > data_end) {
			return XDP_PASS;
		}

		uint32_t size = htonl(amqphdr->size);
		bpf_printk("AMQP frame header: type=%u, channel=%u, size=%u", amqphdr->type, htons(amqphdr->channel), size);

		if (size > AMQP_MAX_FRAME_SIZE) {
			// We ignore frames that are too large
			return XDP_PASS;
		}

		switch (amqphdr->type) {
			case AMQP_FRAME_HEADER_TYPE_METHOD: {
				struct amqp_method_frame *method = (void*)(amqphdr + 1);
				if ((void*)(method + 1) > data_end) {
					return XDP_PASS;
				}

				bpf_printk("AMQP method frame: class_id=%u, method_id=%u", htons(method->class_id), htons(method->method_id));
				break;
			}
			case AMQP_FRAME_HEADER_TYPE_HEADER: {
				struct amqp_header_frame *header = (void*)(amqphdr + 1);
				if ((void*)(header + 1) > data_end) {
					return XDP_PASS;
				}

				bpf_printk("AMQP header frame: class_id=%u, weight=%u, body_size=%llu", htons(header->class_id), htons(header->weight), htobe64(header->body_size));
				break;
			}
			case AMQP_FRAME_HEADER_TYPE_BODY: {
				bpf_printk("AMQP body frame");
				break;
			}
			default:
				break;
		}

		// Collect the AMQP frame header
		__u32 packet_map_index = 0;
		__u32 *packet_count = bpf_map_lookup_elem(&packets_count_map, &packet_map_index);
		if (!packet_count) {
			__u32 zero = 0;
			bpf_map_update_elem(&packets_count_map, &packet_map_index, &zero, BPF_ANY);
			packet_count = bpf_map_lookup_elem(&packets_count_map, &packet_map_index);
		}

		if (packet_count) {
			if (*packet_count < MAX_ENTRIES) {
				bpf_map_update_elem(&amqp_frames_map, packet_count, amqphdr, BPF_ANY);
			}
			__sync_fetch_and_add(packet_count, 1);
		}

		// Check AMQP frame end byte
		uint8_t *amqp_frame_end = (void*)(amqphdr + 1) + size;
		if ((void*)(amqp_frame_end + 1) > data_end) {
			return XDP_PASS;
		}
		if (*amqp_frame_end != AMQP_FRAME_END) {
			bpf_printk("AMQP frame end mismatch");
			return XDP_PASS;
		}

		// Move to the next AMQP frame
		amqp_frame_start = (void*)(amqphdr + 1) + size + 1;
	}

	return XDP_PASS;
}
