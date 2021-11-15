package nutanix

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kdomanski/iso9660"
	"github.com/pkg/errors"
	nutanixClient "github.com/terraform-providers/terraform-provider-nutanix/client"
	nutanixClientV3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
	"k8s.io/klog"
)

const (
	DiskLabel        = "config-2"
	ISOFile          = "bootstrap-ign.iso"
	MetadataFilePath = "openstack/latest/meta_data.json"
	UserDataFilePath = "openstack/latest/user_data"
	SleepTime        = 10 * time.Second
	Timeout          = 5 * time.Minute
)

type MetadataCloudInit struct {
	UUID string `json: "uuid"`
}

func CreateNutanixClient(ctx context.Context, prismCentral, port, username, password string, insecure bool) (*nutanixClientV3.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	cred := nutanixClient.Credentials{
		URL:      fmt.Sprintf("%s:%s", prismCentral, port),
		Username: username,
		Password: password,
		Port:     port,
		Endpoint: prismCentral,
		Insecure: insecure,
	}

	cli, err := nutanixClientV3.NewV3Client(cred)
	if err != nil {
		klog.Errorf("Failed to create the nutanix client. error: %v", err)
		return nil, err
	}

	return cli, nil

}

func GenerateBootISOImageName(infraID string) string {
	return fmt.Sprintf("%s-%s", infraID, ISOFile)

}

func CreateBootstrapISO(infraID, userData string) (string, error) {
	id := uuid.New()
	metaObj := &MetadataCloudInit{
		UUID: id.String(),
	}
	fullISOFile := GenerateBootISOImageName(infraID)
	metadata, err := json.Marshal(metaObj)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed marshal metadata struct to json"))
	}
	writer, err := iso9660.NewWriter()
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to create writer: %s", err))
	}
	defer writer.Cleanup()

	userDataReader := strings.NewReader(userData)
	err = writer.AddFile(userDataReader, UserDataFilePath)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to add file: %s", err))
	}

	metadataReader := strings.NewReader(string(metadata))
	err = writer.AddFile(metadataReader, MetadataFilePath)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to add file: %s", err))
	}

	outputFile, err := os.OpenFile(fullISOFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to create file: %s", err))
	}

	err = writer.WriteTo(outputFile, DiskLabel)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to write ISO image: %s", err))
	}

	err = outputFile.Close()
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to close output file: %s", err))
	}
	return fullISOFile, nil
}

func WaitForTasks(clientV3 nutanixClientV3.Service, taskUUIDs []string) error {
	for _, t := range taskUUIDs {
		err := WaitForTask(clientV3, t)
		if err != nil {
			return err
		}
	}
	return nil
}

func WaitForTask(clientV3 nutanixClientV3.Service, taskUUID string) error {
	finished := false
	var err error
	for start := time.Now(); time.Since(start) < Timeout; {
		finished, err = isTaskFinished(clientV3, taskUUID)
		if err != nil {
			return err
		}
		if finished {
			break
		}
		time.Sleep(SleepTime)
	}
	if !finished {
		return errors.Errorf("timeout while waiting for task UUID: %s", taskUUID)
	}

	return nil
}

func isTaskFinished(clientV3 nutanixClientV3.Service, taskUUID string) (bool, error) {
	isFinished := map[string]bool{
		"QUEUED":    false,
		"RUNNING":   false,
		"SUCCEEDED": true,
	}
	status, err := getTaskStatus(clientV3, taskUUID)
	if err != nil {
		return false, err
	}
	if val, ok := isFinished[status]; ok {
		return val, nil
	}
	return false, errors.Errorf("Retrieved unexpected task status: %s", status)

}

func getTaskStatus(clientV3 nutanixClientV3.Service, taskUUID string) (string, error) {
	v, err := clientV3.GetTask(taskUUID)

	if err != nil {
		return "", err
	}

	if *v.Status == "INVALID_UUID" || *v.Status == "FAILED" {
		return *v.Status, errors.Errorf("error_detail: %s, progress_message: %s", utils.StringValue(v.ErrorDetail), utils.StringValue(v.ProgressMessage))
	}
	return *v.Status, nil
}
