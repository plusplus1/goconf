package goconf

import (
	"strings"
)

const (
	ERRNO_OK               int = 0
	ERRNO_ADD_KEY_EMPTY    int = 1
	ERRNO_ADD_TYPE_INVALID int = 2
)
const (
	key_trim_cutset = "/ \t"
)

type Configuration map[string]interface{}

type confError struct {
	Errno   int
	Message string
}

// Create an empty configuration
func NewConfiguration() Configuration {
	cfg := make(map[string]interface{})
	return cfg
}

// Add an element to configuration
// 	@param: key string, use "/" seprator
//	@param: value interface, may be any data type
func (cfg Configuration) Add(key string, value interface{}) Configuration {
	listKeyPath := strings.Split(strings.Trim(key, key_trim_cutset), "/") // split the key path
	count := len(listKeyPath)

	if 0 == count { // When the key is invalid, return directely, add element fail
		return cfg
	}

	if 1 == count { // Add key value directely, and return
		cfg[listKeyPath[0]] = value
		return cfg
	}

	// Add key value step by step
	var tempCfg Configuration = cfg
	var i int = 0
	for i = 0; i+1 < count; i++ {
		val := tempCfg.Get(listKeyPath[i])
		if nil != val {
			tempVal, isOk := val.(Configuration)
			if isOk {
				tempCfg = tempVal
				continue
			}
			return cfg // Value type invalid, add fail
		} else {
			tempVal := NewConfiguration()
			tempCfg[listKeyPath[i]] = tempVal
			tempCfg = tempVal
		}
	}
	tempCfg[listKeyPath[i]] = value
	return cfg
}

// Get an element from configuration
func (cfg Configuration) Get(key string) interface{} {
	// split the key path
	listKeyPath := strings.Split(strings.Trim(key, key_trim_cutset), "/")
	count := len(listKeyPath)
	if 0 == count {
		return nil
	}
	if 1 == count {
		val, ok := cfg[key]
		if ok {
			return val
		}
		return nil
	}

	// find step by step
	var tempCfg Configuration = cfg
	var i int = 0
	for i = 0; i+1 < count; i++ {
		val, ok := tempCfg[listKeyPath[i]]
		if !ok {
			return nil
		}
		tempVal, typeOk := val.(Configuration)
		if !typeOk {
			return nil
		}
		tempCfg = tempVal
	}
	val, ok := tempCfg[listKeyPath[i]]
	if ok {
		return val
	}
	return nil
}
