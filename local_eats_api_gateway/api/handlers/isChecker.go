package handlers

import (
	"github.com/google/uuid"
	"strconv"
)

func IsNumber(number string) (int, error) {
	number1, err := strconv.Atoi(number)
	if err != nil {
		return 0, err
	}

	return number1, nil
}
func Parse(id string) bool {
	_, err := uuid.Parse(id)
	return !(err == nil)
}
