package googlecloud

import (
	"context"
	"fmt"
	"hash/crc32"
	"regexp"
	"strings"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/gumi-tsd/secret-env-manager/internal/model"
	"google.golang.org/api/iterator"
)

func AccessSecretVersion(name string) (*string, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %w", err)
	}

	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		return nil, fmt.Errorf("Data corruption detected.")
	}

	data := string(result.Payload.Data)

	// env format
	data = strings.ReplaceAll(data, "\n", "")
	data = strings.ReplaceAll(data, "\r", "")

	return &data, nil
}

func ListSecrets(project string) (*model.Secrets, error) {
	parent := fmt.Sprintf("projects/%s", project)

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	req := &secretmanagerpb.ListSecretsRequest{
		Parent: parent,
	}

	it := client.ListSecrets(ctx, req)

	gcpSecrets := model.Secrets{}
	re := regexp.MustCompile(`[^/]+$`)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to list secrets: %w", err)
		}
		result := re.FindString(resp.Name)

		unixTime := time.Unix(resp.CreateTime.GetSeconds(), int64(resp.CreateTime.GetNanos()))
		rfc3339Time := unixTime.Format(time.RFC3339)

		gcpSecret := model.Secret{
			Name:      result,
			CreatedAt: rfc3339Time,
			Version:   "latest",
		}

		gcpSecrets.Secrets = append(gcpSecrets.Secrets, gcpSecret)
	}

	return &gcpSecrets, nil
}
