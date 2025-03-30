TARGET_NAME ?= Untitled


CMD_CLANG ?= clang
CMD_GO ?= go
CMD_RM ?= rm
CMD_BPFTOOL ?= bpftool
BUILD_PATH ?= ./build
LINUX_ARCH = arm64
BUILD_TAGS ?=
TARGET_ARCH = $(LINUX_ARCH)
ifeq ($(BUILD_TAGS),forarm)
BUILD_TAGS := -tags forarm
TARGET_ARCH = arm
endif

.PHONY: all
all: ebpf_module genbtf assets build 
	@echo $(shell date)

.PHONY: clean
clean:
	$(CMD_RM) -f assets/*.d
	$(CMD_RM) -f assets/*.o
	$(CMD_RM) -f assets/ebpf_probe.go
	$(CMD_RM) -f bin/eHook_$(TARGET_NAME)

.PHONY: ebpf_module
ebpf_module:
	clang \
	-D__TARGET_ARCH_$(TARGET_ARCH) \
	--target=bpf \
	-c \
	-nostdlibinc \
	-no-canonical-prefixes \
	-O2 \
	-I       libbpf/src \
	-I       ebpf_module \
	-g \
	-o assets/ebpf_module.o \
	user/include/ebpf_module.c

.PHONY: assets
assets:
	$(CMD_GO) run github.com/shuLhan/go-bindata/cmd/go-bindata -pkg assets -o "assets/ebpf_probe.go" $(wildcard ./config/config_syscall_*.json ./assets/*.o ./assets/*_min.btf)

.PHONY: genbtf
genbtf:
	cd assets && ./$(CMD_BPFTOOL) gen min_core_btf rock5b-5.10-f9d1b1529-arm64.btf rock5b-5.10-arm64_min.btf ebpf_module.o
	cd assets && ./$(CMD_BPFTOOL) gen min_core_btf a12-5.10-arm64.btf a12-5.10-arm64_min.btf ebpf_module.o

.PHONY: build
build:
	GOARCH=arm64 GOOS=android $(CMD_GO) build $(BUILD_TAGS) -ldflags "-w -s -extldflags '-Wl,--hash-style=sysv'" -o bin/eHook_$(TARGET_NAME) .
