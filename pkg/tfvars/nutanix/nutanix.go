package nutanix

import (
	"encoding/json"
	"fmt"
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
	DiskGiB                int32  `json:"nutanix_control_plane_disk_gib"`
	NumCPUs                int32  `json:"nutanix_control_plane_num_cpus"`
	NumCoresPerSocket      int32  `json:"nutanix_control_plane_cores_per_socket"`
	PrismElement           string `json:"nutanix_prism_element"`
	Insecure               bool   `json:"insecure"`
	Subnet                 string `json:"nutanix_subnet"`
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
	PrismElement          string
	ImageURL              string
	Insecure              bool
	BootstrapIgnitionData string
	ClusterID             string
	//TODO: fetch from control plane configs
	Subnet string
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
	cfg := &config{
		Port:                   sources.Port,
		Insecure:               sources.Insecure,
		PrismCentral:           sources.PrismCentral,
		Username:               sources.Username,
		Password:               sources.Password,
		MemoryMiB:              16384, //TODO: need to fetch from control plane config
		DiskGiB:                120,   //TODO: need to fetch from control plane config
		NumCPUs:                4,     //TODO: need to fetch from control plane config
		NumCoresPerSocket:      4,     //TODO: need to fetch from control plane config
		PrismElement:           sources.PrismElement,
		Subnet:                 sources.Subnet, //TODO: need to fetch from control plane config
		Image:                  generateImageName(sources.ClusterID), //"rhcos-manual", //TODO: need to fetch from control plane config
		ImageFilePath:          cachedImage,
		BootstrapIgnitionImage: bootstrapIgnitionImage,
	}
	return json.MarshalIndent(cfg, "", "  ")
}

func generateImageName(infraID string ) string {
	return fmt.Sprintf("%s-rhcos",infraID)
}
