#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

char _license[] SEC("license") = "GPL";

SEC("xdp/drop")
int xdp_drop(struct xdp_md *ctx) {
	return XDP_DROP;
}
