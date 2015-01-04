package hv

// #cgo CFLAGS: -framework Hypervisor
// #cgo LDFLAGS: -framework Hypervisor
// #include <stddef.h>
// #include <Hypervisor/hv.h>
// #include <Hypervisor/hv_vmx.h>
import (
	"C"
)

import (
	"fmt"
	"unsafe"
)

const (
	HV_MEMORY_READ  = C.HV_MEMORY_READ
	HV_MEMORY_WRITE = C.HV_MEMORY_WRITE
	HV_MEMORY_EXEC  = C.HV_MEMORY_EXEC
)

type HvVCPU struct {
	vcpuid C.hv_vcpuid_t
}

func HvVmCreate() error {
	if err := C.hv_vm_create(C.HV_VM_DEFAULT); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

func HvVmDestroy() error {
	if err := C.hv_vm_destroy(); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

func HvVmMap(uva uintptr, gpa uintptr, size uint64, flags int) error {
	if err := C.hv_vm_map(C.hv_uvaddr_t(uva), C.hv_gpaddr_t(gpa), C.size_t(size), C.hv_memory_flags_t(flags)); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

func HvVmUnmap(gpa uintptr, size uint64) error {
	if err := C.hv_vm_unmap(C.hv_gpaddr_t(gpa), C.size_t(size)); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

func HvVmProtect(gpa uintptr, size uint64, flags int) error {
	if err := C.hv_vm_protect(C.hv_gpaddr_t(gpa), C.size_t(size), C.hv_memory_flags_t(flags)); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

func HvVmSyncTSC(tsc uint64) error {
	if err := C.hv_vm_sync_tsc(C.uint64_t(tsc)); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

func HvVCPUCreate() (*HvVCPU, error) {
	var vcpuid C.hv_vcpuid_t
	if err := C.hv_vcpu_create(&vcpuid, C.HV_VCPU_DEFAULT); err != C.HV_SUCCESS {
		return nil, hvError(err)
	}
	return &HvVCPU{
		vcpuid: vcpuid,
	}, nil
}

func (vcpu *HvVCPU) Destroy() error {
	if err := C.hv_vcpu_destroy(vcpu.vcpuid); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

func (vcpu *HvVCPU) ReadRegister(reg int) (uint64, error) {
	var value uint64
	if err := C.hv_vcpu_read_register(vcpu.vcpuid, C.hv_x86_reg_t(reg), (*C.uint64_t)(&value)); err != C.HV_SUCCESS {
		return 0, hvError(err)
	}
	return value, nil
}

func (vcpu *HvVCPU) WriteRegister(reg int, value uint64) error {
	if err := C.hv_vcpu_write_register(vcpu.vcpuid, C.hv_x86_reg_t(reg), C.uint64_t(value)); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

func (vcpu *HvVCPU) ReadFPState(buffer []byte) error {
	if err := C.hv_vcpu_read_fpstate(vcpu.vcpuid, unsafe.Pointer(&buffer[0]), C.size_t(len(buffer))); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

func (vcpu *HvVCPU) WriteFPState(buffer []byte) error {
	if err := C.hv_vcpu_write_fpstate(vcpu.vcpuid, unsafe.Pointer(&buffer[0]), C.size_t(len(buffer))); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

func (vcpu *HvVCPU) EnableNativeMSR(msr uint32, enable bool) error {
	if err := C.hv_vcpu_enable_native_msr(vcpu.vcpuid, C.uint32_t(msr), C.bool(enable)); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

func (vcpu *HvVCPU) ReadMSR(msr uint32) (uint64, error) {
	var value uint64
	if err := C.hv_vcpu_read_msr(vcpu.vcpuid, C.uint32_t(msr), (*C.uint64_t)(&value)); err != C.HV_SUCCESS {
		return ^uint64(0), hvError(err)
	}
	return value, nil
}

func (vcpu *HvVCPU) WriteMSR(msr uint32, value uint64) error {
	if err := C.hv_vcpu_write_msr(vcpu.vcpuid, C.uint32_t(msr), C.uint64_t(value)); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

func (vcpu *HvVCPU) Flush() error {
	if err := C.hv_vcpu_flush(vcpu.vcpuid); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

func (vcpu *HvVCPU) InvalidateTLB() error {
	if err := C.hv_vcpu_invalidate_tlb(vcpu.vcpuid); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

func (vcpu *HvVCPU) Run() error {
	if err := C.hv_vcpu_run(vcpu.vcpuid); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

func HvVCPUInterrupt(vcpus []HvVCPU) error {
	vcpuids := make([]C.hv_vcpuid_t, len(vcpus))
	for i, vcpu := range vcpus {
		vcpuids[i] = vcpu.vcpuid
	}
	if err := C.hv_vcpu_interrupt((*C.hv_vcpuid_t)(&vcpuids[0]), C.uint(len(vcpuids))); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

func (vcpu *HvVCPU) GetExecTime() (uint64, error) {
	var time uint64
	if err := C.hv_vcpu_get_exec_time(vcpu.vcpuid, (*C.uint64_t)(&time)); err != C.HV_SUCCESS {
		return ^uint64(0), hvError(err)
	}
	return time, nil
}

func hvError(err C.hv_return_t) error {
	return fmt.Errorf("%s (%d)", hvErrorString(err), err)
}

func hvErrorString(err C.hv_return_t) string {
	switch err {
	case C.HV_SUCCESS:
		return "Success"
	case C.HV_ERROR:
		return "Error"
	case C.HV_BUSY:
		return "Busy"
	case C.HV_BAD_ARGUMENT:
		return "Bad argument"
	case C.HV_NO_RESOURCES:
		return "No resources"
	case C.HV_NO_DEVICE:
		return "No device"
	case C.HV_UNSUPPORTED:
		return "Unsupported"
	default:
		return "Unknown"
	}
}
