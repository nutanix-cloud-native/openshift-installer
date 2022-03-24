package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types/nutanix"
)

func validPlatform() *nutanix.Platform {
	return &nutanix.Platform{
		PrismCentral: nutanix.NutanixPrismCentral{
			Endpoint: configv1.NutanixPrismEndpoint{Address: "test-pc", Port: 8080},
			Username: "test-username-pc",
			Password: "test-password-pc",
		},
		PrismElements: []nutanix.NutanixPrismElement{{
			Name:     "test-pe-name",
			UUID:     "test-pe-uuid",
			Endpoint: configv1.NutanixPrismEndpoint{Address: "test-pe", Port: 8081},
			Username: "test-username-pe",
			Password: "test-password-pe",
		}},
		SubnetUUID: "test-subnet",
	}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name          string
		platform      *nutanix.Platform
		expectedError string
	}{
		{
			name:     "minimal",
			platform: validPlatform(),
		},
		{
			name: "missing Prism Central address",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral.Endpoint.Address = ""
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.endpoint\.address: Required value: must specify the Prism Central endpoint address$`,
		},
		{
			name: "missing Prism Central username",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral.Username = ""
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.username: Required value: must specify the Prism Central username$`,
		},
		{
			name: "missing password",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral.Password = ""
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.password: Required value: must specify the Prism Central password$`,
		},
		{
			name: "missing prism element name",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismElements[0].Name = ""
				return p
			}(),
			expectedError: `^test-path\.prismElement[0]\.name: Required value: must specify the Prism Element (cluster) name$`,
		},
		{
			name: "missing prism element address",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismElements[0].Endpoint.Address = ""
				return p
			}(),
			expectedError: `^test-path\.prismElement[0]\.endpoint\.addres: Required value: must specify the Prism element (cluster) endpoint address$`,
		},
		{
			name: "valid VIPs",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.2"
				p.IngressVIP = "192.168.111.3"
				return p
			}(),
		},
		{
			name: "missing API VIP",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.APIVIP = ""
				p.IngressVIP = "192.168.111.3"
				return p
			}(),
			expectedError: `^test-path\.apiVIP: Required value: must specify a VIP for the API$`,
		},
		{
			name: "missing Ingress VIP",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.2"
				p.IngressVIP = ""
				return p
			}(),
			expectedError: `^test-path\.ingressVIP: Required value: must specify a VIP for Ingress$`,
		},
		{
			name: "Invalid API VIP",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111"
				p.IngressVIP = "192.168.111.2"
				return p
			}(),
			expectedError: `^test-path\.apiVIP: Invalid value: "192.168.111": "192.168.111" is not a valid IP$`,
		},
		{
			name: "Invalid Ingress VIP",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.1"
				p.IngressVIP = "192.168.111"
				return p
			}(),
			expectedError: `^test-path\.ingressVIP: Invalid value: "192.168.111": "192.168.111" is not a valid IP$`,
		},
		{
			name: "Same API and Ingress VIP",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.APIVIP = "192.168.111.1"
				p.IngressVIP = "192.168.111.1"
				return p
			}(),
			expectedError: `^test-path\.apiVIP: Invalid value: "192.168.111.1": IPs for both API and Ingress should not be the same$`,
		},
		{
			name: "Capital letters in Prism Central address",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral.Endpoint.Address = "tEsT-PrismCentral"
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.endpoint\.address: Invalid value: "tEsT-PrismCentral": must be the domain name or IP address of the Prism Central$`,
		},
		{
			name: "URL as Prism Central",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral = "https://test-pc"
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.endpoint\.address: Invalid value: "https://test-pc": must be the domain name or IP address of the Prism Central$`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, field.NewPath("test-path")).ToAggregate()
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expectedError, err)
			}
		})
	}
}
