package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Global variable for DB connection
var db *gorm.DB

// Initialize the database
func initDB() {
	var err error
	// Define the Data Source Name (DSN)
	dsn := "root:password@tcp(127.0.0.1:3306)/go_app?charset=utf8&parseTime=True&loc=Local"

	// Open a connection to the MySQL database
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Migrate the schema (create the table if it doesn't exist)
	db.AutoMigrate(&User{})
}

// var users = []User{
// 	{ID: 1, Name: "John", Email: "john@example.com"},
// 	{ID: 2, Name: "Jane", Email: "jane@example.com"},
// }

func main() {
	// Initialize the database
	initDB()

	r := gin.Default()

	// Define a GET route
	r.GET("/users", func(c *gin.Context) {
		var users []User
		// Find all users from the database
		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
			return
		}

		// Return the users
		c.JSON(http.StatusOK, users)
	})

	// Define a GET route with id parameter
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		// for _, user := range users {
		// 	if userID, err := strconv.Atoi(id); err == nil && user.ID == uint(userID) {
		// 		c.JSON(http.StatusOK, user)
		// 		return
		// 	}
		// }
		// c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		var user User
		// Find the user with the given ID
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	r.POST("/users", func(c *gin.Context) {
		// var newUser User
		// if err := c.ShouldBindJSON(&newUser); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }
		// newUser.ID = uint(len(users) + 1)
		// users = append(users, newUser)
		// c.JSON(http.StatusCreated, newUser)

		var user User

		// Bind the incoming JSON to the User struct
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Insert the user data into the MySQL database
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// Return the created user data
		c.JSON(http.StatusCreated, user)
	})
	r.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		// for i, user := range users {
		// 	if userID, err := strconv.Atoi(id); err == nil && user.ID == uint(userID) {
		// 		users = slices.Delete(users, i, i+1)
		// 		c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
		// 		return
		// 	}
		// }
		// c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		var user User
		// Find the user with the given ID
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		// Delete the user from the database
		if err := db.Delete(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Error deleting user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
	})

	r.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		// for i, user := range users {
		// 	if userID, err := strconv.Atoi(id); err == nil && user.ID == uint(userID) {
		// 		var updatedUser User
		// 		if err := c.ShouldBindJSON(&updatedUser); err != nil {
		// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 			return
		// 		}
		// 		updatedUser.ID = user.ID
		// 		users[i] = updatedUser
		// 		c.JSON(http.StatusOK, updatedUser)
		// 		return
		// 	}
		// }
		// c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		var user User
		// Find the user with the given ID
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		// Bind the JSON data to the user struct
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if userID, err := strconv.Atoi(id); err == nil {
			user.ID = uint(userID)
		}
		// Update the user in the database
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	// Start the server
	r.Run(":8080") // Run on port 8080
}
