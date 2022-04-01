package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
// TODO(nutanix): Revisit for further expanding the validation logic
func ValidatePlatform(p *nutanix.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, validatePrismCentral(p, fldPath)...)
	allErrs = append(allErrs, validatePrismElements(p, fldPath)...)

	if len(p.SubnetUUID) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("subnet"), "must specify the subnet"))
	}

	// If all VIPs are empty, skip IP validation.  All VIPs are required to be defined together.
	if p.APIVIP != "" || p.IngressVIP != "" {
		allErrs = append(allErrs, validateVIPs(p, fldPath)...)
	}

	return allErrs
}

func validatePrismCentral(p *nutanix.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(p.PrismCentral.Endpoint.Address) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismCentral").Child("endpoint").Child("address"),
			"must specify the Prism Central endpoint address"))
	} else {
		if err := validate.Host(p.PrismCentral.Endpoint.Address); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("prismCentral").Child("endpoint").Child("address"),
				p.PrismCentral.Endpoint.Address, "must be the domain name or IP address of the Prism Central"))
		}
	}

	if p.PrismCentral.Endpoint.Port < 1 || p.PrismCentral.Endpoint.Port > 65535 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("prismCentral").Child("endpoint").Child("port"),
			p.PrismCentral.Endpoint.Port, "The Prism Central endpoint port is invalid, must be in the range of 1 to 65535"))
	}

	if len(p.PrismCentral.Username) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismCentral").Child("username"),
			"must specify the Prism Central username"))
	}

	if len(p.PrismCentral.Password) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismCentral").Child("password"),
			"must specify the Prism Central password"))
	}

	return allErrs
}

// validateVIPs checks that all required VIPs are provided and are valid IP addresses.
func validateVIPs(p *nutanix.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(p.APIVIP) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("apiVIP"), "must specify a VIP for the API"))
	} else if err := validate.IP(p.APIVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, err.Error()))
	}

	if len(p.IngressVIP) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("ingressVIP"), "must specify a VIP for Ingress"))
	} else if err := validate.IP(p.IngressVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, err.Error()))
	}

	if p.APIVIP == p.IngressVIP {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, "IPs for both API and Ingress should not be the same"))
	}

	return allErrs
}

func validatePrismElements(p *nutanix.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(p.PrismElements[0].Endpoint.Name) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismElements[0]").Child("endpoint").Child("name"),
			"must specify the Prism Element (cluster) name"))
	}

	if len(p.PrismElements[0].Endpoint.Endpoint.Address) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismElements[0]").Child("endpoint").Child("endpoint").Child("address"),
			"must specify the Prism element (cluster) endpoint address"))
	} else {
		if err := validate.Host(p.PrismElements[0].Endpoint.Endpoint.Address); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("prismElements[0]").Child("endpoint").Child("endpoint").Child("address"),
				p.PrismElements[0].Endpoint.Endpoint.Address, "must be the domain name or IP address of the Prism Element (cluster)"))
		}
	}

	if p.PrismElements[0].Endpoint.Endpoint.Port < 1 || p.PrismElements[0].Endpoint.Endpoint.Port > 65535 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismElements[0]").Child("endpoint").Child("port"),
			"The Prism Element (cluster) endpoint port is invalid, must be in the range of 1 to 65535"))
	}

	if len(p.PrismElements[0].Username) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismElements[0]").Child("username"),
			"must specify the Prism Element (cluster) username"))
	}

	if len(p.PrismElements[0].Password) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismElements[0]").Child("password"),
			"must specify the Prism Element (cluster) password"))
	}

	return allErrs
}
