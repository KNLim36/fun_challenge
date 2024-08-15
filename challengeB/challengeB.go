package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

var alphabets = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var digits = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

func main() {
	// Get file from challenge A
	data, errs := ioutil.ReadFile("challengeA.txt")
	if errs != nil {
		fmt.Println("Failed to read file size:", errs)
		return
	}

	content := string(data)
	allStrings := strings.Split(content, ",")

	for i := 0; i < len(allStrings); i++ {
		checkAndOutputObjectType(allStrings[i])
	}
}

func checkAndOutputObjectType(input string) {
	isConfirmed := false

	// Calculate the length of the input string
	length := len(input)

	// Determine the middle index of the string, rounding up if necessary
	halfLength := int(math.Ceil(float64(length) / 2))

	// Split the string into two halves based on the calculated half length
	firstHalf := input[:halfLength]  // Slice from the beginning to the half length
	secondHalf := input[halfLength:] // Slice from the half length to the end

	var currentType string

	if !isConfirmed {
		for i := 0; i < len(firstHalf); i++ {
			firstChar := string(firstHalf[i])
			currentType, isConfirmed = verifyObjectType(firstChar)

			if isConfirmed {
				fmt.Printf("Object: '%s' has type %s\n", input, currentType)
				return
			}

		}

		for i := len(secondHalf) - 1; i >= 0; i-- {
			if secondHalf[i] != ' ' {
				lastChar := string(secondHalf[i])
				currentType, isConfirmed = verifyObjectType(lastChar)

				if isConfirmed {
					fmt.Printf("Object: '%s' has type %s\n", input, currentType)
					return
				}
			}
		}
	}

	fmt.Printf("Object: '%s' has type %s\n", input, currentType)
}

func verifyObjectType(input string) (objectType string, isConfirmed bool) {
	// Check the easy ones: alphanumerical and real numbers first

	// If it is an empty space, definitely an alphanumerics
	if input == " " {
		return "Alphanumerics", true
	}

	if input == "." {
		return "Real numbers", true
	}

	if containsString(alphabets, input) {
		return "Alphabets", false
	}

	if containsString(digits, input) {
		return "Integers", false
	}

	return "Type: Unknown", false
}

func containsString(haystack []string, needle string) bool {
	// Create a map to efficiently check for string existence
	lookupTable := make(map[string]bool)

	// Populate the lookup table with elements from the haystack
	for _, element := range haystack {
		lookupTable[element] = true
	}

	// Check if the needle exists in the lookup table
	return lookupTable[needle]
}
