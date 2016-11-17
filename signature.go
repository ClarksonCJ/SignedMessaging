package main

import (
	"crypto/hmac"
	"crypto/sha256"
)

type data string

func (d *data) compute(key string) []byte {
	// Return nil if no key passed for computation
	if len(key) == 0 {
		return nil
	}

	// return nil if no value set to attached object
	if len(*d) == 0 {
		return nil
	}
	message := []byte(*d)

	mac := hmac.New(sha256.New, []byte(key))
	mac.Write(message)
	return mac.Sum(nil)
}

func (d *data) compare(key string, generated []byte) bool {
	// Return false when no key passed as parameter
	if len(key) == 0 {
		return false
	}
	// Return false when no generated hash passed for comparison
	if len(generated) == 0 {
		return false
	}
	currentHash := d.compute(key)
	return hmac.Equal(currentHash, generated)
}
