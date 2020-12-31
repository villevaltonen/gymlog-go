package app

import (
	"bytes"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	// correct credentials
	var jsonStr1 = []byte(`{"username":"user1", "password": "password1"}`)
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonStr1))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	// incorrect password
	var jsonStr2 = []byte(`{"username":"user1", "password": "passwordnotcorrect"}`)
	req, _ = http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonStr2))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

	// incorrect username
	var jsonStr3 = []byte(`{"username":"usernotfound", "password": "password1"}`)
	req, _ = http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonStr3))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestHeartbeat(t *testing.T) {
	// without JWT
	req1, _ := http.NewRequest("GET", "/api/heartbeat", nil)
	response1 := executeRequest(req1)
	checkResponseCode(t, http.StatusUnauthorized, response1.Code)

	// with JWT
	req3, _ := http.NewRequest("GET", "/api/heartbeat", nil)
	req3.AddCookie(authenticate())
	response3 := executeRequest(req3)
	checkResponseCode(t, http.StatusOK, response3.Code)
}

func authenticate() *http.Cookie {
	var jsonStr = []byte(`{"username":"user1", "password": "password1"}`)
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonStr))
	response := executeRequest(req)
	cookie := response.Result().Cookies()[0]
	return cookie
}
