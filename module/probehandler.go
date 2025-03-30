package module

import (
    "bytes"
	"fmt"
    "eHook/utils"
    "eHook/assets"
    "eHook/controller"
    "path/filepath"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/btf"
    "golang.org/x/sys/unix"
    "eHook/user"
    "math"
	manager "github.com/gojue/ebpfmanager"
)

type ProbeHandler struct {
	bpfManager        *manager.Manager
    bpfManagerOptions manager.Options
    BTF_File          string
}

func CreateProbeHandler(BTF_File string) *ProbeHandler {
    return &ProbeHandler{BTF_File: BTF_File}
}

func (this *ProbeHandler) OnEvent(cpu int, data []byte, perfmap *manager.PerfMap, manager *manager.Manager) {
    fmt.Printf("%s\n", data)
}

func (this *ProbeHandler) SetupManagerOptions() error {
    // 对于没有开启 CONFIG_DEBUG_INFO_BTF 的加载额外的 btf.Spec
    if this.BTF_File != "" {
        byteBuf, err := assets.Asset("assets/" + this.BTF_File)
        if err != nil {
            return fmt.Errorf("SetupManagerOptions failed, err:%v", err)
        }
        spec, err := btf.LoadSpecFromReader((bytes.NewReader(byteBuf)))
        if err != nil {
            return fmt.Errorf("SetupManagerOptions failed, err:%v", err)
        }
        this.bpfManagerOptions = manager.Options{
            DefaultKProbeMaxActive: 512,
            VerifierOptions: ebpf.CollectionOptions{
                Programs: ebpf.ProgramOptions{
                    LogSize:     2097152,
                    KernelTypes: spec,
                },
            },
            RLimit: &unix.Rlimit{
                Cur: math.MaxUint64,
                Max: math.MaxUint64,
            },
        }
    } else {
        this.bpfManagerOptions = manager.Options{
            DefaultKProbeMaxActive: 512,
            VerifierOptions: ebpf.CollectionOptions{
                Programs: ebpf.ProgramOptions{
                    LogSize:     2097152,
                },
            },
            RLimit: &unix.Rlimit{
                Cur: math.MaxUint64,
                Max: math.MaxUint64,
            },
        }
    }
    return nil
}

func (this *ProbeHandler) SetupManager(LibInfo *controller.LibraryInfo) error {
    probes := []*manager.Probe{}
    if user.Enter_Offset != 0 {
        sym := utils.RandStringBytes(8)
        probes = append(probes, &manager.Probe{
            Section:          "uprobe/probe_enter",
            EbpfFuncName:     "probe_enter",
            AttachToFuncName: sym,
            RealFilePath:     LibInfo.RealFilePath,
            BinaryPath:       LibInfo.LibPath,
            NonElfOffset:     LibInfo.NonElfOffset,
            UAddress:         user.Enter_Offset,
        })
    }

    if user.Leave_Offset != 0 {
        sym := utils.RandStringBytes(8)
        probes = append(probes, &manager.Probe{
            Section:          "uprobe/probe_leave",
            EbpfFuncName:     "probe_leave",
            AttachToFuncName: sym,
            RealFilePath:     LibInfo.RealFilePath,
            BinaryPath:       LibInfo.LibPath,
            NonElfOffset:     LibInfo.NonElfOffset,
            UAddress:         user.Leave_Offset,
        })
    }

    this.bpfManager = &manager.Manager{
        Probes: probes,
        PerfMaps: []*manager.PerfMap{
            &manager.PerfMap{
                Map: manager.Map{
                    Name: "log_maps",
                },
                PerfMapOptions: manager.PerfMapOptions{
                    DataHandler: this.OnEvent,
                },
            },
            &manager.PerfMap{
                Map: manager.Map{
                    Name: "events",
                },
                PerfMapOptions: manager.PerfMapOptions{
                    DataHandler: user.OnEvent,
                },
            },
        },
    }
    return nil
}

func (this *ProbeHandler) Run(LibInfo *controller.LibraryInfo) error {
    var bpfFileName = filepath.Join("assets", "ebpf_module.o")
    byteBuf, err := assets.Asset(bpfFileName)

    if err != nil {
        return fmt.Errorf("ProbeHandler.Run(): couldn't find asset %v .", err)
    }

    if err = this.SetupManager(LibInfo); err != nil {
        return fmt.Errorf("ProbeHandler.Run(): couldn't setup bootstrap manager %v .", err)
    }

    if err = this.bpfManager.InitWithOptions(bytes.NewReader(byteBuf), this.bpfManagerOptions); err != nil {
        return fmt.Errorf("ProbeHandler.Run(): couldn't init manager %v", err)
    }

    if err = this.bpfManager.Start(); err != nil {
        return fmt.Errorf("ProbeHandler.Run(): couldn't start bootstrap manager %v .", err)
    }
    return nil
}

func (this *ProbeHandler) Stop() error {
    return this.bpfManager.Stop(manager.CleanAll)
}

