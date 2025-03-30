module eHook

go 1.18

require (
	github.com/cilium/ebpf v0.17.3
	github.com/gojue/ebpfmanager v0.4.6
	github.com/shuLhan/go-bindata v4.0.0+incompatible
	golang.org/x/exp v0.0.0-20230224173230-c95f2b4c22f2
	golang.org/x/sys v0.30.0
)

replace github.com/gojue/ebpfmanager => ./ebpfmanager

replace github.com/cilium/ebpf => ./ebpf

require (
	github.com/avast/retry-go v3.0.0+incompatible // indirect
	github.com/florianl/go-tc v0.4.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/josharian/native v1.1.0 // indirect
	github.com/mdlayher/netlink v1.7.2 // indirect
	github.com/mdlayher/socket v0.4.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/vishvananda/netlink v1.1.0 // indirect
	github.com/vishvananda/netns v0.0.0-20191106174202-0a2b9b5464df // indirect
	golang.org/x/net v0.36.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
)
