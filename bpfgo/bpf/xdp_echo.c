#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <arpa/inet.h>

char _license[] SEC("license") = "GPL";

SEC("xdp/echo")
int xdp_echo(struct xdp_md *ctx) {
	void *data = (void *)(long)ctx->data;
	void *data_end = (void *)(long)ctx->data_end;

	struct ethhdr *eth = data;
	if ((void*)(eth + 1) > data_end) {
		return XDP_PASS;
	}

	bpf_printk("Packet received from %02x:%02x:%02x:%02x:%02x:%02x to %02x:%02x:%02x:%02x:%02x:%02x", eth->h_source[0], eth->h_source[1], eth->h_source[2], eth->h_source[3], eth->h_source[4], eth->h_source[5], eth->h_dest[0], eth->h_dest[1], eth->h_dest[2], eth->h_dest[3], eth->h_dest[4], eth->h_dest[5]);

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

	void* payload = (void*)(ip + 1) + tcp->doff * 4;
	if (payload > data_end) {
		return XDP_PASS;
	}
	uint32_t payload_size = data_end - payload;

	unsigned char* saddr = (unsigned char*)&ip->saddr;
	unsigned char* daddr = (unsigned char*)&ip->daddr;
	bpf_printk("TCP Packet received size %d from %d.%d.%d.%d:%d to %d.%d.%d.%d:%d", payload_size, saddr[0], saddr[1], saddr[2], saddr[3], ntohs(tcp->source), daddr[0], daddr[1], daddr[2], daddr[3], ntohs(tcp->dest));

	// Swap the source and destination MAC addresses
	__u8 tmp[ETH_ALEN];
	__builtin_memcpy(tmp, eth->h_dest, ETH_ALEN);
	__builtin_memcpy(eth->h_dest, eth->h_source, ETH_ALEN);
	__builtin_memcpy(eth->h_source, tmp, ETH_ALEN);

	// Swap the source and destination IP addresses
	__u32 tmp_ip = ip->daddr;
	ip->daddr = ip->saddr;
	ip->saddr = tmp_ip;

	// Swap the source and destination TCP ports
	__u16 tmp_port = tcp->dest;
	tcp->dest = tcp->source;
	tcp->source = tmp_port;

	// Set TCP ACK
	bpf_printk("Received TCP ACK %u, SEQ %u", ntohl(tcp->ack_seq), ntohl(tcp->seq));
	__u32 tmp_ack = tcp->ack_seq;
	tcp->ack_seq = htonl(ntohl(tcp->seq) + payload_size);
	tcp->seq = tmp_ack;
	bpf_printk("Sending TCP ACK %u, SEQ %u", ntohl(tcp->ack_seq), ntohl(tcp->seq));
	tcp->psh = 0;
	tcp->ack = 1;

	// Remove the payload
	ip->tot_len = htons(ntohs(ip->tot_len) - payload_size);
	if (bpf_xdp_adjust_tail(ctx, -payload_size)) {
		bpf_printk("Failed to adjust tail");
		return XDP_PASS;
	}

	return XDP_TX;
}
