package nutanix

import (
	"encoding/json"

	nutanixapis "github.com/nutanix-cloud-native/machine-api-provider-nutanix/pkg/apis/nutanixprovider/v1beta1"
	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	"github.com/pkg/errors"
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
	Insecure               bool   `json:"insecure"`
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
	Insecure              bool
	BootstrapIgnitionData string
	ClusterID             string
	ControlPlaneConfigs   []*nutanixapis.NutanixMachineProviderConfig
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
		Insecure:               sources.Insecure,
		PrismCentral:           sources.PrismCentral,
		Username:               sources.Username,
		Password:               sources.Password,
		MemoryMiB:              controlPlaneConfig.MemorySizeMib,
		DiskSizeMib:            controlPlaneConfig.DiskSizeMib,
		NumCPUs:                controlPlaneConfig.NumSockets,
		NumCoresPerSocket:      controlPlaneConfig.NumVcpusPerSocket,
		PrismElementUUID:       controlPlaneConfig.ClusterReferenceUUID,
		SubnetUUID:             controlPlaneConfig.SubnetUUID,
		Image:                  controlPlaneConfig.ImageName,
		ImageFilePath:          cachedImage,
		BootstrapIgnitionImage: bootstrapIgnitionImage,
	}
	return json.MarshalIndent(cfg, "", "  ")
}
