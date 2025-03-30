package user

import (
	"fmt"
	manager "github.com/gojue/ebpfmanager"
)
func OnEvent(cpu int, data []byte, perfmap *manager.PerfMap, manager *manager.Manager) {
    fmt.Printf("%s\n", data)
}