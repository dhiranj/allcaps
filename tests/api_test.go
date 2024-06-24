package main_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock HTTP server
var serverURL = "http://localhost:8080"

func TestBookClassAPI(t *testing.T) {
	userID := createUser(t, "Test User")
	classID := 3

	// Successful Booking
	resp := makePOSTRequest(t, "/book", map[string]interface{}{
		"user_id":  userID,
		"class_id": classID,
	})
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")

	// User Already Booked
	resp = makePOSTRequest(t, "/book", map[string]interface{}{
		"user_id":  userID,
		"class_id": classID,
	})
	assert.Equal(t, http.StatusConflict, resp.StatusCode, "Expected status code 409")

	userID1 := createUser(t, "User 1")
	userID2 := createUser(t, "User 2")
	userID3 := createUser(t, "User 3")

	classID = 1

	resp = makePOSTRequest(t, "/book", map[string]interface{}{
		"user_id":  userID1,
		"class_id": classID,
	})
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200 for User 1")

	resp = makePOSTRequest(t, "/book", map[string]interface{}{
		"user_id":  userID2,
		"class_id": classID,
	})
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200 for User 2")

	resp = makePOSTRequest(t, "/book", map[string]interface{}{
		"user_id":  userID3,
		"class_id": classID,
	})
	assert.Equal(t, http.StatusConflict, resp.StatusCode, "Expected status code 409 for User 3 waitlisted")

	resp = makePOSTRequest(t, "/cancel", map[string]interface{}{
		"user_id":  userID2,
		"class_id": classID,
	})
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200 for User 2 cancellation")

	resp = makePOSTRequest(t, "/cancel", map[string]interface{}{
		"user_id":  userID3,
		"class_id": classID,
	})
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200 for User 3 cancellation")

}

func createUser(t *testing.T, name string) string {
	// Create a new user
	resp := makePOSTRequest(t, "/add_user", map[string]interface{}{
		"name": name,
	})
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")

	// Parse the response body to get the user ID
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	userMap, ok := data["user"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to get user from response")
	}

	userID, ok := userMap["id"].(string)
	if !ok {
		t.Fatalf("Failed to get user ID from response")
	}

	return userID
}

func makePOSTRequest(t *testing.T, path string, data map[string]interface{}) *http.Response {
	payload, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal request data: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, serverURL+path, bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	return resp
}
