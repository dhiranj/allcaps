package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var logger *log.Logger

func main() {
	// Setup logging
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	r := gin.Default()

	// Initialize sample data
	initSampleData()

	go cleanUpWaitlist()

	r.POST("/book", bookClassHandler)
	r.POST("/cancel", cancelBookingHandler)
	r.POST("/add_user", addUserHandler) // New route for adding users

	r.Run(":8080")
}

func initSampleData() {
	// Initialize classes
	classes[1] = Class{ID: 1, Type: Yoga, Capacity: 2, StartTime: time.Now().Add(1 * time.Hour)}
	classes[2] = Class{ID: 2, Type: Gym, Capacity: 3, StartTime: time.Now().Add(2 * time.Hour)}
	classes[3] = Class{ID: 3, Type: Dance, Capacity: 1, StartTime: time.Now().Add(2 * time.Hour)}

	// Initialize users
	a := uuid.New()
	b := uuid.New()
	c := uuid.New()
	users[a] = User{ID: a, Name: "Alice"}
	users[b] = User{ID: b, Name: "Bob"}
	users[c] = User{ID: c, Name: "Charlie"}

	logger.Println("Initialized sample data")
}
