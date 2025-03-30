# eHook

> 本项目正在开发中，欢迎试用

## ✨ 介绍

- 快速、方便地构建你的 uprobe hook 模块。
- 提供一些方便的封装，锐意开发中。

## 🚀 运行环境

- 目前仅支持 ARM64 架构的 Android 系统，需要 ROOT 权限，推荐搭配 [KernelSU](https://github.com/tiann/KernelSU) 使用
- 系统内核版本5.10+ （可执行`uname -r`查看）

## 💕 使用

工作目录在 `/user` 目录下

- 在 `config.go` 里设置目标信息，如：

  ```go
  const PackageName = "com.android.myapplication"
  const LibraryName = "libc.so"
  // The offset to which the onEnter function will be attached. 0 if not used.
  const Enter_Offset = 0xAFF44
  // The offset to which the onLeave function will be attached. 0 if not used.
  const Leave_Offset = 0x9C158
  ```

- 在 `user.c` 中编写你的 eBPF 模块。可以使用任意的 eBPF API 如读写。

  ```c++
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
      LOG(b)
  }
  ```

  为了避免 eBPF 栈大小限制，尽量用 `VARIABLES_POOL` 定义变量池，用 GET 和 SET 操作变量。

  使用 LOG(char*) 打印信息，SUBMIT(buffer, len) 提交你的自定义数据到 OnEvent

- 如果需要，在 `listener.go` 中编写数据处理函数。

  ```
  func OnEvent(cpu int, data []byte, perfmap *manager.PerfMap, manager *manager.Manager) {
  	// Write your data handler here
      fmt.Printf("%s\n", data)
  }
  ```

## 🛫 构建

1. 环境准备

   本项目在 x86 Linux 下交叉编译

   ```
   sudo apt-get update
   sudo apt-get install golang==1.18
   sudo apt-get install clang==14
   export GOPROXY=https://goproxy.cn,direct
   export GO111MODULE=on
   
   git clone --recursive https://github.com/ShinoLeah/eDBG.git
   ./build_env.sh
   ```

2. 编译

   ```
   make
   ```

​	可以在 Makefile 中指定项目名称，产物在 `bin/` 目录下。