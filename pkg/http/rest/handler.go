// This file contains the REST API handlers for the location and locationHistory services.
package rest

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/surtexx/locsearch/pkg/services/location"
	"github.com/surtexx/locsearch/pkg/services/locationHistory"
	"github.com/surtexx/locsearch/pkg/validating"
)

// Handler returns a new gin engine with the REST API handlers.
func Handler() *gin.Engine {
	r := gin.Default()
	r.PUT("/updateLocation", updateLocation)
	r.GET("/searchUsers", searchUsers)
	r.GET("/getDistanceTraveled", getDistanceTraveled)
	return r
}

func updateLocation(c *gin.Context) {
	username := c.Query("username")
	newLocation := c.Query("newLocation")

	if !validating.ValidateUsername(username) {
		c.JSON(400, gin.H{"error": "Invalid username"})
		return
	}
	if !validating.ValidateCoordinates(newLocation) {
		c.JSON(400, gin.H{"error": "Invalid coordinates"})
		return
	}
	location.UpdateLocation(username, newLocation)
	c.JSON(200, gin.H{"message": "Location updated"})
}

func searchUsers(c *gin.Context) {
	coordinates := c.Query("coordinates")
	if !validating.ValidateCoordinates(coordinates) {
		c.JSON(400, gin.H{"error": "Invalid coordinates"})
		return
	}

	radiusStr := c.Query("radius")
	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid radius"})
		return
	}

	users := location.SearchUsers(coordinates, radius)
	c.JSON(200, gin.H{"users": users})
}

func getDistanceTraveled(c *gin.Context) {
	username := c.Query("username")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	if !validating.ValidateUsername(username) {
		c.JSON(400, gin.H{"error": "Invalid username"})
		return
	}
	if !validating.ValiDate(startDate) {
		c.JSON(400, gin.H{"error": "Invalid date. Use ISO 8601 format (YYY-MM-DDYTHH:MM:SS+HH:MM)"})
		return
	}
	if !validating.ValiDate(endDate) {
		c.JSON(400, gin.H{"error": "Invalid date. Use ISO 8601 format (YYY-MM-DDYTHH:MM:SS+HH:MM)"})
		return
	}

	distance := locationHistory.GetDistanceTraveled(username, startDate, endDate)
	c.JSON(200, gin.H{"distance": distance})
}
