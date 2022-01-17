package nutanix

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in Nutanix.
var PlatformStages = []terraform.Stage{
	stages.NewStage("nutanix", "pre-bootstrap"),
	stages.NewStage("nutanix", "bootstrap", stages.WithNormalBootstrapDestroy()),
	stages.NewStage("nutanix", "master"),
}
