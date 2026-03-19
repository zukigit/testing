package lib

import "fmt"

func GetEnv(m map[string]string, key, fallback string) string {
	if v, ok := m[key]; ok {
		return v
	}
	return fallback
}

func CheckEmptyValues(m map[string]string, keys []string) error {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			if v == "" {
				return fmt.Errorf("environment variable '%s' is empty", key)
			}
		} else {
			return fmt.Errorf("environment variable '%s' is not set", key)
		}
	}
	return nil
}
