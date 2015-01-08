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

const (
	HV_X86_RAX       = C.HV_X86_RAX
	HV_X86_RCX       = C.HV_X86_RCX
	HV_X86_RDX       = C.HV_X86_RDX
	HV_X86_RBX       = C.HV_X86_RBX
	HV_X86_RSP       = C.HV_X86_RSP
	HV_X86_RBP       = C.HV_X86_RBP
	HV_X86_RSI       = C.HV_X86_RSI
	HV_X86_RDI       = C.HV_X86_RDI
	HV_X86_R8        = C.HV_X86_R8
	HV_X86_R9        = C.HV_X86_R9
	HV_X86_R10       = C.HV_X86_R10
	HV_X86_R11       = C.HV_X86_R11
	HV_X86_R12       = C.HV_X86_R12
	HV_X86_R13       = C.HV_X86_R13
	HV_X86_R14       = C.HV_X86_R14
	HV_X86_R15       = C.HV_X86_R15
	HV_X86_CS        = C.HV_X86_CS
	HV_X86_SS        = C.HV_X86_SS
	HV_X86_DS        = C.HV_X86_DS
	HV_X86_ES        = C.HV_X86_ES
	HV_X86_FS        = C.HV_X86_FS
	HV_X86_GS        = C.HV_X86_GS
	HV_X86_RIP       = C.HV_X86_RIP
	HV_X86_RFLAGS    = C.HV_X86_RFLAGS
	HV_X86_GDT_BASE  = C.HV_X86_GDT_BASE
	HV_X86_GDT_LIMIT = C.HV_X86_GDT_LIMIT
	HV_X86_IDT_BASE  = C.HV_X86_IDT_BASE
	HV_X86_IDT_LIMIT = C.HV_X86_IDT_LIMIT
	HV_X86_LDTR      = C.HV_X86_LDTR
	HV_X86_LDT_BASE  = C.HV_X86_LDT_BASE
	HV_X86_LDT_LIMIT = C.HV_X86_LDT_LIMIT
	HV_X86_LDT_AR    = C.HV_X86_LDT_AR
	HV_X86_TR        = C.HV_X86_TR
	HV_X86_TSS_BASE  = C.HV_X86_TSS_BASE
	HV_X86_TSS_LIMIT = C.HV_X86_TSS_LIMIT
	HV_X86_TSS_AR    = C.HV_X86_TSS_AR
	HV_X86_CR0       = C.HV_X86_CR0
	HV_X86_CR1       = C.HV_X86_CR1
	HV_X86_CR2       = C.HV_X86_CR2
	HV_X86_CR3       = C.HV_X86_CR3
	HV_X86_CR4       = C.HV_X86_CR4
	HV_X86_DR0       = C.HV_X86_DR0
	HV_X86_DR1       = C.HV_X86_DR1
	HV_X86_DR2       = C.HV_X86_DR2
	HV_X86_DR3       = C.HV_X86_DR3
	HV_X86_DR4       = C.HV_X86_DR4
	HV_X86_DR5       = C.HV_X86_DR5
	HV_X86_DR6       = C.HV_X86_DR6
	HV_X86_DR7       = C.HV_X86_DR7
	HV_X86_TPR       = C.HV_X86_TPR
	HV_X86_XCR0      = C.HV_X86_XCR0
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

func (vcpu *HvVCPU) VmxReadVMCS(field uint32) (uint64, error) {
	var value uint64
	if err := C.hv_vmx_vcpu_read_vmcs(vcpu.vcpuid, (C.uint32_t)(field), (*C.uint64_t)(&value)); err != C.HV_SUCCESS {
		return ^uint64(0), hvError(err)
	}
	return value, nil
}

func (vcpu *HvVCPU) VmxWriteVMCS(field uint32, value uint64) error {
	if err := C.hv_vmx_vcpu_write_vmcs(vcpu.vcpuid, (C.uint32_t)(field), (C.uint64_t)(value)); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
}

const (
	HV_VMX_CAP_PINBASED   = C.HV_VMX_CAP_PINBASED
	HV_VMX_CAP_PROCBASED  = C.HV_VMX_CAP_PROCBASED
	HV_VMX_CAP_PROCBASED2 = C.HV_VMX_CAP_PROCBASED2
	HV_VMX_CAP_ENTRY      = C.HV_VMX_CAP_ENTRY
)

func HvVmxReadCapability(field int) (uint64, error) {
	var value uint64
	if err := C.hv_vmx_read_capability(C.hv_vmx_capability_t(field), (*C.uint64_t)(&value)); err != C.HV_SUCCESS {
		return ^uint64(0), hvError(err)
	}
	return value, nil
}

func (vcpu *HvVCPU) VmxSetAPICAddress(gpa uintptr) error {
	if err := C.hv_vmx_vcpu_set_apic_address(vcpu.vcpuid, (C.hv_gpaddr_t)(gpa)); err != C.HV_SUCCESS {
		return hvError(err)
	}
	return nil
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
