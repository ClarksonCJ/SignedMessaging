/// Package to host a http endpoint which returns  a signed message object on GET
/// Verifies signed message on POST

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/braintree/manners"
)

type myHandler struct {
}

type message struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Telephone string `json:"phone"`
}

func (h myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling Request")
	testKey := "ThisIsATest"

	switch r.Method {
	case "GET":
		// Get predefined message object
		msg := getMessage()
		log.Printf("Message Obj: %s\n", msg)
		// Compute HMAC-SHA256
		hash := msg.compute(testKey)
		// Set Header for HMAC Signature
		w.Header().Set("X-HMAC-Signature", base64.StdEncoding.EncodeToString(hash))
		// Set Status Code
		w.WriteHeader(http.StatusOK)
		// Write Message Data to response body
		w.Write([]byte(msg))
	case "POST":
		// Read Message Body from Req, log error if raised
		msgJSON, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
		}
		if len(msgJSON) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Message Body length is zero")
			fmt.Fprintf(w, "No Message Body")
			return
		}

		// Get signature from request header
		recSig := r.Header.Get("X-HMAC-Signature")
		log.Printf("Signature: %s\n", recSig)
		if recSig == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("No Hmac Signature in Header")
			fmt.Fprintf(w, "No Hmac Signature in Header")
			return
		}
		log.Printf("Header Signature: %s\n", recSig)
		// Decode Base64 header value
		recSigBytes, _ := base64.StdEncoding.DecodeString(recSig)
		// Validate signature against payload and respond accordingly
		result := validateMessagePost(string(msgJSON), testKey, recSigBytes)
		switch result {
		case true:
			w.WriteHeader(http.StatusOK)
			log.Println("Signature Matched")
			fmt.Fprintf(w, "Request Successful")
		case false:
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Signature Mismatch")
			fmt.Fprintf(w, "Request Failed")
		}
	}
}

func main() {
	log.Println("HMAC Signature handler starting up")
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, os.Kill)
		<-sigchan
		log.Println("Shutting Down")
		manners.Close()
	}()
	log.Fatal(manners.ListenAndServe(":8181", myHandler{}))
}

func getMessage() data {
	msg := message{Name: "Chris Clarkson",
		Email:     "chris.clarkson@hitachicapital.co.uk",
		Telephone: "0123456789"}
	if dataObj, err := json.Marshal(msg); err != nil {
		panic("Unable to parse object to JSON")
	} else {
		return data(dataObj)
	}
}

func validateMessagePost(bodyMsg string, key string, recSig []byte) bool {
	log.Printf("Check Key: %s\n", key)
	log.Printf("Check Sig: %s\n", base64.StdEncoding.EncodeToString(recSig))
	checkObj := data(bodyMsg)
	retval := checkObj.compare(key, recSig)
	return retval
}
