package exifSort

import (
	"fmt"
	"sync/atomic"
)

type walkStateType struct {
	skippedCount     uint64
	validCount       uint64
	invalidCount     uint64
	transferErrCount uint64
	walkDoPrint      bool
	walkErrMsgs      map[string]string
	transferErrMsgs  map[string]string
}

func (w *walkStateType) skipped() uint64 {
	return w.skippedCount
}

func (w *walkStateType) valid() uint64 {
	return w.validCount
}

func (w *walkStateType) invalid() uint64 {
	return w.invalidCount
}

func (w *walkStateType) walkErrs() map[string]string {
	return w.walkErrMsgs
}

func (w *walkStateType) transferErrs() map[string]string {
	return w.transferErrMsgs
}

func (w *walkStateType) total() uint64 {
	return w.skippedCount + w.validCount + w.invalidCount
}

func (w *walkStateType) storeValid() {
	atomic.AddUint64(&w.validCount, 1)
}

// We don't check if you have a path duplicate
func (w *walkStateType) storeInvalid(path string, errStr string) {
	atomic.AddUint64(&w.invalidCount, 1)
	w.walkErrMsgs[path] = errStr
}

// We don't check if you have a path duplicate
func (w *walkStateType) storeTransferErr(path string, errStr string) {
	atomic.AddUint64(&w.transferErrCount, 1)
	w.transferErrMsgs[path] = errStr
}

func (w *walkStateType) storeSkipped() {
	atomic.AddUint64(&w.skippedCount, 1)
}

// has to be a global so it can be accessed via walk routines
var walkState walkStateType

func (w *walkStateType) init(walkDoPrint bool) {
	w.skippedCount = 0
	w.validCount = 0
	w.invalidCount = 0
	w.walkDoPrint = walkDoPrint
	walkState.walkErrMsgs = make(map[string]string)
	walkState.transferErrMsgs = make(map[string]string)
}

func (w *walkStateType) walkPrintf(s string, params ...interface{}) {
	if !w.walkDoPrint {
		return
	}

	if len(params) == 0 {
		fmt.Print(s)
	}

	fmt.Printf(s, params...)
}

func walkErrMsg(path string, errMsg string) string {
	return fmt.Sprintf("%s with (%s)", path, errMsg)
}
