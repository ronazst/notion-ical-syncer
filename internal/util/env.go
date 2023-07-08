package util

import (
	"fmt"
	"os"
	"strings"
)

func GetOsEnv(envKey string) (string, error) {
	envValue := strings.TrimSpace(os.Getenv(envKey))
	if len(envValue) == 0 {
		return "", fmt.Errorf("invalid %s, please check your lambda environment variables", envKey)
	}
	return envValue, nil
}
