package model

import (
	"fmt"
	"github.com/ronazst/notion-ical-syncer/internal/util"
	"strings"
)

type QueryItem struct {
	DatabaseID   string
	DateFieldKey string
}

func ParseQueryItems(rawItems string) ([]QueryItem, error) {
	var result []QueryItem
	strItems := strings.Split(rawItems, ",")
	for i := 0; i < len(strItems); i++ {
		itemPair := strings.Split(strItems[i], "/")
		if len(itemPair) != 2 {
			return nil, fmt.Errorf("failed to read part of input, value: %v", strItems[i])
		}
		if util.IsBlank(itemPair[0]) || util.IsBlank(itemPair[1]) {
			return nil, fmt.Errorf("the input pair can't contains blank string, value: %v", strItems[i])
		}
		result = append(result, QueryItem{
			DatabaseID:   itemPair[0],
			DateFieldKey: itemPair[1],
		})
	}
	return result, nil
}
