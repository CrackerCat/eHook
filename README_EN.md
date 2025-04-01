# eHook

## ✨ Introduction
- Quickly and conveniently build your uprobe hook modules.
- Provides some convenient wrappers.

## 🚀 Requirements
- Currently only supports ARM64 architecture Android systems with ROOT permissions. Recommended to use with [KernelSU](https://github.com/tiann/KernelSU)[1](@ref)
- Kernel version 5.10+ (Check with `uname -r`)

## 💕 Usage
Work in the `/user` directory:

- Configure target information in `config.go`:
  ```go
  const PackageName = "com.android.myapplication"
  const LibraryName = "libc.so"
  // The offset to which the onEnter function will be attached. 0 if not used.
  const Enter_Offset = 0xAFF44
  // The offset to which the onLeave function will be attached. 0 if not used.
  const Leave_Offset = 0x9C158

- You can specify absolute library paths in `LibraryName`.

- Write your eBPF module in `user.c`. Use provided wrappers (see "API" section) or any eBPF APIs. Refer to eBPF Docs

  ```c
  struct data_t {
      int a;
      char b;
  };
  VARIABLES_POOL(data_t);
  
  static __always_inline void onEnter(struct pt_regs* ctx) {
      SET(a, 1);
      SET(b, 'c');
  }
  
  static __always_inline void onLeave(struct pt_regs* ctx) {
      char s = GET(b);
      LOG(&s, 1);
  }
  ```

- Implement data handlers in `listener.go` if needed:

  ```go
  func OnEvent(cpu int, data []byte, perfmap *manager.PerfMap, manager *manager.Manager) {
      // Write your data handler here
      fmt.Printf("%s\n", data)
  }
  ```

## 💭 API

- `VARIABLES_POOL`: Defines global variable pool for shared/large variables

  ```c
  struct data_t {
      int a;
      char b;
  };
  VARIABLES_POOL(data_t);
  ```

- `GET(name)`: Retrieves variable from pool (equivalent to `data->name`)

  ```c
  int a = GET(a);
  __builtin_memcpy(GET(b), "xxxxx", 5);
  ```

- `SET(name, var)`: Sets variable (equivalent to `data->name = var`). Use `GET` for string variables.

- `READ_KERN(x)`: Reads non-userspace variables (e.g., `READ_KERN(ctx->regs[0])`)

- `READ(ptr, len)`: Reads specified userspace address

- `WRITE(ptr, content)`: Writes to specified userspace address (must be writable)

- `LOG(char*, len)`: Outputs to console

- `SUBMIT(void*, len)`: Submits data to `OnEvent` in `listen.go`

## 🧑💻 Example

Bypass adb property check implementation:

```go
// config.go
const PackageName = "com.android.myapplication"
const LibraryName = "libc.so"
const Enter_Offset = 0xAFF44 //__system_property_get
const Leave_Offset = 0x9C158
```

```c
// user.c
#include "include/eHook.h"

struct data_t {
    __u64 X1;
};
VARIABLES_POOL(data_t);

static __always_inline void onEnter(struct pt_regs* ctx) {
    // Do not modify the name of 'onEnter' 'onLeave' or 'ctx'
    __u64 prop_name_addr = READ_KERN(ctx->regs[0]);
    char* s = READ(prop_name_addr, 14);
    if(!__builtin_memcmp(s, "init.svc.adbd", 14)) {
        SET(X1, READ_KERN(ctx->regs[1]));
    } else {
        SET(X1, 0);
    }
}

static __always_inline void onLeave(struct pt_regs* ctx) {
    if(GET(X1) != 0) {
        LOG("modified.\n", 10);
        WRITE(GET(X1), "stopped");
    }
}
```

## ⚠️ Notes

- eBPF functions have strict limitations (e.g., cannot use libc). Please research these constraints.
- `PackageName` parameter is for library targeting. Hooking system libraries like libc will affect all processes.

## 🛫 Building

1. **Environment Setup** (Cross-compile on x86 Linux):

   ```shell
   sudo apt-get update
   sudo apt-get install golang-1.18
   sudo apt-get install clang-15
   export GOPROXY=https://goproxy.cn,direct
   export GO111MODULE=on
   
   git clone --recursive https://github.com/ShinoLeah/eHook.git
   ./build_env.sh
   ```

2. **Compilation**:

   ```shell
   make
   ```

   - Modify `Makefile` for project naming. Outputs in `bin/`. 

3. **Run**:

   ```shell
   adb push bin/eHook_Untitled /data/local/tmp
   adb shell
   su
   cd /data/local/tmp
   chmod +x eHook_Untitled
   ./eHook_Untitled
   ```

## ❤️🩹 Others

- Star the repo if you find it useful 🌟
- Issues and PRs are welcome!