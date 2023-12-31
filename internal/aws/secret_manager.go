package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gumi-tsd/secret-env-manager/internal/model"
)

func ListSecrets(profile string) (*model.Secrets, error) {
	region := "ap-northeast-1"

	config, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.ListSecretsInput{}

	result, err := svc.ListSecrets(context.TODO(), input)
	if err != nil {
		log.Fatal(err.Error())
	}

	secrets := model.Secrets{}
	for _, secret := range result.SecretList {
		secrets.Secrets = append(secrets.Secrets, model.Secret{
			Name:      *secret.Name,
			CreatedAt: secret.CreatedDate.String(),
			Version:   "AWSCURRENT",
		})
	}

	return &secrets, nil
}

func AccessSecret(env model.Env) (string, error) {
	region := "ap-northeast-1"

	config, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(env.Account),
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(env.SecretName),
		VersionStage: aws.String(env.Version),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal(err.Error())
	}

	return *result.SecretString, nil
}
