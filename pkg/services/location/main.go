package location

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/umahmood/haversine"
)

type User struct {
	Username string `json:"username"`
	Location string `json:"location"`
}

func UpdateLocation(username, newLocation string) {
	mongoClient := connectToMongo()
	defer mongoClient.Disconnect(context.Background())

	collection := mongoClient.Database("locsearch").Collection("locations")

	filter := bson.D{{"username", username}}
	update := bson.D{
		{"$set", bson.D{
			{"location", newLocation},
		}}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		panic(err)
	}

	fmt.Println("Location updated")
}

func SearchUsers(coordinates string, radius float64) []string {
	mongoClient := connectToMongo()
	defer mongoClient.Disconnect(context.Background())

	collection := mongoClient.Database("locsearch").Collection("locations")

	users := []User{}

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		panic(err)
	}

	for cursor.Next(context.Background()) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			panic(err)
		}

		users = append(users, user)
	}

	usernames := []string{}

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
