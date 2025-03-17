package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Value string
}

func TestCache_AddAndGet(t *testing.T) {
	cache := NewCache[TestStruct]()

	// Test adding and getting a value
	key := "testKey"
	value := &TestStruct{Value: "testValue"}
	err := cache.Add(key, value)
	assert.NoError(t, err)

	retrievedValue, err := cache.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, retrievedValue)
}

func TestCache_GetNonExistentKey(t *testing.T) {
	cache := NewCache[TestStruct]()

	// Test getting a non-existent key
	key := "nonExistentKey"
	retrievedValue, err := cache.Get(key)
	assert.Error(t, err)
	assert.Nil(t, retrievedValue)
}

func TestCache_Remove(t *testing.T) {
	cache := NewCache[TestStruct]()

	// Test adding and removing a value
	key := "testKey"
	value := &TestStruct{Value: "testValue"}
	err := cache.Add(key, value)
	assert.NoError(t, err)

	cache.Remove(key)
	retrievedValue, err := cache.Get(key)
	assert.Error(t, err)
	assert.Nil(t, retrievedValue)
}
