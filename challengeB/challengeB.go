package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

var alphabets = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var capitalAlphabets = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
var digits = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var file *os.File

const shouldOutputToConsole = false
const shouldSaveToFile = true
const challengeAFileName = "challengeA.txt"
const challengeBFileName = "challengeB.txt"

func main() {
	// Get file from challenge A
	data, errs := ioutil.ReadFile(challengeAFileName)
	if errs != nil {
		fmt.Println("Failed to read file size:", errs)
		return
	}

	content := string(data)
	allStrings := strings.Split(content, ",")

	if shouldSaveToFile {
		// Creates a file
		fileName := challengeBFileName
		file, errs = os.Create(fileName)
		if errs != nil {
			fmt.Println("Failed to create file:", errs)
			return
		}

		// Make sure the file closes
		defer file.Close()
	}

	for i := 0; i < len(allStrings); i++ {
		determineAndOutputObjectType(allStrings[i])
	}

	defer fmt.Println("Wrote to file", challengeBFileName, ".")
}

func determineAndOutputObjectType(input string) {
	isConfirmed := false

	// Determine the middle index of the string, rounding up if necessary
	halfLength := int(math.Ceil(float64(len(input)) / 2))

	// Split the string into two halves based on the calculated half length
	firstHalf := input[:halfLength]  // Slice from the beginning to the half length
	secondHalf := input[halfLength:] // Slice from the half length to the end

	var currentType string

	if !isConfirmed {
		// Check first half
		for i := 0; i < len(firstHalf); i++ {
			firstChar := string(firstHalf[i])
			currentType, isConfirmed = verifyObjectType(firstChar)

			if isConfirmed {
				if shouldOutputToConsole {
					fmt.Printf("Object: '%s' is %s\n", strings.Trim(input, " "), currentType)
				}

				if shouldSaveToFile {
					fmt.Fprintf(file, "Object: '%s' is %s\n", strings.Trim(input, " "), currentType)

				}
				return
			}

		}

		// Check second half
		for i := len(secondHalf) - 1; i >= 0; i-- {
			lastChar := string(secondHalf[i])
			currentType, isConfirmed = verifyObjectType(lastChar)

			if isConfirmed {
				if shouldOutputToConsole {
					fmt.Printf("Object: '%s' is %s\n", strings.Trim(input, " "), currentType)
				}

				if shouldSaveToFile {
					fmt.Fprintf(file, "Object: '%s' is %s\n", strings.Trim(input, " "), currentType)

				}
				return
			}

		}
	}

	// If it's still undetermined, must be an integer in this case.
	if shouldOutputToConsole {
		fmt.Printf("Object: '%s' is %s\n", strings.Trim(input, " "), currentType)
	}

	if shouldSaveToFile {
		fmt.Fprintf(file, "Object: '%s' is %s\n", strings.Trim(input, " "), currentType)

	}
}

func verifyObjectType(input string) (objectType string, isConfirmed bool) {
	// Check the easy ones: alphanumerical and real numbers first

	// If it has a space, definitely an alphanumerics
	if input == " " {
		return "alphanumerical", true
	}

	// If it has a decimal, definitely real numbers
	if input == "." {
		return "real", true
	}

	// If it has alphabets (small or capital), definitely alphabets (since alphanumerics are out)
	if containsString(alphabets, input) || containsString(capitalAlphabets, input) {
		return "alphabetical", false
	}

	// If it has digits, might be an integer, might be a real number
	if containsString(digits, input) {
		return "integral", false
	}

	return "unknown", false
}

func containsString(haystack []string, needle string) bool {
	// Iterate over each element in the haystack
	for _, element := range haystack {
		// Check if the current element matches the needle
		if element == needle {
			return true
		}
	}
	// If no match is found, return false
	return false
}
