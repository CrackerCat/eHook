#include "include/eHook.h"

struct data_t {
    
};
VARIABLES_POOL(data_t);

static __always_inline void onEnter(struct pt_regs* ctx) {
    // Do not modify the name of 'onEnter' 'onLeave' or 'ctx'
}

static __always_inline void onLeave(struct pt_regs* ctx) {
}