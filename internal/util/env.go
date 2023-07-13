package util

import (
	"fmt"
	"os"
	"strings"
)

func GetOsEnv(envKey string) string {
	return strings.TrimSpace(os.Getenv(envKey))
}

func CheckRequiredEnv() error {
	if value := GetOsEnv(EnvStackId); len(value) == 0 {
		return fmt.Errorf("invalid %s, please check your lambda environment variables", EnvStackId)
	}
	if value := GetOsEnv(EnvDdbTable); len(value) == 0 {
		return fmt.Errorf("invalid %s, please check your lambda environment variables", EnvDdbTable)
	}
	return nil
}
