package envssm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func extractKey(paramName string) string {
	// This assumes the path prefix follows the "/service-name/prod/" format
	parts := strings.Split(paramName, "/")
	return parts[len(parts)-1] // Return the last part as the key (e.g., "DB_HOST")
}

func createSsmClient() *ssm.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("[error] unable to load SDK config, %v", err)
	}
	return ssm.NewFromConfig(cfg)
}

func loadConfigFromAws(path string) []types.Parameter {
	fmt.Println("[info] loading config from AWS SSM...")

	// Create an SSM client
	ssmClient := createSsmClient()

	// Define parameters for retrieval
	withDecryption := true

	// Retrieve the parameters
	params, err := ssmClient.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
		Path:           &path,
		WithDecryption: &withDecryption,
	})

	if err != nil {
		log.Fatalf("[error] unable to load config from AWS SSM, %v", err)
	}

	fmt.Println("[success] Done loading config from AWS SSM")

	return params.Parameters
}

func deleteSsmParams(keys []string) {
	fmt.Println("[info] deleting SSM params...")

	ssmClient := createSsmClient()

	for _, key := range keys {
		_, err := ssmClient.DeleteParameter(context.TODO(), &ssm.DeleteParameterInput{Name: &key})
		if err != nil {
			log.Fatalf("[error] unable to delete SSM param %s, %v", key, err)
		}
	}

	fmt.Println("[success] Done deleting SSM params")
}

func addSsmParams(path string, envs map[string]string) {
	fmt.Println("[info] adding SSM params...")

	ssmClient := createSsmClient()
	overwrite := true

	for key, value := range envs {
		keyWithPath := path + key

		_, err := ssmClient.PutParameter(
			context.TODO(),
			&ssm.PutParameterInput{
				Name:      &keyWithPath,
				Value:     &value,
				Type:      types.ParameterTypeSecureString,
				Overwrite: &overwrite,
			},
		)
		if err != nil {
			log.Fatalf("[error] unable to add SSM param %s, %v", key, err)
		}
	}

	fmt.Println("[success] Done adding SSM params")
}
