package rest

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/surtexx/locsearch/pkg/services/location"
	"github.com/surtexx/locsearch/pkg/validating"
)

func Handler() *gin.Engine {
	r := gin.Default()
	r.PUT("/updateLocation", updateLocation)
	r.GET("/searchUsers", searchUsers)
	r.GET("/getDistanceTraveled", getDistanceTraveled)
	return r
}

func updateLocation(c *gin.Context) {
	if !validating.ValidateUsername(c.Query("username")) {
		c.JSON(400, gin.H{"error": "Invalid username"})
	}
	if !validating.ValidateCoordinates(c.Query("newLocation")) {
		c.JSON(400, gin.H{"error": "Invalid coordinates"})
	}
	if !validating.ValiDate(c.Query("timestamp")) {
		c.JSON(400, gin.H{"error": "Invalid date. Use ISO 8601 format (YYY-MM-DDYTHH:MM:SS+HH:MM)"})
	}
	username := c.Query("username")
	newLocation := c.Query("newLocation")

	location.UpdateLocation(username, newLocation)
}

func searchUsers(c *gin.Context) {
	if !validating.ValidateCoordinates(c.Query("coordinates")) {
		c.JSON(400, gin.H{"error": "Invalid coordinates"})
	}
	coordinates := c.Query("coordinates")
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
	if !validating.ValidateUsername(c.Query("username")) {
		c.JSON(400, gin.H{"error": "Invalid username"})
	}
	if !validating.ValiDate(c.Query("startDate")) {
		c.JSON(400, gin.H{"error": "Invalid date. Use ISO 8601 format (YYY-MM-DDYTHH:MM:SS+HH:MM)"})
	}
	if !validating.ValiDate(c.Query("endDate")) {
		c.JSON(400, gin.H{"error": "Invalid date. Use ISO 8601 format (YYY-MM-DDYTHH:MM:SS+HH:MM)"})
	}
	username := c.Query("username")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	distance := locationHistory.GetDistanceTraveled(username, startDate, endDate)
	c.JSON(200, gin.H{"distance": distance})
}
