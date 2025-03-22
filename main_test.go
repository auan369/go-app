package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Helper function to set up a test router
func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, users)
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		// for _, user := range users {
		// 	if id == string(rune(user.ID)) {
		// 		c.JSON(http.StatusOK, user)
		// 		return
		// 	}
		// }
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

	return r
}

// Test GET /users
func TestGetUsers(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John") // Check if "Alice" exists in response
}

// Test GET /users/:id (valid ID)
func TestGetUserByID_Valid(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// fmt.Println("test2", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John") // Check if "Alice" is returned
}

// Test GET /users/:id (invalid ID)
func TestGetUserByID_Invalid(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/users/99", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}

// Test POST /users
func TestCreateUser(t *testing.T) {
	router := setupRouter()

	newUserJSON := `{"name": "John Doe", "email": "john@example.com"}`
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte(newUserJSON)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "John Doe")
}

// Test PUT /users/:id (valid ID)
func TestUpdateUser_Valid(t *testing.T) {
	router := setupRouter()

	updatedUserJSON := `{"name": "John Ray", "email": "asd@fd.com"}`
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer([]byte(updatedUserJSON)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John Ray")
}

// Test PUT /users/:id (invalid ID)
func TestUpdateUser_Invalid(t *testing.T) {
	router := setupRouter()

	updatedUserJSON := `{"name": "John Ray", "email": "dsa@asd.com"}`
	req, _ := http.NewRequest("PUT", "/users/99", bytes.NewBuffer([]byte(updatedUserJSON)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}

// Test DELETE /users/:id (valid ID)
func TestDeleteUser_Valid(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User deleted")
}

// Test DELETE /users/:id (invalid ID)
func TestDeleteUser_Invalid(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("DELETE", "/users/99", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}
