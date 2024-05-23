//go:build integration
// +build integration

package rest

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/surtexx/locsearch/pkg/db/dynamodbclient"
	"github.com/surtexx/locsearch/pkg/services/location"
)

func TestSearchUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/searchUsers?coordinates=39.13354,27.14438&radius=10.0", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := Handler()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestUpdateLocation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, err := http.NewRequest("PUT", "http://127.0.0.1:8080/updateLocation?username=test&newLocation=39.13354,29.24438", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := Handler()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	username := "test"
	oldLocation := "39.13354,27.14438"
	newLocation := "39.13354,29.24438"

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

	_, err = svc.DeleteItem(input)
	if err != nil {
		panic(err)
	}

	location.UpdateLocation(username, newLocation)
	location.UpdateLocation(username, oldLocation)
}

func TestGetDistanceTraveled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/getDistanceTraveled?username=test&startDate=2024-01-01T00:00:00+00:00&endDate=2024-12-30T00:00:00+00:00", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := Handler()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
