package nutanix

import (
	"encoding/json"

	machinev1 "github.com/openshift/api/machine/v1"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
)

type config struct {
	PrismCentral           string `json:"prism_central"`
	Port                   string `json:"port"`
	Username               string `json:"username"`
	Password               string `json:"password"`
	MemoryMiB              int64  `json:"nutanix_control_plane_memory_mib"`
	DiskSizeMib            int64  `json:"nutanix_control_plane_disk_mib"`
	NumCPUs                int64  `json:"nutanix_control_plane_num_cpus"`
	NumCoresPerSocket      int64  `json:"nutanix_control_plane_cores_per_socket"`
	PrismElementUUID       string `json:"nutanix_prism_element_uuid"`
	SubnetUUID             string `json:"nutanix_subnet_uuid"`
	Image                  string `json:"nutanix_image"`
	ImageFilePath          string `json:"nutanix_image_filepath"`
	BootstrapIgnitionImage string `json:"nutanix_bootstrap_ignition_image"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	PrismCentral          string
	Port                  string
	Username              string
	Password              string
	ImageURL              string
	BootstrapIgnitionData string
	ClusterID             string
	ControlPlaneConfigs   []*machinev1.NutanixMachineProviderConfig
}

//TFVars generate Nutanix-specific Terraform variables
func TFVars(sources TFVarsSources) ([]byte, error) {
	cachedImage, err := cache.DownloadImageFile(sources.ImageURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to use cached nutanix image")
	}
	bootstrapIgnitionImage, err := nutanixtypes.CreateBootstrapISO(sources.ClusterID, sources.BootstrapIgnitionData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bootstrap ignition iso")
	}
	controlPlaneConfig := sources.ControlPlaneConfigs[0]
	cfg := &config{
		Port:                   sources.Port,
		PrismCentral:           sources.PrismCentral,
		Username:               sources.Username,
		Password:               sources.Password,
		MemoryMiB:              controlPlaneConfig.MemorySize.Value() / (1024 * 1024),
		DiskSizeMib:            controlPlaneConfig.SystemDiskSize.Value() / (1024 * 1024),
		NumCPUs:                int64(controlPlaneConfig.VCPUSockets),
		NumCoresPerSocket:      int64(controlPlaneConfig.VCPUsPerSocket),
		PrismElementUUID:       *controlPlaneConfig.Cluster.UUID,
		SubnetUUID:             *controlPlaneConfig.Subnet.UUID,
		Image:                  *controlPlaneConfig.Image.Name,
		ImageFilePath:          cachedImage,
		BootstrapIgnitionImage: bootstrapIgnitionImage,
	}
	return json.MarshalIndent(cfg, "", "  ")
}
