package nutanix

import (
	"context"

	"github.com/pkg/errors"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/nutanix/validation"
)

// Validate executes platform-specific validation.
func Validate(ic *types.InstallConfig) error {
	if ic.Platform.Nutanix == nil {
		return errors.New(field.Required(field.NewPath("platform", "nutanix"), "Nutanix validation requires a Nutanix platform configuration").Error())
	}

	p := ic.Platform.Nutanix
	if errs := validation.ValidatePlatform(p, field.NewPath("platform").Child("nutanix")); len(errs) != 0 {
		return errs.ToAggregate()
	}

	client, err := nutanixtypes.CreateNutanixClient(context.TODO(),
		p.PrismCentral,
		p.Port,
		p.Username,
		p.Password,
		p.Insecure)

	if err != nil {
		return errors.New(field.InternalError(field.NewPath("platform", "nutanix"), errors.Wrapf(err, "unable to connect to Prism Central %s.", p.PrismCentral)).Error())
	}
	return validateResources(client, ic)
}

func validateResources(client *v3.Client, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	return allErrs.ToAggregate()
}

// ValidateForProvisioning performs platform validation specifically for installer-
// provisioned infrastructure. In this case, self-hosted networking is a requirement
// when the installer creates infrastructure for nutanix clusters.
func ValidateForProvisioning(ic *types.InstallConfig) error {
	if ic.Platform.Nutanix == nil {
		return errors.New(field.Required(field.NewPath("platform", "nutanix"), "Nutanix validation requires a Nutanix platform configuration").Error())
	}

	p := ic.Platform.Nutanix
	_, err := nutanixtypes.CreateNutanixClient(context.TODO(),
		p.PrismCentral,
		p.Port,
		p.Username,
		p.Password,
		p.Insecure)

	if err != nil {
		return errors.New(field.InternalError(field.NewPath("platform", "nutanix"), errors.Wrapf(err, "unable to connect to Prism Central %s.", p.PrismCentral)).Error())
	}
	//TODO: add validation
	return nil
}
