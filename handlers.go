package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	maxNameLength = 50
	minNameLength = 2
)

func bookClassHandler(c *gin.Context) {
	var req struct {
		UserID  uuid.UUID `json:"user_id"`
		ClassID int       `json:"class_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Printf("Received booking request: UserID=%s, ClassID=%d", req.UserID, req.ClassID)

	class, err := bookClass(req.UserID, req.ClassID)
	if err != nil {
		if err.Error() == "class is full, added to waiting list" || err.Error() == "waitlist is full" {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		} else if err.Error() == "user not found" || err.Error() == "class not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Booking successful",
			"class_id": class.ID,
			"time":     class.StartTime.Format(time.RFC3339),
		})
	}
}

func cancelBookingHandler(c *gin.Context) {
	var req struct {
		UserID  uuid.UUID `json:"user_id"`
		ClassID int       `json:"class_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Printf("Received cancellation request: UserID=%s, ClassID=%d", req.UserID, req.ClassID)

	if err := cancelBooking(req.UserID, req.ClassID); err != nil {
		if err.Error() == "user not found" || err.Error() == "class not found" || err.Error() == "booking not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Cancellation successful"})
	}
}

func addUserHandler(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Name) < minNameLength {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name must be at least 2 characters long"})
		return
	}

	if len(req.Name) > maxNameLength {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is too long"})
		return
	}

	user := addUser(req.Name)
	logger.Printf("Added new user: ID=%s, Name=%s", user.ID, user.Name)
	c.JSON(http.StatusOK, gin.H{"message": "User added successfully", "user": user})
}
