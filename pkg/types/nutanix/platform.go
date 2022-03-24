package nutanix

import (
	configv1 "github.com/openshift/api/config/v1"
)

// Platform stores any global configuration used for Nutanix platforms.
type Platform struct {
	// PrismCentral is the endpoint (address and port) and credentials to
	// connect to the Prism Central.
	PrismCentral NutanixPrismCentral `json:"prismCentral"`

	// PrismElement is the endpoint (address and port) and credentials to
	// connect to the Prism Elements (clusters). Currently we only support one
	// Prism Element (cluster) for an openshift cluster, where all the Nutanix resources (VMs, subnet, etc.)
	// used in the Openshift cluster locate. In the future, we may support the Nutanix resources (VMs, etc.)
	// used in the Openshift cluster can come from multiple Prism Elements (clusters) of the Prism Cental.
	PrismElements []NutanixPrismElement `json:"prismElements"`

	// Insecure disables certificate checking when connecting to Prism Central.
	Insecure bool `json:"insecure"`

	// DefaultStorageContainer is the default datastore to use for provisioning volumes.
	// DefaultStorageContainer string `json:"defaultStorageContainer"`

	// ClusterOSImage overrides the url provided in rhcos.json to download the RHCOS Image
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

	// SubnetUUID identifies the network subnet to be used by the cluster.
	SubnetUUID string `json:"subnetUuid,omitempty"`
}

// NutanixPrismCentral holds the endpoint and credentials data used to connect to the Prism Central
type NutanixPrismCentral struct {
	// Endpoint holds the address and port of the Prism Central
	Endpoint configv1.NutanixPrismEndpoint `json:"endpoint"`

	// Username is the name of the user to connect to the Prism Central
	Username string `json:"username"`

	// Password is the password for the user to connect to the Prism Central
	Password string `json:"password"`
}

//
type NutanixPrismElement struct {
	// name is the name of the Prism Element (cluster)
	Name string `json:"name"`

	// uuid
	UUID string `json:"uuid"`

	// endpoint holds the address and port of the Prism Element
	Endpoint configv1.NutanixPrismEndpoint `json:"endpoint"`

	// Username is the name of the user to connect to the Prism Element
	Username string `json:"username"`

	// Password is the password for the user to connect to the Prism Element
	Password string `json:"password"`
}
