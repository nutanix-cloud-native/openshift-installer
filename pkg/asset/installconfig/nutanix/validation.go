package nutanix

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
)

// Validate executes platform-specific validation.
func Validate(ic *types.InstallConfig) error {
	if ic.Platform.Nutanix == nil {
		return field.Required(field.NewPath("platform", "nutanix"), "nutanix validation requires a nutanix platform configuration")
	}

	p := ic.Platform.Nutanix
	nc, err := nutanixtypes.CreateNutanixClient(context.TODO(), p.PrismCentral.Endpoint.Address, strconv.Itoa(int(p.PrismCentral.Endpoint.Port)), p.PrismCentral.Username, p.PrismCentral.Password)
	if err != nil {
		return field.InternalError(field.NewPath("platform", "nutanix"), errors.Wrapf(err, "unable to connect to Prism Central %q", p.PrismCentral.Endpoint.Address))
	}

	// validate whether a prism element with the UUID actually exists
	_, err = nc.V3.GetCluster(p.PrismElements[0].UUID)
	if err != nil {
		return field.InternalError(field.NewPath("platform", "nutanix", "prismElementUUID"), errors.Wrapf(err, "prism element UUID %s does not correspond to a valid prism element in Prism", p.PrismElements[0].UUID))
	}

	return nil
}

// ValidateForProvisioning performs platform validation specifically for installer-
// provisioned infrastructure. In this case, self-hosted networking is a requirement
// when the installer creates infrastructure for nutanix clusters.
func ValidateForProvisioning(ic *types.InstallConfig) error {
	if ic.Platform.Nutanix == nil {
		return field.Required(field.NewPath("platform", "nutanix"), "nutanix validation requires a nutanix platform configuration")
	}

	p := ic.Platform.Nutanix
	nc, err := nutanixtypes.CreateNutanixClient(context.TODO(), p.PrismCentral.Endpoint.Address, strconv.Itoa(int(p.PrismCentral.Endpoint.Port)), p.PrismCentral.Username, p.PrismCentral.Password)
	if err != nil {
		return field.InternalError(field.NewPath("platform", "nutanix"), errors.Wrapf(err, "unable to connect to Prism Central %q", p.PrismCentral.Endpoint.Address))
	}

	// validate whether a subnet with the UUID actually exists
	_, err = nc.V3.GetSubnet(p.SubnetUUID)
	if err != nil {
		return field.InternalError(field.NewPath("platform", "nutanix", "subnetUUID"), errors.Wrapf(err, "subnet UUID %s does not correspond to a valid subnet in Prism", p.SubnetUUID))
	}

	return nil
}
