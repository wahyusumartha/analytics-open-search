package infra

import (
	"github.com/opensearch-project/opensearch-go"
)

func NewOpenSearchClient(username, password, endpoint string) (*opensearch.Client, error) {
	client, err := opensearch.NewClient(
		opensearch.Config{
			Addresses: []string{endpoint},
			Username:  username,
			Password:  password,
		},
	)

	if err != nil {
		return nil, err
	}

	return client, err
}
