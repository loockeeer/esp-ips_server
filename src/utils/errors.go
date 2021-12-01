package utils

import (
	"log"
	"strconv"
)

func Atoi(input string, message string) int {
	n, err := strconv.Atoi(input)

	if err != nil {
		log.Panicln(message)
	}
	return n
}
