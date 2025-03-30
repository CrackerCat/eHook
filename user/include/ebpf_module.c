#include "../user.c"

char __license[] SEC("license") = "GPL";
__u32 _version SEC("version") = 0xFFFFFFFE;

SEC("uprobe/probe_enter") 
int probe_enter(struct pt_regs* ctx) {
    onEnter(ctx);
    return 0;
}

SEC("uprobe/probe_leave") 
int probe_leave(struct pt_regs* ctx) {
    onLeave(ctx);
    return 0;
}