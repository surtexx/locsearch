//go:build integration
// +build integration

package location

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/surtexx/locsearch/pkg/db/dynamodbclient"
)

func TestUpdateLocation(t *testing.T) {
	username := "test"
	oldLocation := "39.13354,27.14438"
	newLocation := "39.13354,29.24438"

	UpdateLocation(username, newLocation)

	users := SearchUsers(newLocation, 0.0)
	if len(users) == 0 {
		t.Error("Expected location to be updated")
	}

	// cleanup
	svc := dynamodbclient.Connect()
	tableName := os.Getenv("LOCATION_HISTORY_TABLE")

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {
				S: aws.String(username),
			},
		},
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		panic(err)
	}

	UpdateLocation(username, newLocation)
	UpdateLocation(username, oldLocation)
}
