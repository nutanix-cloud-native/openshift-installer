package nutanix

// MachinePool stores the configuration for a machine pool installed
// on Nutanix.
type MachinePool struct {
	// NumCPUs is the total number of virtual processor cores to assign a vm.
	//
	// +optional
	NumCPUs int64 `json:"cpus"`

	// NumCoresPerSocket is the number of cores per socket in a vm. The number
	// of vCPUs on the vm will be NumCPUs/NumCoresPerSocket.
	//
	// +optional
	NumCoresPerSocket int64 `json:"coresPerSocket"`

	// Memory is the size of a VM's memory in MB.
	//
	// +optional
	MemoryMiB int64 `json:"memoryMB"`

	// OSDisk defines the storage for instance.
	//
	// +optional
	OSDisk `json:"osDisk"`
}

// OSDisk defines the disk for a virtual machine.
type OSDisk struct {
	// DiskSizeMib defines the size of disk in MiB.
	//
	// +optional
	DiskSizeMib int64 `json:"diskSizeMib"`
}

// Set sets the values from `required` to `p`.
func (p *MachinePool) Set(required *MachinePool) {
	if required == nil || p == nil {
		return
	}

	if required.NumCPUs != 0 {
		p.NumCPUs = required.NumCPUs
	}

	if required.NumCoresPerSocket != 0 {
		p.NumCoresPerSocket = required.NumCoresPerSocket
	}

	if required.MemoryMiB != 0 {
		p.MemoryMiB = required.MemoryMiB
	}

	if required.OSDisk.DiskSizeMib != 0 {
		p.OSDisk.DiskSizeMib = required.OSDisk.DiskSizeMib
	}
}