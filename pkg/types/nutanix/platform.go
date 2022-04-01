package nutanix

import configv1 "github.com/openshift/api/config/v1"

// Platform stores any global configuration used for Nutanix platforms.
type Platform struct {
	// PrismCentral is the endpoint (address and port) and credentials to
	// connect to the Prism Central.
	PrismCentral PrismCentral `json:"prismCentral"`

	// PrismElement is the endpoint (address and port) and credentials to
	// connect to the Prism Elements (clusters). Currently we only support one
	// Prism Element (cluster) for an openshift cluster, where all the Nutanix resources (VMs, subnet, etc.)
	// used in the Openshift cluster locate. In the future, we may support the Nutanix resources (VMs, etc.)
	// used in the Openshift cluster can come from multiple Prism Elements (clusters) of the Prism Cental.
	PrismElements []PrismElement `json:"prismElements"`

	// ClusterOSImage overrides the url provided in rhcos.json to download the RHCOS Image
	//
	// +optional
	ClusterOSImage string `json:"clusterOSImage,omitempty"`

	// APIVIP is the virtual IP address for the api endpoint
	//
	// +kubebuilder:validation:format=ip
	// +optional
	APIVIP string `json:"apiVIP,omitempty"`

	// IngressVIP is the virtual IP address for ingress
	//
	// +kubebuilder:validation:format=ip
	// +optional
	IngressVIP string `json:"ingressVIP,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on Nutanix for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// SubnetUUID specifies the UUID of the subnet to be used by the cluster.
	SubnetUUID string `json:"subnetUUID,omitempty"`
}

// PrismCentral holds the endpoint and credentials data used to connect to the Prism Central
type PrismCentral struct {
	// Endpoint holds the address and port of the Prism Central
	Endpoint configv1.NutanixPrismEndpoint `json:"endpoint"`

	// Username is the name of the user to connect to the Prism Central
	Username string `json:"username"`

	// Password is the password for the user to connect to the Prism Central
	Password string `json:"password"`
}

// PrismElement holds the endpoint and credentials data used to connect to the Prism Element
type PrismElement struct {
	// UUID is a v4 UUID for Nutanix Prism element endpoint
	UUID string `json:"uuid"`

	// endpoint holds the address and port of the Prism Element
	Endpoint configv1.NutanixPrismElementEndpoint `json:"endpoint"`

	// Username is the name of the user to connect to the Prism Element
	Username string `json:"username"`

	// Password is the password for the user to connect to the Prism Element
	Password string `json:"password"`
}
