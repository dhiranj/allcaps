package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func bookClass(userID uuid.UUID, classID int) error {

	classesMutex.Lock()
	defer classesMutex.Unlock()
	usersMutex.RLock()
	defer usersMutex.RUnlock()
	bookingsMutex.Lock()
	defer bookingsMutex.Unlock()
	waitingListMutex.Lock()
	defer waitingListMutex.Unlock()

	class, exists := classes[classID]
	if !exists {
		logger.Printf("Class not found: ClassID=%d", classID)
		return errors.New("class not found")
	}

	if _, exists := users[userID]; !exists {
		logger.Printf("User not found: UserID=%s", userID)
		return errors.New("user not found")
	}

	// Check if the user has already booked this class
	for _, booking := range bookings[classID] {
		if booking.UserID == userID {
			logger.Printf("User already booked this class: UserID=%s, ClassID=%d", userID, classID)
			return errors.New("user has already booked this class")
		}
	}

	// Check if the user has already booked any class in the same time slot
	for _, b := range bookings {
		for _, booking := range b {
			if booking.UserID == userID {
				otherClass := classes[booking.ClassID]
				if otherClass.StartTime.Equal(class.StartTime) {
					logger.Printf("User has already booked another class in the same time slot: UserID=%s, ClassID=%d, OtherClassID=%d", userID, classID, booking.ClassID)
					return fmt.Errorf("user has already booked another class (%s) in this time slot", otherClass.Type)
				}
			}
		}
	}

	if len(bookings[classID]) >= class.Capacity {
		// Check if the waitlist has reached the capacity of the class
		if len(waitingList[classID]) >= class.Capacity {
			logger.Printf("Waitlist is full: ClassID=%d", classID)
			return errors.New("waitlist is full")
		}
		waitingList[classID] = append(waitingList[classID], WaitingList{UserID: userID, ClassID: classID})
		logger.Printf("Class at capacity, added user %s to waiting list for class %d", userID, classID)
		return errors.New("class is full, added to waiting list")
	}

	bookings[classID] = append(bookings[classID], Booking{UserID: userID, ClassID: classID})
	logger.Printf("Booked class %d for user %s", classID, userID)
	return nil
}

func cancelBooking(userID uuid.UUID, classID int) error {

	classesMutex.Lock()
	defer classesMutex.Unlock()
	usersMutex.RLock()
	defer usersMutex.RUnlock()
	bookingsMutex.Lock()
	defer bookingsMutex.Unlock()
	waitingListMutex.Lock()
	defer waitingListMutex.Unlock()

	class, exists := classes[classID]
	if !exists {
		logger.Printf("Class not found: ClassID=%d", classID)
		return fmt.Errorf("class not found")
	}

	if _, exists := users[userID]; !exists {
		logger.Printf("User not found: UserID=%s", userID)
		return errors.New("user not found")
	}

	if time.Until(class.StartTime) < 30*time.Minute {
		logger.Printf("Cancellation attempt within 30 minutes: UserID=%s, ClassID=%d", userID, classID)
		return fmt.Errorf("cannot cancel within 30 minutes of the class start time")
	}

	for i, booking := range bookings[classID] {
		if booking.UserID == userID {
			bookings[classID] = append(bookings[classID][:i], bookings[classID][i+1:]...)
			logger.Printf("Booking canceled: UserID=%s, ClassID=%d", userID, classID)

			if len(waitingList[classID]) > 0 {
				nextUser := waitingList[classID][0]
				waitingList[classID] = waitingList[classID][1:]
				bookings[classID] = append(bookings[classID], Booking{UserID: nextUser.UserID, ClassID: classID})
				notifyUser(nextUser.UserID, classID)
			}

			return nil
		}
	}

	logger.Printf("Booking not found for cancellation: UserID=%s, ClassID=%d", userID, classID)
	return fmt.Errorf("booking not found for user %v in class %v", userID, classID)
}

func notifyUser(userID uuid.UUID, classID int) {
	user, exists := users[userID]
	if exists {
		logger.Printf("Notification sent: UserID=%s, ClassID=%d", userID, classID)
		fmt.Printf("Notification: User %s has been moved from waitlist to booking for class %d.\n", user.Name, classID)
	}
}

func addUser(name string) User {

	usersMutex.Lock()
	defer usersMutex.Unlock()

	id := uuid.New()
	user := User{ID: id, Name: name}
	users[id] = user
	return user
}

func cleanUpWaitlist() {
	for {
		time.Sleep(10 * time.Minute) // Run every 10 minutes

		now := time.Now()
		for classID, class := range classes {

			if now.After(class.StartTime) {

				if len(waitingList[classID]) > 0 {
					logger.Printf("Cleaning up waitlist for class %d", classID)
					waitingList[classID] = []WaitingList{} // Clear the waitlist
				}
			}

		}
	}
}
