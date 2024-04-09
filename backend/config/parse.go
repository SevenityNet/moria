package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getBoolEnv(key string) (bool, error) {
	valStr, exists := os.LookupEnv(key)
	if !exists {
		return false, fmt.Errorf("%s not set", key)
	}

	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return false, fmt.Errorf("%s: invalid boolean value", key)
	}

	return val, nil
}

func getIntEnv(key string) (int, error) {
	valStr, exists := os.LookupEnv(key)
	if !exists {
		return 0, fmt.Errorf("%s not set", key)
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("%s: invalid integer value", key)
	}

	return val, nil
}

func getStringEnv(key string, validValues []string) (string, error) {
	valStr, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("%s not set", key)
	}

	if len(validValues) > 0 {
		valid := false
		for _, v := range validValues {
			if valStr == v {
				valid = true
				break
			}
		}
		if !valid {
			return "", fmt.Errorf("%s: invalid value. Must be one of %v", key, validValues)
		}
	}

	return valStr, nil
}

func getStringArrayEnv(key, delimiter string) ([]string, error) {
	valStr, exists := os.LookupEnv(key)
	if !exists {
		return []string{}, fmt.Errorf("%s not set", key)
	}

	return strings.Split(valStr, delimiter), nil
}
