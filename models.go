package main

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ClassType string

const (
	Yoga  ClassType = "yoga"
	Gym   ClassType = "gym"
	Dance ClassType = "dance"
)

type Class struct {
	ID        int       `json:"id"`
	Type      ClassType `json:"type"`
	Capacity  int       `json:"capacity"`
	StartTime time.Time `json:"start_time"`
}

type Booking struct {
	UserID  uuid.UUID `json:"user_id"`
	ClassID int       `json:"class_id"`
}

type WaitingList struct {
	UserID  uuid.UUID `json:"user_id"`
	ClassID int       `json:"class_id"`
}

var (
	users       = map[uuid.UUID]User{}
	classes     = map[int]Class{}
	bookings    = map[int][]Booking{}
	waitingList = map[int][]WaitingList{}

	usersMutex       sync.RWMutex
	classesMutex     sync.RWMutex
	bookingsMutex    sync.RWMutex
	waitingListMutex sync.RWMutex
)
