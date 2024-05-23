// This file contains the functions to update a user's location and search for users within a certain radius.
package location

import (
	"os"
	"strconv"
	"strings"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/surtexx/locsearch/pkg/db/dynamodbclient"
	"github.com/surtexx/locsearch/pkg/services/locationHistory"
	"github.com/umahmood/haversine"
)

type User struct {
	Username  string
	Location  string
	Timestamp string
}

func UpdateLocation(username, newLocation string) {
	svc := dynamodbclient.Connect()

	tableName := os.Getenv("LOCATIONS_TABLE")

	currentTime := time.Now().Format(time.RFC3339)

	// Check if username exists
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {
				S: aws.String(username),
			},
		},
	}

	result, err := svc.GetItem(getItemInput)
	if err != nil {
		panic(err)
	}
	if result.Item == nil {
		item := &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]*dynamodb.AttributeValue{
				"Username": {
					S: aws.String(username),
				},
				"Location": {
					S: aws.String(newLocation),
				},
				"Timestamp": {
					S: aws.String(currentTime),
				},
			},
		}

		_, err = svc.PutItem(item)
		if err != nil {
			panic(err)
		}
	} else {
		updateItemInput := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":l": {
					S: aws.String(newLocation),
				},
				":t": {
					S: aws.String(currentTime),
				},
			},
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"Username": {
					S: aws.String(username),
				},
			},
			ReturnValues:     aws.String("UPDATED_NEW"),
			UpdateExpression: aws.String("set #loc = :l, #ts = :t"),
			ExpressionAttributeNames: map[string]*string{
				"#loc": aws.String("Location"),
				"#ts":  aws.String("Timestamp"),
			},
		}

		_, err = svc.UpdateItem(updateItemInput)
		if err != nil {
			panic(err)
		}
	}
	locationHistory.UpdateLocationHistory(username, newLocation, currentTime)
}

func SearchUsers(coordinates string, radius float64) []string {
	svc := dynamodbclient.Connect()

	tableName := os.Getenv("LOCATIONS_TABLE")

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := svc.Scan(input)
	if err != nil {
		panic(err)
	}

	var users []User
	for _, item := range result.Items {
		user := User{}

		err = dynamodbattribute.UnmarshalMap(item, &user)
		if err != nil {
			panic(err)
		}

		users = append(users, user)
	}

	coordinates = strings.TrimSpace(coordinates)
	coordinatesSplit := strings.Split(coordinates, ",")
	coordinatesSplit[0] = strings.TrimSpace(coordinatesSplit[0])
	coordinatesSplit[1] = strings.TrimSpace(coordinatesSplit[1])

	latitude, err := strconv.ParseFloat(coordinatesSplit[0], 64)
	if err != nil {
		panic(err)
	}

	longitude, err := strconv.ParseFloat(coordinatesSplit[1], 64)
	if err != nil {
		panic(err)
	}

	coords := haversine.Coord{Lat: latitude, Lon: longitude}

	var usernames []string
	for _, user := range users {
		userCoordinates := strings.Split(user.Location, ",")
		userCoordinates[0] = strings.TrimSpace(userCoordinates[0])
		userCoordinates[1] = strings.TrimSpace(userCoordinates[1])

		userLatitude, err := strconv.ParseFloat(userCoordinates[0], 64)
		if err != nil {
			panic(err)
		}

		userLongitude, err := strconv.ParseFloat(userCoordinates[1], 64)
		if err != nil {
			panic(err)
		}

		userCoords := haversine.Coord{Lat: userLatitude, Lon: userLongitude}

		km, _ := haversine.Distance(coords, userCoords)

		if km <= radius {
			usernames = append(usernames, user.Username)
		}
	}
	return usernames
}
