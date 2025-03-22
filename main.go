package main

import (
	"net/http"
	"strconv"

	"slices"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{
	{ID: 1, Name: "John", Email: "john@example.com"},
	{ID: 2, Name: "Jane", Email: "jane@example.com"},
}

func main() {
	r := gin.Default()

	// Define a GET route
	r.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, users)
	})

	// Define a GET route with id parameter
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, user := range users {
			if userID, err := strconv.Atoi(id); err == nil && user.ID == uint(userID) {
				c.JSON(http.StatusOK, user)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
	})

	r.POST("/users", func(c *gin.Context) {
		var newUser User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newUser.ID = uint(len(users) + 1)
		users = append(users, newUser)
		c.JSON(http.StatusCreated, newUser)
	})

	r.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		for i, user := range users {
			if userID, err := strconv.Atoi(id); err == nil && user.ID == uint(userID) {
				users = slices.Delete(users, i, i+1)
				c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
	})

	r.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		for i, user := range users {
			if userID, err := strconv.Atoi(id); err == nil && user.ID == uint(userID) {
				var updatedUser User
				if err := c.ShouldBindJSON(&updatedUser); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				updatedUser.ID = user.ID
				users[i] = updatedUser
				c.JSON(http.StatusOK, updatedUser)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
	})

	// Start the server
	r.Run(":8080") // Run on port 8080
}
