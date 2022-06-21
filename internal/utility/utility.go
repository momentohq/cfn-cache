package utility

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

// GetSecret fetches a secret from secrets manager.
func GetSecret(svc secretsmanageriface.SecretsManagerAPI, secretName string) (string, error) {
	secretValue, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get secret %s err=%+v", secretName, err)
	}

	return *secretValue.SecretString, nil
}
