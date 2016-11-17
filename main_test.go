package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

/// Test to ensure that the GET Request to the service returns a valid JSON message in the body of the response
/// Assertions to validate the correct values are Unmarshaled from the message body
func TestGETShouldReturnJsonMsgInBody(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler := myHandler{}

	handler.ServeHTTP(w, r)

	resp := message{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)

	if err != nil {
		t.Error(err)
	}

	if resp.Name != "Chris Clarkson" {
		t.Errorf("Name value not expected, Received: %s, Expected: Chris Clarkson", resp.Name)
	}

	if resp.Email != "chris.clarkson@hitachicapital.co.uk" {
		t.Errorf("Email Value not expected, recieved: %s, Expected: chris.clarkson@hitachicapital.co.uk", resp.Email)
	}

	if resp.Telephone != "0123456789" {
		t.Errorf("Telephone Value not expected, received: %s, Expected: 0123456789", resp.Telephone)
	}
}

/// Ensure that the GET Request returns a HMAC-SHA256 hash in the X-HMAC-Signature header
func TestGETShouldReturnXHMACSIGNATUREinHeader(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler := myHandler{}
	handler.ServeHTTP(w, r)
	log.Printf(w.Header().Get("X-HMAC-Signature"))
	if w.Header().Get("X-HMAC-Signature") == "" {
		t.Errorf("Expected HMAC Signature in Header")
	}
}

/// A successfully execute GET request to the service should acknowledge the request with a 200 status code in the response
func TestGETShouldReturnStatusCode200(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler := myHandler{}
	handler.ServeHTTP(w, r)
	statusCode := w.Code
	if statusCode != 200 {
		t.Errorf("Expected Status Code 200, received: %d", statusCode)
	}
}

/// a Post request without the HMAC signature in the header should response with status 400
func TestPOSTShouldReturn400ErrorIfNotReceiveHMACSignatureInHeader(t *testing.T) {
	body := getMessage()
	r, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(string(body)))
	w := httptest.NewRecorder()

	handler := myHandler{}
	handler.ServeHTTP(w, r)

	statusCode := w.Code
	if statusCode != 400 {
		t.Errorf("Expected Status Code 400, received: %d", statusCode)
	}
}

/// POST Should return 200 if the HMAC-SHA26 signature in the header is validated against the request body
func TestPOSTShouldReturn200IfValidHashPresentedForBody(t *testing.T) {
	body := getMessage()
	r, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(string(body)))
	dataBody := data(body)
	hash := dataBody.compute("ThisIsATest")
	r.Header.Set("X-HMAC-Signature", base64.StdEncoding.EncodeToString(hash))
	w := httptest.NewRecorder()

	handler := myHandler{}
	handler.ServeHTTP(w, r)

	statusCode := w.Code
	if statusCode != 200 {
		t.Errorf("Expected Status Code 200, received: %d", statusCode)
	}
}

/// An Invalid HMAC-SHA256 signature in the header should return a 400 status code in the response
func TestPOSTShouldReturn400IfInvalidHashPresentedForBody(t *testing.T) {
	body := getMessage()
	r, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(string(body)))
	r.Header.Set("X-HMAC-Signature", "AU8MRaXUE/xTxv852nZSDWyq62HlKiWQ67R7vo1cd4s=")
	w := httptest.NewRecorder()

	handler := myHandler{}
	handler.ServeHTTP(w, r)

	statusCode := w.Code
	if statusCode != 400 {
		t.Errorf("Expected Status Code 400, received: %d", statusCode)
	}
}

/// A POST request with an empty body should return a 400 status code in the response
func TestPOSTShouldReturn400IfNoMessagebodyPresentedDuringRequest(t *testing.T) {
	body := ""
	r, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(string(body)))
	r.Header.Set("X-HMAC-Signature", "AU8MRaXUE/xTxv852nZSDWyq62HlKiWQ67R7vo1cd4s=")
	w := httptest.NewRecorder()

	handler := myHandler{}
	handler.ServeHTTP(w, r)

	statusCode := w.Code
	if statusCode != 400 {
		t.Errorf("Expected Status code 400, received: %d", statusCode)
	}
}

/// The GetMessage function should return a valid JSON object which can be unmarshalled and provide the correct information
func TestGetMessageReturnsDataObject(t *testing.T) {
	testObj := getMessage()
	log.Println(testObj)
	result := message{}
	json.Unmarshal([]byte(testObj), &result)

	if result.Name != "Chris Clarkson" {
		t.Errorf("Name value not expected, Received: %s, Expected: Chris Clarkson", result.Name)
	}

	if result.Email != "chris.clarkson@hitachicapital.co.uk" {
		t.Errorf("Email Value not expected, recieved: %s, Expected: chris.clarkson@hitachicapital.co.uk", result.Email)
	}

	if result.Telephone != "0123456789" {
		t.Errorf("Telephone Value not expected, received: %s, Expected: 0123456789", result.Telephone)
	}
}

/// Validation function should return false if no Message value passed
func TestThatValidateMessagePostReturnsFalseIfNoMsgPassed(t *testing.T) {
	result := validateMessagePost("", "ThisIsATest", []byte("AU8MRaXUE/xTxv852nZSDWyq62HlKiWQ67R7vo1cd4s="))
	if result != false {
		t.Errorf("Expected false result, received: %t", result)
	}
}

/// Validation function should return false if key is not passed
func TestThatValidateMessagePostReturnsFalseIfNoKeyPassed(t *testing.T) {
	result := validateMessagePost("{\"name\":\"Chris Clarkson\",\"email\":\"chris.clarkson@hitachicapital.co.uk\",\"phone\":\"0123456789\"}", "", []byte("AU8MRaXUE/xTxv852nZSDWyq62HlKiWQ67R7vo1cd4s="))
	if result != false {
		t.Errorf("Expected false result, received: %t", result)
	}
}

/// Validation function should return false if no Hash is passed for comparison
func TestThatValidateMessagePostReturnsFalseIfNoComparisonHashPassed(t *testing.T) {
	result := validateMessagePost("{\"name\":\"Chris Clarkson\",\"email\":\"chris.clarkson@hitachicapital.co.uk\",\"phone\":\"0123456789\"}", "ThisIsATest", []byte(""))
	if result != false {
		t.Errorf("Expected false result, received: %t", result)
	}
}

/// Validation function should return true if validation processes successfully
func TestThatValidateMessagePostReturnsTrueIfAllComponentsValidateSuccessfully(t *testing.T) {
	hash, _ := base64.StdEncoding.DecodeString("xU8MRaXUE/XTxv76JnZSDWyq62HlKiWQ67R7vo1cd4s=")
	result := validateMessagePost("{\"name\":\"Chris Clarkson\",\"email\":\"chris.clarkson@hitachicapital.co.uk\",\"phone\":\"0123456789\"}",
		"ThisIsATest",
		hash)
	if result != true {
		t.Errorf("Expected true result, received: %t", result)
	}
}
