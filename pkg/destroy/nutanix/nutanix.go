package nutanix

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/destroy/providers"
	installertypes "github.com/openshift/installer/pkg/types"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	nutanixClientV3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

const (
	Description = "Created By OpenShift Installer"
	EmptyFilter = ""
)

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	ClusterID string
	InfraID   string
	V3Client  *nutanixClientV3.Client
	Logger    logrus.FieldLogger
}

type ExpectedCategory struct {
	Key   string
	Value string
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

	return &ClusterUninstaller{
		ClusterID: metadata.ClusterID,
		InfraID:   metadata.InfraID,
		V3Client:  v3Client,
		Logger:    logger,
	}, nil
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() (*installertypes.ClusterQuota, error) {
	ec := createExpectedCategory(o.InfraID)
	bootISOImageName := nutanixtypes.BootISOImageName(o.InfraID)
	rhcosImageName := nutanixtypes.RHCOSImageName(o.InfraID)
	imagesToDelete := []string{
		bootISOImageName,
		rhcosImageName,
	}
	o.Logger.Infof("Starting deletion of Nutanix infrastructure for Openshift cluster %s", o.InfraID)
	clientV3 := o.V3Client.V3
	//Delete VMs
	err := cleanupVMs(clientV3, ec, o)
	if err != nil {
		return nil, err
	}

	//Delete image
	err = cleanupImages(clientV3, ec, o, imagesToDelete)
	if err != nil {
		return nil, err
	}
	//Delete category
	err = cleanupCategory(clientV3, ec, o)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func cleanupVMs(clientV3 nutanixClientV3.Service, ec *ExpectedCategory, o *ClusterUninstaller) error {
	matchedVirtualMachineList := make([]*nutanixClientV3.VMIntentResource, 0)
	allVMsRaw, err := clientV3.ListAllVM(EmptyFilter)
	if err != nil {
		return err
	}
	allVMs := allVMsRaw.Entities
	for _, v := range allVMs {
		vmName := *v.Spec.Name
		if strings.HasPrefix(vmName, o.InfraID) {
			if hasCategoryAssigned(v.Metadata, ec, o.Logger) {
				if v.Spec.Description != nil {
					if *v.Spec.Description == Description {
						matchedVirtualMachineList = append(matchedVirtualMachineList, v)
					}
				}
			}
		}
	}

	amountOfTBD := len(matchedVirtualMachineList)
	if amountOfTBD == 0 {
		o.Logger.Infof("No VMs found that require deletion for cluster %s", o.ClusterID)
	} else {
		logToBeDeletedVMs(matchedVirtualMachineList, o.Logger)
		err := deleteVMs(clientV3, matchedVirtualMachineList, o.Logger)
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanupImages(clientV3 nutanixClientV3.Service, ec *ExpectedCategory, o *ClusterUninstaller, expectedImageNames []string) error {
	var found bool
	allImagesRaw, err := clientV3.ListAllImage(EmptyFilter)
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
				if hasCategoryAssigned(i.Metadata, ec, o.Logger) {
					if i.Spec.Description != nil {
						if *i.Spec.Description == Description {
							found = true
							o.Logger.Infof("Deleting image %s with UUID %s", imageName, imageUUID)
							response, err := clientV3.DeleteImage(imageUUID)
							if err != nil {
								return err
							}
							nutanixtypes.WaitForTask(clientV3, response.Status.ExecutionContext.TaskUUID.(string))
							break
						}
					}
				}
			}
		}
		if !found {
			o.Logger.Infof("No image with name %s was found", expectedImageName)
		}
	}
	return nil
}

func cleanupCategory(clientV3 nutanixClientV3.Service, ec *ExpectedCategory, o *ClusterUninstaller) error {
	_, err := clientV3.GetCategoryKey(ec.Key)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "does not exist.") {
			//Already deleted
			return nil
		}
		return err
	}

	allCategoryValuesRaw, err := clientV3.ListAllCategoryValues(ec.Key, EmptyFilter)
	if err != nil {
		return err
	}

	l := o.Logger
	allCategoryValues := allCategoryValuesRaw.Entities
	if len(allCategoryValues) == 1 {
		l.Infof("Deleting category value : %s", ec.Value)
		err := clientV3.DeleteCategoryValue(ec.Key, ec.Value)
		if err != nil {
			return err
		}
	}

	if len(allCategoryValues) > 1 {
		return errors.Errorf("multiple category values found for category %s", ec.Key)
	}

	l.Infof("Deleting category key : %s", ec.Key)
	return clientV3.DeleteCategoryKey(ec.Key)
}

func createExpectedCategory(infraID string) *ExpectedCategory {
	return &ExpectedCategory{
		Key:   fmt.Sprintf("openshift-%s", infraID),
		Value: "openshift-ipi-installations",
	}
}

func logToBeDeletedVMs(tbd []*nutanixClientV3.VMIntentResource, l logrus.FieldLogger) {
	l.Info("VMs scheduled to be deleted: ")
	for _, v := range tbd {
		l.Infof("- %s", *v.Spec.Name)
	}
}

func deleteVMs(clientV3 nutanixClientV3.Service, tbd []*nutanixClientV3.VMIntentResource, l logrus.FieldLogger) error {
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

func hasCategoryAssigned(metadata *nutanixClientV3.Metadata, ec *ExpectedCategory, l logrus.FieldLogger) bool {
	if metadata.Categories != nil {
		if val, ok := metadata.Categories[ec.Key]; ok {
			if val == ec.Value {
				return true
			}
		}
	}
	return false
}
