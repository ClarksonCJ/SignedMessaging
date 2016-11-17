package main

import (
	"testing"
)

/// If an Empty key is passed to the hash computation function it should panic
func TestComputeMethodReturnsNilWithEmptyKey(t *testing.T) {
	testObj := data("Test Object")
	retval := testObj.compute("")
	if retval != nil {
		t.Errorf("Expected a nil response with empty key, received: %s", retval)
	}
}

/// If the attached object has no content, computation should return nil
func TestComputeMethodReturnsNilIfAttachedStructHasNoContent(t *testing.T) {
	testObj := data("")
	retval := testObj.compute("ThisisATestKey")
	if retval != nil {
		t.Errorf("Expected a nil response when object has no content, received: %s", retval)
	}
}

/// Check that returned hash value is of a fixed length of 32
func TestLengthOfReturnedHashIs32(t *testing.T) {
	testObj := data("{\"name\":\"Chris Clarkson\",\"email\":\"chris.clarkson@hitachicapital.co.uk\",\"phone\":\"0123456789\"}")
	retval := testObj.compute("ThisisATestKey")
	if len(retval) != 32 {
		t.Errorf("Expected the hash length to be 32 bytes, received: %d", len(retval))
	}
}

/// Ensure that a negative result is returned if no Key is passed to the compare function
func TestCompareMethodReturnsFalseWhenNoKeyIsPassed(t *testing.T) {
	testObj := data("{\"name\":\"Chris Clarkson\",\"email\":\"chris.clarkson@hitachicapital.co.uk\",\"phone\":\"0123456789\"}")
	hash := testObj.compute("ThisIsATestKey")
	result := testObj.compare("", hash)
	if result == true {
		t.Errorf("Expected False result when empty key passed to compare method, received: %t", result)
	}
}

/// Emsure that a Negative result is returned if no hash is passed for comparison
func TestCompareMethodReturnsFalseWhenNoHashIsPassedForComparison(t *testing.T) {
	testObj := data("{\"name\":\"Chris Clarkson\",\"email\":\"chris.clarkson@hitachicapital.co.uk\",\"phone\":\"0123456789\"}")
	result := testObj.compare("ThisIsATestKey", []byte(""))
	if result == true {
		t.Errorf("Expected False result when empty hash passed to compare method, received: %t", result)
	}
}

/// Ensure that a negative result is returned when no parameters are passed
func TestCompareMethodReturnsFalseWhenNoHashOrKeyPassed(t *testing.T) {
	testObj := data("{\"name\":\"Chris Clarkson\",\"email\":\"chris.clarkson@hitachicapital.co.uk\",\"phone\":\"0123456789\"}")
	result := testObj.compare("", []byte(""))
	if result == true {
		t.Errorf("Expected False result when empty hash and no key passed to compare method, received: %t", result)
	}
}

/// Ensure tht a posibive result is returned when a key and hash are passed which validate successfully.
func TestCompareMethodReturnsTrueWhenKeyAndGeneratedHashArePassedAndValidateSuccessfully(t *testing.T) {
	testObj := data("{\"name\":\"Chris Clarkson\",\"email\":\"chris.clarkson@hitachicapital.co.uk\",\"phone\":\"0123456789\"}")
	hash := testObj.compute("ThisIsATestKey")
	result := testObj.compare("ThisIsATestKey", hash)
	if result == false {
		t.Errorf("Expected treu result when key and hash passed to compare method and validate successfully, received: %t", result)
	}
}
