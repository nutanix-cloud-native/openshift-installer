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
	category := createExpectedCategory(o.infraID)

	cleanupFuncs := []struct {
		name    string
		execute func(nutanixclientv3.Service, *expectedCategory, *clusterUninstaller) error
	}{
		{name: "VMs", execute: cleanupVMs},
		{name: "Images", execute: cleanupImages},
		{name: "Categories", execute: cleanupCategories},
	}
	done := true
	for _, cleanupFunc := range cleanupFuncs {
		if done {
			err := cleanupFunc.execute(o.v3Client.V3, category, o)
			if err != nil {
				o.logger.Debugf("%s: %v", cleanupFunc.name, err)
				done = false
			}
		}
	}
	return done, nil
}

func cleanupVMs(clientV3 nutanixclientv3.Service, ec *expectedCategory, o *clusterUninstaller) error {
	matchedVirtualMachineList := make([]*nutanixclientv3.VMIntentResource, 0)
	allVMsRaw, err := clientV3.ListAllVM(emptyFilter)
	if err != nil {
		return err
	}
	allVMs := allVMsRaw.Entities
	for _, v := range allVMs {
		vmName := *v.Spec.Name
		if strings.HasPrefix(vmName, o.infraID) {
			if hasCategoryAssigned(v.Metadata, ec) {
				matchedVirtualMachineList = append(matchedVirtualMachineList, v)
			}
		}
	}

	if len(matchedVirtualMachineList) == 0 {
		o.logger.Infof("No VMs found that require deletion for cluster %s", o.clusterID)
	} else {
		logToBeDeletedVMs(matchedVirtualMachineList, o.logger)
		err := deleteVMs(clientV3, matchedVirtualMachineList, o.logger)
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanupImages(clientV3 nutanixclientv3.Service, ec *expectedCategory, o *clusterUninstaller) error {
	bootISOImageName := nutanixtypes.BootISOImageName(o.infraID)
	rhcosImageName := nutanixtypes.RHCOSImageName(o.infraID)
	expectedImageNames := []string{
		bootISOImageName,
		rhcosImageName,
	}

	var found bool
	allImagesRaw, err := clientV3.ListAllImage(emptyFilter)
	if err != nil {
		return err
	}
	allImages := allImagesRaw.Entities
	for _, expectedImageName := range expectedImageNames {
		found = false
		for _, i := range allImages {
			imageName := *i.Spec.Name
			imageUUID := *i.Metadata.UUID
			if imageName == expectedImageName {
				if hasCategoryAssigned(i.Metadata, ec) {
					found = true
					o.logger.Infof("Deleting image %s with UUID %s", imageName, imageUUID)
					response, err := clientV3.DeleteImage(imageUUID)
					if err != nil {
						return err
					}
					err = nutanixtypes.WaitForTask(clientV3, response.Status.ExecutionContext.TaskUUID.(string))
					if err != nil {
						return err
					}
					break
				}
			}
		}
		if !found {
			o.logger.Infof("No image with name %s was found", expectedImageName)
		}
	}
	return nil
}

func cleanupCategories(clientV3 nutanixclientv3.Service, ec *expectedCategory, o *clusterUninstaller) error {
	_, err := clientV3.GetCategoryKey(ec.key)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			//Already deleted
			return nil
		}
		return err
	}

	for _, val := range ec.values {
		o.logger.Infof("Deleting category value : %s", val)
		err := clientV3.DeleteCategoryValue(ec.key, val)
		if err != nil {
			return err
		}
	}

	o.logger.Infof("Deleting category key : %s", ec.key)
	return clientV3.DeleteCategoryKey(ec.key)
}

func createExpectedCategory(infraID string) *expectedCategory {
	return &expectedCategory{
		key:    fmt.Sprintf(expectedCategoryKeyFormat, infraID),
		values: []string{expectedCategoryValueOwned, expectedCategoryValueShared},
	}
}

func logToBeDeletedVMs(tbd []*nutanixclientv3.VMIntentResource, l logrus.FieldLogger) {
	l.Info("VMs scheduled to be deleted: ")
	for _, v := range tbd {
		l.Infof("- %s", *v.Spec.Name)
	}
}

func deleteVMs(clientV3 nutanixclientv3.Service, tbd []*nutanixclientv3.VMIntentResource, l logrus.FieldLogger) error {
	taskUUIDs := make([]string, 0)
	for _, v := range tbd {
		l.Infof("Deleting vm %s with ID %s", *v.Spec.Name, *v.Metadata.UUID)
		response, err := clientV3.DeleteVM(*v.Metadata.UUID)
		if err != nil {
			return err
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
