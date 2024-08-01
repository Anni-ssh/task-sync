package handler

import "strconv"

func parseQueryInt(value string) int {
	if value == "" {
		return 0
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return parsedValue
}
