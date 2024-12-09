package faas

import (
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/harlow/go-micro-services/worker/config"
	"github.com/harlow/go-micro-services/worker/ipc"
	"github.com/harlow/go-micro-services/worker/types"
	"github.com/harlow/go-micro-services/worker/worker"
)

func Serve(factory types.FuncHandlerFactory) {
	runtime.GOMAXPROCS(1)
	ipc.SetRootPathForIpc(os.Getenv("FAAS_ROOT_PATH_FOR_IPC"))
	funcId, err := strconv.Atoi(os.Getenv("FAAS_FUNC_ID"))
	if err != nil {
		log.Fatal("[FATAL] Failed to parse FAAS_FUNC_ID")
	}
	clientId, err := strconv.Atoi(os.Getenv("FAAS_CLIENT_ID"))
	if err != nil {
		log.Fatal("[FATAL] Failed to parse FAAS_CLIENT_ID")
	}

	log.Printf("[INFO] Starting worker for funcId: %d, clientId: %d\n", funcId, clientId)
	log.Printf("[INFO] FAAS_ROOT_PATH_FOR_IPC: %s\n", os.Getenv("FAAS_ROOT_PATH_FOR_IPC"))

	// msgPipeFd, err := strconv.Atoi(os.Getenv("FAAS_MSG_PIPE_FD"))
	// if err != nil {
	// 	log.Fatal("[FATAL] Failed to parse FAAS_MSG_PIPE_FD")
	// }

	// msgPipe := os.NewFile(uintptr(msgPipeFd), "msg_pipe")
	// payloadSizeBuf := make([]byte, 4)
	// nread, err := msgPipe.Read(payloadSizeBuf)
	// if err != nil || nread != len(payloadSizeBuf) {
	// 	log.Fatalf("[FATAL] Failed to read payload size, nread: %d, len(payloadSizeBuf): %d, err: %v\n", nread, len(payloadSizeBuf), err)
	// }

	// payloadSize := binary.LittleEndian.Uint32(payloadSizeBuf)
	// payload := make([]byte, payloadSize)
	// nread, err = msgPipe.Read(payload)
	// if err != nil || nread != len(payload) {
	// 	log.Fatal("[FATAL] Failed to read payload")
	// }

	// err = config.InitFuncConfig(payload)
	err = config.InitFuncConfig()
	if err != nil {
		log.Fatalf("[FATAL] InitFuncConfig failed: %s", err)
	}

	w, err := worker.NewFuncWorker(uint16(funcId), uint16(clientId), factory)
	if err != nil {
		log.Fatal("[FATAL] Failed to create FuncWorker: ", err)
	}

	go func(w *worker.FuncWorker) {
		w.Run()
	}(w)

	numWorkers := 1
	maxProcFactor, err := strconv.Atoi(os.Getenv("FAAS_GO_MAX_PROC_FACTOR"))
	if err != nil {
		maxProcFactor = 8
	}

	for {
		// message := protocol.NewEmptyMessage()
		// nread, err := msgPipe.Read(message)
		// if err != nil || (nread != protocol.MessageFullByteSize && nread != 0) {
		// 	log.Fatalf("[FATAL] Failed to read message, nread: %d, len(message): %d, err: %v\n", nread, protocol.MessageFullByteSize, err)
		// }
		// if protocol.IsCreateFuncWorkerMessage(message) {
		// clientId := protocol.GetClientIdFromMessage(message)
		w, err := worker.NewFuncWorker(uint16(funcId), uint16(clientId+1), factory)
		if err != nil {
			log.Fatal("[FATAL] Failed to create FuncWorker: ", err)
		}
		numWorkers += 1
		planedMaxProcs := (numWorkers-1)/maxProcFactor + 1
		currentMaxProcs := runtime.GOMAXPROCS(0)
		if planedMaxProcs > currentMaxProcs {
			runtime.GOMAXPROCS(planedMaxProcs)
			log.Printf("[INFO] Current GOMAXPROCS is %d", planedMaxProcs)
		}
		go func(w *worker.FuncWorker) {
			w.Run()
		}(w)
	}
	//  else {
	// 	log.Fatal("[FATAL] Unknown message type")
	// }
	// }
}
