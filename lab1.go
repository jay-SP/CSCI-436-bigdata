package main

import (
	"fmt"
	"strconv"
)

func compressMessage(msg string) string {
	compressedMsg := ""
	currentChar := msg[0]
	count := 1

	for i := 1; i < len(msg); i++ {
		if msg[i] == currentChar {
			count++
		} else {
			if count > 1 {
				compressedMsg += string(currentChar) + strconv.Itoa(count)
			} else {
				compressedMsg += string(currentChar)
			}
			currentChar = msg[i]
			count = 1
		}
	}

	// Add the last character's count
	if count > 1 {
		compressedMsg += string(currentChar) + strconv.Itoa(count)
	} else {
		compressedMsg += string(currentChar)
	}

	return compressedMsg
}

func main() {
	// sample Input
	msg1 := "abcaaabbb"
	msg2 := "abbcccddd"

	//Output
	fmt.Println(compressMessage(msg1))
	fmt.Println(compressMessage(msg2))
}
