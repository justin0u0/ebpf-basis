#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

char _license[] SEC("license") = "GPL";

SEC("xdp/pass")
int xdp_pass(struct xdp_md *ctx) {
	return XDP_PASS;
}
