#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>
#include <linux/if_ether.h>
#include <linux/if_packet.h>
#include <linux/ip.h>

char _license[] SEC("license") = "GPL";

struct {
	__uint(type, BPF_MAP_TYPE_SOCKMAP);
	__uint(max_entries, 20);
	__type(key, __u32);   // array index
	__type(value, __u32); // socket FD
} sockmap SEC(".maps");

SEC("sk_skb/parser/redirect")
int sk_skb(struct __sk_buff *skb) {
	bpf_printk("[sk_skb/parser] skb->local_port [%d], skb->remote_port [%d]",
		skb->local_port, bpf_ntohs(skb->remote_port));
	return skb->len;
}

SEC("sk_skb/verdict/redirect")
int sk_skb_verdict(struct __sk_buff *skb) {
	bpf_printk("[sk_skb/verdict] skb->local_port [%d], skb->remote_port [%d]",
		skb->local_port, bpf_ntohs(skb->remote_port));
	
	return bpf_sk_redirect_map(skb, &sockmap, 0, 0);
}

SEC("sk_lookup/redirect")
int sk_lookup(struct bpf_sk_lookup *ctx) {
	bpf_printk("[sk_lookup] ctx->local_port [%d], ctx->remote_port [%d]",
		ctx->local_port, bpf_ntohs(ctx->remote_port));

	return SK_PASS;
}

SEC("sockops/redirect")
int sockops(struct bpf_sock_ops *skops) {
	switch (skops->op) {
	case BPF_SOCK_OPS_PASSIVE_ESTABLISHED_CB: // SYN
		bpf_printk("[BPF_SOCK_OPS_PASSIVE_ESTABLISHED_CB] local port [%d], remote port [%d]",
			skops->local_port, bpf_ntohs(skops->remote_port));
		break;
	case BPF_SOCK_OPS_ACTIVE_ESTABLISHED_CB: // SYN-ACK
		bpf_printk("[BPF_SOCK_OPS_ACTIVE_ESTABLISHED_CB] local port [%d], remote port [%d]",
			skops->local_port, bpf_ntohs(skops->remote_port));
		break;
	default:
		break;
	}

	return 0;
}
