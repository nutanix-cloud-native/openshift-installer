package nutanix

import (
	"encoding/json"
	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	nutanixapis "github.com/openshift/machine-api-provider-nutanix/pkg/apis/nutanixprovider/v1beta1"
	"github.com/pkg/errors"
	"os"
)

const (
	nutanixOSImageOverrideEnvVar                = "NUTANIX_OS_IMAGE_URI"
	nutanixBootstrapIgnitionImageOverrideEnvVar = "NUTANIX_IGNITION_IMAGE_URI"
)

type config struct {
	PrismCentralAddress            string `json:"nutanix_prism_central_address"`
	Port                           string `json:"nutanix_prism_central_port"`
	Username                       string `json:"nutanix_username"`
	Password                       string `json:"nutanix_password"`
	MemoryMiB                      int64  `json:"nutanix_control_plane_memory_mib"`
	DiskSizeMiB                    int64  `json:"nutanix_control_plane_disk_mib"`
	NumCPUs                        int64  `json:"nutanix_control_plane_num_cpus"`
	NumCoresPerSocket              int64  `json:"nutanix_control_plane_cores_per_socket"`
	PrismElementUUID               string `json:"nutanix_prism_element_uuid"`
	SubnetUUID                     string `json:"nutanix_subnet_uuid"`
	Image                          string `json:"nutanix_image"`
	ImageFilePath                  string `json:"nutanix_image_filepath,omitempty"`
	ImageURI                       string `json:"nutanix_image_uri,omitempty"`
	BootstrapIgnitionImage         string `json:"nutanix_bootstrap_ignition_image"`
	BootstrapIgnitionImageFilePath string `json:"nutanix_bootstrap_ignition_image_filepath,omitempty"`
	BootstrapIgnitionImageURI      string `json:"nutanix_bootstrap_ignition_image_uri,omitempty"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	PrismCentralAddress   string
	Port                  string
	Username              string
	Password              string
	ImageURL              string
	BootstrapIgnitionData string
	ClusterID             string
	ControlPlaneConfigs   []*nutanixapis.NutanixMachineProviderConfig
}

//TFVars generate Nutanix-specific Terraform variables
func TFVars(sources TFVarsSources) ([]byte, error) {
	controlPlaneConfig := sources.ControlPlaneConfigs[0]
	bootstrapIgnitionImageName := nutanixtypes.BootISOImageName(sources.ClusterID)
	cfg := &config{
		Port:                   sources.Port,
		PrismCentralAddress:    sources.PrismCentralAddress,
		Username:               sources.Username,
		Password:               sources.Password,
		MemoryMiB:              controlPlaneConfig.MemorySizeMib,
		DiskSizeMiB:            controlPlaneConfig.DiskSizeMib,
		NumCPUs:                controlPlaneConfig.NumSockets,
		NumCoresPerSocket:      controlPlaneConfig.NumVcpusPerSocket,
		PrismElementUUID:       controlPlaneConfig.ClusterReferenceUUID,
		SubnetUUID:             controlPlaneConfig.SubnetUUID,
		Image:                  controlPlaneConfig.ImageName,
		BootstrapIgnitionImage: bootstrapIgnitionImageName,
	}

	osImageOverride := os.Getenv(nutanixOSImageOverrideEnvVar)
	if osImageOverride != "" {
		cfg.ImageURI = osImageOverride
	} else {
		cachedImage, err := cache.DownloadImageFile(sources.ImageURL)
		if err != nil {
			return nil, errors.Wrap(err, "failed to use cached nutanix image")
		}
		cfg.ImageFilePath = cachedImage
	}

	bootstrapIgnitionImagePath, err := nutanixtypes.CreateBootstrapISO(sources.ClusterID, sources.BootstrapIgnitionData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bootstrap ignition iso")
	}

	ignitionImageOverride := os.Getenv(nutanixBootstrapIgnitionImageOverrideEnvVar)
	if ignitionImageOverride != "" {
		cfg.BootstrapIgnitionImageURI = ignitionImageOverride
	} else {
		cfg.BootstrapIgnitionImageFilePath = bootstrapIgnitionImagePath
	}

	return json.MarshalIndent(cfg, "", "  ")
}
