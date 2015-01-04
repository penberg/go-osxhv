package main

// #include <stdlib.h> // for valloc()
import (
	"C"
)

import (
	"fmt"
	. "github.com/penberg/go-osxhv"
	"os"
	"runtime"
)

func main() {
	runtime.LockOSThread()

	if err := HvVmCreate(); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to create VM: %v\n", err)
		os.Exit(1)
	}

	vmMemorySize := uint64(1 * 1024 * 1024)
	vmMemoryUva := uintptr(C.valloc(C.size_t(vmMemorySize)))
	vmMemoryGpa := uintptr(0)

	if err := HvVmMap(vmMemoryUva, vmMemoryGpa, vmMemorySize, 0); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to map VM memory: %v\n", err)
		os.Exit(1)
	}

	if err := HvVmProtect(vmMemoryGpa, vmMemorySize, HV_MEMORY_READ|HV_MEMORY_WRITE|HV_MEMORY_EXEC); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to map VM memory: %v\n", err)
		os.Exit(1)
	}

	vcpu, err := HvVCPUCreate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to create VCPU: %v\n", err)
		os.Exit(1)
	}

	for {
		err := vcpu.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: unable to create VCPU: %v\n", err)
			os.Exit(1)
		}
	}

	if err := HvVmDestroy(); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to destroy VM: %v\n", err)
		os.Exit(1)
	}
}
