#include "include/eHook.h"

struct data_t {
	__u64 regs[31];
	__u64 sp;
	__u64 pc;
    __u64 pstate;
};

static __always_inline data_t* getData() {
    __u32 zero = 0;
    struct data_t *data = bpf_map_lookup_elem(&event_map, &zero);
    if (!data) return;
    return data;
}

static __always_inline void onEnter(struct pt_regs* ctx) {
    // Do not modify this name
    LOG("OnEnter Triggered.");
    SUBMIT("xxxx", 4);
    struct data_t *data = getData();
    for(int i = 0; i < 31; ++i) {
        bpf_probe_read_kernel(&data->regs[i], sizeof(data->regs[i]), &ctx->regs[i]);
    }
    bpf_probe_read_kernel(&data->sp, sizeof(data->sp), &ctx->sp);
    bpf_probe_read_kernel(&data->pc, sizeof(data->pc), &ctx->pc);
    bpf_probe_read_kernel(&data->pstate, sizeof(data->pstate), &ctx->pstate);
    char buffer[200];
    bpf_probe_read_user_str(buffer, 200, (char*) ctx->regs[0]);
    LOG(buffer);
}

static __always_inline void onLeave(struct pt_regs* ctx) {
    LOG("OnLeave Triggered.");
    SUBMIT("xxxx", 4);
    struct data_t *data = getData();
    char new_value[] = "stopped";
    int ret = bpf_probe_write_user((char*)data->regs[1], new_value, sizeof(new_value));
    if(ret < 0) {
        LOG("Write Memory Failed.");
    } else {
        LOG("Write Memory Success.");
    }
}