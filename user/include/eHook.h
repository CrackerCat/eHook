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


#define LOG(ptr, len) bpf_perf_event_output(ctx, &log_maps, BPF_F_CURRENT_CPU, ptr, len);

#define SUBMIT(ptr, len) bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, ptr, len);

#define GET(name) \
({\
    __u32 ___zero = 0;\
    struct data_t *___data = bpf_map_lookup_elem(&event_map, &___zero);\
    if (!___data) return;\
    ___data->name;\
})

#define SET(name, var) \
({\
    __u32 ___zero = 0;\
    struct data_t *___data = bpf_map_lookup_elem(&event_map, &___zero);\
    if (!___data) return;\
    ___data->name = var;\
})

#define READ_KERN(ptr)                                                                         \
({                                                                                         \
    typeof(ptr) _val;                                                                      \
    __builtin_memset((void *) &_val, 0, sizeof(_val));                                     \
    bpf_probe_read((void *) &_val, sizeof(_val), &ptr);                                    \
    _val;                                                                                  \
})

#define WRITE(addr, content)                                                                         \
({                                                                                         \
    bpf_probe_write_user((void*) addr, content, sizeof(content));                                                                                 \
})

#define READ(ptr, len)                                                                         \
({                                                                                         \
    char _val[len+1];                                                                      \
    __builtin_memset((void *) &_val, 0, sizeof(_val));                                     \
    bpf_probe_read_user((void *) _val, len, (void*) ptr);                               \
    _val;                                                                                  \
})
