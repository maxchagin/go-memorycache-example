package memorycache

import (
	"testing"
	"time"
)

const (
	testKey      string = "cache:test"
	testKeyEmpty string = "cache:empty"
	testValue    string = "Hello Test"
)

// AppCache init new cache
var AppCache = New(10*time.Minute, 1*time.Hour)

// AppCache init new cache
var AppCacheGC = New(10*time.Minute, 1*time.Second)

// TestGet get cache by key
func TestGet(t *testing.T) {

	AppCache.Set(testKey, testValue, 1*time.Minute)

	value, found := AppCache.Get(testKey)

	if value != testValue {
		t.Error("Error: ", "The received value: do not correspond to the expectation:", value, testValue)
	}

	if found != true {
		t.Error("Error: ", "Could not get cache")
	}

	// get cache by key is empty
	value, found = AppCache.Get(testKeyEmpty)

	if value != nil || found != false {
		t.Error("Error: ", "Value does not exist and must be empty", value)
	}
}

// TestDelete delete cache by key
func TestDelete(t *testing.T) {

	AppCache.Set(testKey, testValue, 1*time.Minute)

	error := AppCache.Delete(testKey)

	if error != nil {
		t.Error("Error: ", "Cache delete failed")
	}

	value, found := AppCache.Get(testKey)

	if found {
		t.Error("Error: ", "Should not be found because it was deleted")
	}

	if value != nil {
		t.Error("Error: ", "Value is not nil:", value)
	}

	// repeat deletion of an existing cache
	error = AppCache.Delete(testKeyEmpty)

	if error == nil {
		t.Error("Error: ", "An empty cache should return an error")
	}

}
