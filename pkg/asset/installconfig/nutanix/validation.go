package nutanix

import (
	"context"

	"github.com/pkg/errors"
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

	_, err := nutanixtypes.CreateNutanixClient(context.TODO(),
		p.PrismCentral,
		p.Port,
		p.Username,
		p.Password,
		p.Insecure)

	if err != nil {
		return errors.New(field.InternalError(field.NewPath("platform", "nutanix"), errors.Wrapf(err, "unable to connect to Prism Central %s.", p.PrismCentral)).Error())
	}
	return validateResources(ic)
}

func validateResources(ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	// p := ic.Platform.Nutanix
	// allErrs = append(allErrs, validateNetwork(finder, p, field.NewPath("platform").Child("vsphere").Child("network"))...)
	return allErrs.ToAggregate()
}

// ValidateForProvisioning performs platform validation specifically for installer-
// provisioned infrastructure. In this case, self-hosted networking is a requirement
// when the installer creates infrastructure for vSphere clusters.
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
	return validateProvisioning(ic)
}

func validateProvisioning(ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validation.ValidateForProvisioning(ic.Platform.Nutanix, field.NewPath("platform").Child("nutanix"))...)
	// allErrs = append(allErrs, folderExists(finder, ic, field.NewPath("platform").Child("vsphere").Child("folder"))...)

	return allErrs.ToAggregate()
}

// func validateNetwork( p *nutanix.Platform, fldPath *field.Path) field.ErrorList {
// 	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
// 	defer cancel()
// 	dcName := p.Datacenter
// 	if !strings.HasPrefix(dcName, "/") && !strings.HasPrefix(dcName, "./") {
// 		dcName = "./" + dcName
// 	}

// 	dataCenter, err := finder.Datacenter(ctx, dcName)
// 	if err != nil {
// 		return field.ErrorList{field.Invalid(fldPath, p.Datacenter, err.Error())}
// 	}
// 	networkPath := fmt.Sprintf("%s/network/%s", dataCenter.InventoryPath, p.Network)
// 	_, err = finder.Network(ctx, networkPath)
// 	if err != nil {
// 		return field.ErrorList{field.Invalid(fldPath, p.Network, "unable to find network provided")}
// 	}
// 	return nil
// }
