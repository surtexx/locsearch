package locationHistory

import (
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

	tableName := "locationHistory"

	filt := expression.Name("Username").Equal(expression.Value(username)).And(expression.Name("Timestamp").Between(expression.Value(startDate), expression.Value(endDate)))

	proj := expression.NamesList(expression.Name("Location"))

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
	for _, i := range result.Items {
		var location string
		err = dynamodbattribute.UnmarshalMap(i, &location)
		if err != nil {
			panic(err)
		}

		locationSplit := strings.Split(location, ",")
		locationSplit[0] = strings.Trim(locationSplit[0], " ")
		locationSplit[1] = strings.Trim(locationSplit[1], " ")

		latitude, err := strconv.ParseFloat(locationSplit[0], 64)
		if err != nil {
			panic(err)
		}

		longitude, err := strconv.ParseFloat(locationSplit[1], 64)
		if err != nil {
			panic(err)
		}

		locations = append(locations, Coordinate{Latitude: latitude, Longitude: longitude})
	}

	// calculate distance traveled
	distanceTraveled := 0.0
	for i := 0; i < len(locations)-1; i++ {
		location1 := haversine.Coord{Lat: locations[i].Latitude, Lon: locations[i].Longitude}
		location2 := haversine.Coord{Lat: locations[i+1].Latitude, Lon: locations[i+1].Longitude}
		km, _ := haversine.Distance(location1, location2)
		distanceTraveled += km
	}

	return distanceTraveled
}

// func UpdateLocationHistory(username, location string) {
// 	svc := dynamodbclient.Connect()

// 	tableName := "locationHistory"

// 	_, err := svc.UpdateItem(input)
// 	if err != nil {
// 		panic(err)
// 	}
// }
