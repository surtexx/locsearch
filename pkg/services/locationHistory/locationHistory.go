package locationHistory

import (
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/surtexx/locsearch/pkg/db/dynamodbclient"
	"github.com/umahmood/haversine"
)

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

func GetDistanceTraveled(username, startDate, endDate string) float64 {
	svc := dynamodbclient.Connect()

	tableName := os.Getenv("LOCATION_HISTORY_TABLE")

	filt := expression.Name("Username").Equal(expression.Value(username))

	proj := expression.NamesList(expression.Name("Locations"), expression.Name("Timestamps"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		panic(err)
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	result, err := svc.Scan(input)
	if err != nil {
		panic(err)
	}

	locations := []Coordinate{}
	timestamps := []string{}
	for _, i := range result.Items {
		var location struct {
			Locations  []string `json:"Locations"`
			Timestamps []string `json:"Timestamps"`
		}
		err = dynamodbattribute.UnmarshalMap(i, &location)
		if err != nil {
			panic(err)
		}
		for _, location := range location.Locations {
			locationSplit := strings.Split(location, ",")
			latitude, _ := strconv.ParseFloat(locationSplit[0], 64)
			longitude, _ := strconv.ParseFloat(locationSplit[1], 64)
			locations = append(locations, Coordinate{Latitude: latitude, Longitude: longitude})
		}
		timestamps = location.Timestamps
	}

	// calculate distance traveled
	distanceTraveled := 0.0
	for i := 0; i < len(locations)-1; i++ {
		if timestamps[i] < startDate {
			continue
		}
		if timestamps[i+1] > endDate {
			break
		}
		location1 := haversine.Coord{Lat: locations[i].Latitude, Lon: locations[i].Longitude}
		location2 := haversine.Coord{Lat: locations[i+1].Latitude, Lon: locations[i+1].Longitude}
		km, _ := haversine.Distance(location1, location2)
		distanceTraveled += km
	}

	return distanceTraveled
}

func UpdateLocationHistory(username, location, currentTime string) {
	svc := dynamodbclient.Connect()

	tableName := os.Getenv("LOCATION_HISTORY_TABLE")

	getInput := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Username": {
				S: aws.String(username),
			},
		},
	}
	getResult, err := svc.GetItem(getInput)
	if err != nil {
		panic(err)
	}

	if getResult.Item == nil {
		addInput := &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]*dynamodb.AttributeValue{
				"Username": {
					S: aws.String(username),
				},
				"Locations": {
					L: []*dynamodb.AttributeValue{
						{
							S: aws.String(location),
						},
					},
				},
				"Timestamps": {
					L: []*dynamodb.AttributeValue{
						{
							S: aws.String(currentTime),
						},
					},
				},
			},
		}

		_, err := svc.PutItem(addInput)
		if err != nil {
			panic(err)
		}
	} else {
		updateInput := &dynamodb.UpdateItemInput{
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"Username": {
					S: aws.String(username),
				},
			},
			UpdateExpression: aws.String("SET #locations = list_append(#locations, :location), #timestamps = list_append(#timestamps, :timestamp)"),
			ExpressionAttributeNames: map[string]*string{
				"#locations":  aws.String("Locations"),
				"#timestamps": aws.String("Timestamps"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":location": {
					L: []*dynamodb.AttributeValue{
						{
							S: aws.String(location),
						},
					},
				},
				":timestamp": {
					L: []*dynamodb.AttributeValue{
						{
							S: aws.String(currentTime),
						},
					},
				},
			},
		}

		_, err := svc.UpdateItem(updateInput)
		if err != nil {
			panic(err)
		}
	}
}
