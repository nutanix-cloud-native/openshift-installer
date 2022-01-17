package nutanix

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	nutanixclientv3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/openshift/installer/pkg/destroy/providers"
	installertypes "github.com/openshift/installer/pkg/types"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
)

const (
	emptyFilter                 = ""
	expectedCategoryKeyFormat   = "kubernetes-io-cluster-%s"
	expectedCategoryValueOwned  = "owned"
	expectedCategoryValueShared = "shared"
)

// clusterUninstaller holds the various options for the cluster we want to delete.
type clusterUninstaller struct {
	clusterID string
	infraID   string
	v3Client  *nutanixclientv3.Client
	logger    logrus.FieldLogger
}

type expectedCategory struct {
	key    string
	values []string
}

// New returns an Nutanix destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *installertypes.ClusterMetadata) (providers.Destroyer, error) {
	v3Client, err := nutanixtypes.CreateNutanixClient(context.TODO(),
		metadata.ClusterPlatformMetadata.Nutanix.PrismCentral,
		metadata.ClusterPlatformMetadata.Nutanix.Port,
		metadata.ClusterPlatformMetadata.Nutanix.Username,
		metadata.ClusterPlatformMetadata.Nutanix.Password,
	)
	if err != nil {
		return nil, err
	}

	return &clusterUninstaller{
		clusterID: metadata.ClusterID,
		infraID:   metadata.InfraID,
		v3Client:  v3Client,
		logger:    logger,
	}, nil
}

// Run is the entrypoint to start the uninstall process.
func (o *clusterUninstaller) Run() (*installertypes.ClusterQuota, error) {
	o.logger.Infof("Starting deletion of Nutanix infrastructure for Openshift cluster %s", o.infraID)
	err := wait.PollImmediateInfinite(time.Second*30, o.destroyCluster)
	if err != nil {
		return nil, errors.Wrap(err, "failed to destroy cluster")
	}

	return nil, nil
}

func (o *clusterUninstaller) destroyCluster() (bool, error) {
	cleanupFuncs := []struct {
		name    string
		execute func(*clusterUninstaller) error
	}{
		{name: "VMs", execute: cleanupVMs},
		{name: "Images", execute: cleanupImages},
		{name: "Categories", execute: cleanupCategories},
	}

	done := true
	for _, cleanupFunc := range cleanupFuncs {
		if err := cleanupFunc.execute(o); err != nil {
			o.logger.Debugf("%s: %q", cleanupFunc.name, err)
			done = false
		}
	}
	return done, nil
}

func cleanupVMs(o *clusterUninstaller) error {
	ec := createExpectedCategory(o.infraID)
	matchedVirtualMachineList := make([]*nutanixclientv3.VMIntentResource, 0)
	allVMs, err := o.v3Client.V3.ListAllVM(emptyFilter)
	if err != nil {
		return err
	}

	for _, v := range allVMs.Entities {
		if hasCategoryAssigned(v.Metadata, ec) {
			matchedVirtualMachineList = append(matchedVirtualMachineList, v)
		}
	}

	if len(matchedVirtualMachineList) == 0 {
		o.logger.Infof("no VMs found that require deletion for cluster %s", o.clusterID)
	} else {
		logToBeDeletedVMs(matchedVirtualMachineList, o.logger)
		err := deleteVMs(o.v3Client.V3, matchedVirtualMachineList, o.logger)
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanupImages(o *clusterUninstaller) error {
	ec := createExpectedCategory(o.infraID)
	bootISOImageName := nutanixtypes.BootISOImageName(o.infraID)
	rhcosImageName := nutanixtypes.RHCOSImageName(o.infraID)
	expectedImageNames := []string{
		bootISOImageName,
		rhcosImageName,
	}

	var found bool
	allImages, err := o.v3Client.V3.ListAllImage(emptyFilter)
	if err != nil {
		return err
	}

	for _, expectedImageName := range expectedImageNames {
		found = false
		for _, i := range allImages.Entities {
			imageName := *i.Spec.Name
			imageUUID := *i.Metadata.UUID
			if imageName == expectedImageName && hasCategoryAssigned(i.Metadata, ec) {
				found = true
				o.logger.Infof("deleting image %s with UUID %s", imageName, imageUUID)
				response, err := o.v3Client.V3.DeleteImage(imageUUID)
				if err != nil {
					o.logger.Errorf("failed to delete image %s: %q", imageUUID, err)
					break
				}

				if err = nutanixtypes.WaitForTask(o.v3Client.V3, response.Status.ExecutionContext.TaskUUID.(string)); err != nil {
					o.logger.Errorf("failed to confirm image deletion %s: %q", imageUUID, err)
				}

				break
			}
		}

		if !found {
			o.logger.Infof("no image with name %s was found", expectedImageName)
		}
	}
	return nil
}

func cleanupCategories(o *clusterUninstaller) error {
	ec := createExpectedCategory(o.infraID)
	_, err := o.v3Client.V3.GetCategoryKey(ec.key)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			//Already deleted
			return nil
		}
		return err
	}

	for _, val := range ec.values {
		o.logger.Infof("deleting category value : %s", val)
		err := o.v3Client.V3.DeleteCategoryValue(ec.key, val)
		if err != nil {
			o.logger.Errorf("failed to delete category value %s: %q", val, err)
		}
	}

	o.logger.Infof("deleting category key : %s", ec.key)
	return o.v3Client.V3.DeleteCategoryKey(ec.key)
}

func createExpectedCategory(infraID string) *expectedCategory {
	return &expectedCategory{
		key:    fmt.Sprintf(expectedCategoryKeyFormat, infraID),
		values: []string{expectedCategoryValueOwned, expectedCategoryValueShared},
	}
}

func logToBeDeletedVMs(tbd []*nutanixclientv3.VMIntentResource, l logrus.FieldLogger) {
	l.Info("virtual machines scheduled to be deleted: ")
	for _, v := range tbd {
		l.Infof("- %s", *v.Spec.Name)
	}
}

func deleteVMs(clientV3 nutanixclientv3.Service, vms []*nutanixclientv3.VMIntentResource, l logrus.FieldLogger) error {
	taskUUIDs := make([]string, 0)
	for _, vm := range vms {
		l.Infof("deleting VM %s with ID %s", *vm.Spec.Name, *vm.Metadata.UUID)
		response, err := clientV3.DeleteVM(*vm.Metadata.UUID)
		if err != nil {
			l.Errorf("failed to delete VM %s: %q", *vm.Metadata.UUID, err)
			continue
		}

		taskUUIDs = append(taskUUIDs, response.Status.ExecutionContext.TaskUUID.(string))
	}

	return nutanixtypes.WaitForTasks(clientV3, taskUUIDs)
}

func hasCategoryAssigned(metadata *nutanixclientv3.Metadata, ec *expectedCategory) bool {
	value, keyExists := metadata.Categories[ec.key]
	return keyExists && stringInSlice(value, ec.values)
}

func stringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}

	return false
}
