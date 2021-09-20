package finder

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kafka"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// FindScramSecrets returns the matching MSK Cluster's associated secrets
func FindScramSecrets(conn *kafka.Kafka, clusterArn string) ([]*string, error) {
	input := &kafka.ListScramSecretsInput{
		ClusterArn: aws.String(clusterArn),
	}

	var scramSecrets []*string
	err := conn.ListScramSecretsPages(input, func(page *kafka.ListScramSecretsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}
		scramSecrets = append(scramSecrets, page.SecretArnList...)
		return !lastPage
	})

	return scramSecrets, err
}
