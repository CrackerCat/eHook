#include "vmlinux_510.h"
#include "bpf_helpers.h"
#include "bpf_tracing.h"

struct {                                                                                       
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);                                                                       
    __uint(max_entries, 1024);                                                         
    __type(key, int);                                                                    
    __type(value, __u32);                                                                
} events SEC(".maps");

struct {                                                                                       
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);                                                                       
    __uint(max_entries, 1024);                                                         
    __type(key, int);                                                                    
    __type(value, __u32);                                                                
} log_maps SEC(".maps");

#define VARIABLES_POOL(name) \
struct {                                                                                     \
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);                                                                     \
    __uint(max_entries, 1);                                                         \
    __type(key, u32);                                                                    \
    __type(value, struct name);                                                                \
} event_map SEC(".maps");


static __always_inline int __strlen(char *ptr) {
    int len = 0;
    for(char *i = ptr; *i != 0; i++) {
        len++;
    }
    return len;
}

#define LOG(ptr) log_internal(ptr, ctx);
static __always_inline void log_internal(char* ptr, struct pt_regs* ctx) {
    bpf_perf_event_output(ctx, &log_maps, BPF_F_CURRENT_CPU, ptr, __strlen(ptr));
}

#define SUBMIT(ptr,len) bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, ptr, len);

#define GET(name) \
({__u32 ___zero = 0;\
struct data_t *___data = bpf_map_lookup_elem(&event_map, &___zero);\
if (!___data) return;\
___data->name;})

#define SET(name, var) \
({__u32 ___zero = 0;\
struct data_t *___data = bpf_map_lookup_elem(&event_map, &___zero);\
if (!___data) return;\
___data->name = var;})