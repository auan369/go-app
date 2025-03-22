// Description: This file contains the test cases for the CRUD operations of the user API.
package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Helper function to set up a test router
func setupRouter() *gin.Engine {
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
		var user User
		// Find the user with the given ID
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	r.POST("/users", func(c *gin.Context) {
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

	return r
}

// Initialize the database for testing
func initTestDB() {
	var err error
	// Define the Data Source Name (DSN)
	dsn := "root:password@tcp(127.0.0.1:3306)/go_app_test?charset=utf8&parseTime=True&loc=Local"

	// Open a connection to the MySQL database
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the test database: %v", err)
	}

	// Migrate the schema (create the table if it doesn't exist)
	db.AutoMigrate(&User{})
}

// Test GET /users
func TestGetUsers(t *testing.T) {
	initTestDB()
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John") // Check if "John" exists in response
}

// Test GET /users/:id (valid ID)
func TestGetUserByID_Valid(t *testing.T) {
	initTestDB()
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/users/2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John") // Check if "John" is returned
}

// Test GET /users/:id (invalid ID)
func TestGetUserByID_Invalid(t *testing.T) {
	initTestDB()
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/users/99", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}

// Test POST /users
func TestCreateUser(t *testing.T) {
	initTestDB()
	router := setupRouter()

	newUserJSON := `{"name": "John Doe", "email": "john@example.com", "password": "password"}`
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte(newUserJSON)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "John Doe")
	//Delete the user based on email
	db.Where("email = ?", "john@example.com").Delete(&User{})
}

// Test PUT /users/:id (valid ID)
func TestUpdateUser_Valid(t *testing.T) {
	initTestDB()
	router := setupRouter()

	updatedUserJSON := `{"name": "John Ray", "email": "johnray@example.com", "password": "newpassword"}`
	req, _ := http.NewRequest("PUT", "/users/4", bytes.NewBuffer([]byte(updatedUserJSON)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John Ray")
}

// Test PUT /users/:id (invalid ID)
func TestUpdateUser_Invalid(t *testing.T) {
	initTestDB()
	router := setupRouter()

	updatedUserJSON := `{"name": "John Ray", "email": "johnray@example.com", "password": "newpassword"}`
	req, _ := http.NewRequest("PUT", "/users/99", bytes.NewBuffer([]byte(updatedUserJSON)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}

// Test DELETE /users/:id (valid ID)
func TestDeleteUser_Valid(t *testing.T) {
	initTestDB()
	router := setupRouter()
	//create a user to delete with id 1
	user := User{
		ID:       1, // Set the user ID explicitly
		Name:     "John Doe",
		Email:    "sd@email.com",
		Password: "password", // Ensure you're hashing the password in your real implementation
	}
	db.Create(&user)
	//delete the user
	// newUserJSON := `{"id": "1", "name": "John Doe", "email": "sd@email.com", "password": "password"}`
	// req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte(newUserJSON)))
	// req.Header.Set("Content-Type", "application/json")
	// w := httptest.NewRecorder()
	// router.ServeHTTP(w, req)

	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User deleted")
}

// Test DELETE /users/:id (invalid ID)
func TestDeleteUser_Invalid(t *testing.T) {
	initTestDB()
	router := setupRouter()

	req, _ := http.NewRequest("DELETE", "/users/99", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}
