// Package nutanix generates Machine objects for nutanix.
package nutanix

import (
	"fmt"

	baremetalprovider "github.com/metal3-io/cluster-api-provider-baremetal/pkg/apis/baremetal/v1alpha1"
	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != nutanix.Name {
		return nil, fmt.Errorf("non nutanix configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != nutanix.Name {
		return nil, fmt.Errorf("non-nutanix machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.Nutanix

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		provider, err := provider(clusterID, platform, osImage, userDataSecret)

		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
		}
		machine := machineapi.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "machine.openshift.io/v1beta1",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "openshift-machine-api",
				Name:      fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, idx),
				Labels: map[string]string{
					"machine.openshift.io/cluster-api-cluster":      clusterID,
					"machine.openshift.io/cluster-api-machine-role": role,
					"machine.openshift.io/cluster-api-machine-type": role,
				},
			},
			Spec: machineapi.MachineSpec{
				ProviderSpec: machineapi.ProviderSpec{
					Value: &runtime.RawExtension{Object: provider},
				},
				// we don't need to set Versions, because we control those via operators.
			},
		}
		machines = append(machines, machine)
	}
	return machines, nil
}

func provider(clusterID string, platform *nutanix.Platform, osImage string, userDataSecret string) (*baremetalprovider.BareMetalMachineProviderSpec, error) {

	return &baremetalprovider.BareMetalMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "baremetal.cluster.k8s.io/v1alpha1",
			Kind:       "BareMetalMachineProviderSpec",
		},
		Image:    baremetalprovider.Image{},
		UserData: &corev1.SecretReference{Name: userDataSecret},
	}, nil
}

// func provider(clusterID string, platform *nutanix.Platform, mpool *nutanix.MachinePool, osImage string, userDataSecret string) (*nutanixapis.NutanixMachineSpec, error) {
// 	return nil, nil
// }

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, clusterID string) {
}
