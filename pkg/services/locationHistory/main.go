package locationHistory

import (
	"context"

	"github.com/umahmood/haversine"
)

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type Locations struct {
	Username  string  `json:"username"`
	Coordinates []Coordinate `json:"coordinates"`

func GetDistanceTraveled(username, startDate, endDate string) float64 {
	mongoClient := connectToMongo()
	defer mongoClient.Disconnect(context.Background())

	collection := mongoClient.Database("locsearch").Collection("locationHistory")

	// find user
	filter := bson.D{{"username", username}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		panic(err)
	}

	for cursor.Next(context.Background()) {
		var locations Locations
		err := cursor.Decode(&locations)
		if err != nil {
			panic(err)
		}
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
