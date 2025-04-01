package user

import (
	manager "github.com/gojue/ebpfmanager"
)
func OnEvent(cpu int, data []byte, perfmap *manager.PerfMap, manager *manager.Manager) {
	// Write your data handler here
}