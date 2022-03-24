package nutanix

import (
	"context"
	"fmt"
	"time"

	nutanixClient "github.com/terraform-providers/terraform-provider-nutanix/client"
	nutanixClientV3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"k8s.io/klog"
)

func CreateNutanixClient(ctx context.Context, prismCentral, port, username, password string) (*nutanixClientV3.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	cred := nutanixClient.Credentials{
		URL:      fmt.Sprintf("%s:%s", prismCentral, port),
		Username: username,
		Password: password,
		Port:     port,
		Endpoint: prismCentral,
		Insecure: true,
	}

	cli, err := nutanixClientV3.NewV3Client(cred)
	if err != nil {
		klog.Errorf("Failed to create the nutanix client. error: %v", err)
		return nil, err
	}

	return cli, nil
}
