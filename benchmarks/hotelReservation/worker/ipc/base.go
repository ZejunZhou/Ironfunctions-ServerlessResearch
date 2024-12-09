package ipc

import (
	"fmt"
	"os"
)

var rootPathForIpc string

const fileCreatMode = 0664

func SetRootPathForIpc(path string) {
	rootPathForIpc = path
}

func GetEngineUnixSocketPath() string {
	// create a unix socket file in the rootPathForIpc directory
	rootPathForIpc := fmt.Sprintf("%s/engine.sock", rootPathForIpc)
	if _, err := os.Stat(rootPathForIpc); os.IsNotExist(err) {
		if err := os.MkdirAll(rootPathForIpc, os.FileMode(0770)); err != nil {
			fmt.Printf("failed to create directory: %v\n", err)
		}
	}

	return rootPathForIpc
}

func GetFuncWorkerInputFifoName(clientId uint16) string {
	return fmt.Sprintf("worker_%d_input", clientId)
}

func GetFuncWorkerOutputFifoName(clientId uint16) string {
	return fmt.Sprintf("worker_%d_output", clientId)
}

func GetFuncCallInputShmName(fullCallId uint64) string {
	return fmt.Sprintf("%d.i", fullCallId)
}

func GetFuncCallOutputShmName(fullCallId uint64) string {
	return fmt.Sprintf("%d.o", fullCallId)
}

func GetFuncCallOutputFifoName(fullCallId uint64) string {
	return fmt.Sprintf("%d.o", fullCallId)
}
