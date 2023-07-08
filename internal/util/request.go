package util

import (
	"errors"
	"strings"
)

func IsEligibleUser(requestPath string, stackId string) error {
	if IsBlank(stackId) {
		return errors.New("invalid stack id")
	}
	requestId := strings.ReplaceAll(requestPath, "/", "")
	if requestId != stackId {
		return errors.New("invalid request id")
	}
	return nil
}
