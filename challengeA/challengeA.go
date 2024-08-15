package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
)

// Create a string builder to construct
var str strings.Builder
var alphabets = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var digits = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var alphanumerics = append(alphabets, digits...)

// // Smaller test sizes for quicker iterations
// const (
// 	fileSizeLimit = 1024 * 30 // Approximately 30 KB
// 	bufferSize    = 4096
// )

// Hard limit for file size: 10 MB = 10 * 1024 KB = 10 * 1024 * 1024 bytes
const (
	fileSizeLimit = 10 * 1024 * 1024 // 10 MB = 10 * 1024 KB = 10 * 1024 * 1024 bytes
	bufferSize    = 65536
)

const writeFrequencyNeeded = fileSizeLimit / bufferSize

func main() {
	var currentByteCount int = 0
	var currentWriteFrequency int = 0
	var data []byte

	// Creates a file first
	fileName := "challengeA.txt"
	file, errs := os.Create(fileName)
	if errs != nil {
		fmt.Println("Failed to create file:", errs)
		return
	}

	// Make sure the file closes
	defer file.Close()

	for {
		// https://www.kelche.co/blog/go/golang-bufio/
		writer := bufio.NewWriterSize(file, bufferSize)

		// We build the random object as byte
		object, byteSize := getRandObject()

		// and record its byte count (to write in 1 go)
		currentByteCount += byteSize

		// if the random object is the first object, we don't need comma
		if currentWriteFrequency == 0 && len(data) == 0 {
			data = append(data, object...)
		} else {
			// we have to consider the comma byte size as well before adding in comma
			comma := []byte(",")
			currentByteCount += len(comma)

			// add in both comma and the random object byte
			data = append(data, comma...)
			data = append(data, object...)
		}

		// If the currentByteCount is larger than the allowed bufferSize
		if currentByteCount >= bufferSize {
			// We write to the memory buffer
			_, errs := writer.Write(data)
			if errs != nil {
				fmt.Println("Failed to write to memory:", errs)
				return
			}

			// followed by writing to the actual file
			errs = writer.Flush()
			if errs != nil {
				fmt.Println("Failed to write to file:", errs)
				return
			}

			// reset the byte counter, byte slice and increment write frequency
			currentByteCount = 0
			data = data[:0]
			currentWriteFrequency++
		}

		// If we're about to exceed the frequency needed, stop
		if currentWriteFrequency >= writeFrequencyNeeded {
			break
		}
	}

	fmt.Println("Wrote to file", fileName, ".")
}

func getRandObject() ([]byte, int) {
	// Create a map to associate random index with object generation function
	randFuncMap := map[int]func(expectedLength int) ([]byte, int){
		0: getAlphabeticalObject,
		1: getRealNumberObject,
		2: getIntegerObject,
		3: getAlphanumericObjectWithSpace,
	}

	// Pick a function with a random index
	randomIndex := rand.IntN(len(randFuncMap))

	// We expect the average string length to be 5
	// e.g. absde, 12345, 1.234, a123d
	return randFuncMap[randomIndex](5)
}

func getAlphabeticalObject(expectedLength int) ([]byte, int) {
	// Iterate over the desired length of the string
	for i := 0; i < expectedLength; i++ {
		// Generate a random index to select a letter from the alphabet
		index := rand.IntN(len(alphabets))

		// 50% upper case, 50% lower case
		if coinFlip() {
			str.WriteString(strings.ToUpper(alphabets[index]))
		} else {
			str.WriteString(alphabets[index])
		}
	}

	result, length := getStringAndLengthFromBuilder()

	// Before resetting the str builder
	str.Reset()

	// Return the byte slice and its length
	return []byte(result), length
}

// Get an alphanumeric object with a random number of spaces before and after it
func getAlphanumericObjectWithSpace(expectedLength int) ([]byte, int) {
	// Generate a random number of spaces between 1 and 10
	spaceCount := rand.IntN(10) + 1

	// Randomly distribute spaces between the beginning and end of the string
	frontCount := rand.IntN(spaceCount + 1)
	backCount := spaceCount - frontCount

	// Add leading spaces
	addSpaces(frontCount)

	hasDigit := false
	hasAlphabet := false

	// Generate the alphanumeric part of the string
	for i := 0; i < expectedLength; i++ {
		index := rand.IntN(len(alphanumerics))

		// This is used to check if an alphabet or digit is getting created
		isAlphabet := index <= len(alphabets)-1
		isDigit := index >= len(alphabets)

		if isAlphabet {
			hasAlphabet = true
			// Coin flip to upper case
			if rand.IntN(1) == 1 {
				str.WriteString(strings.ToUpper(alphanumerics[index]))
			} else {
				str.WriteString(alphanumerics[index])
			}
			continue
		}

		if isDigit {
			hasDigit = true
			str.WriteString(alphanumerics[index])
			continue
		}

	}

	if !hasAlphabet {
		// Do something to make sure we add alphabet
		index := rand.IntN(len(alphabets))
		str.WriteString(alphabets[index]) // Adding in random alphabet
	}

	if !hasDigit {
		// Do something to make sure we add digit
		index := rand.IntN(len(digits))
		str.WriteString(digits[index]) // Adding in random digit
	}

	// Add trailing spaces
	addSpaces(backCount)

	// Return the string result as byte and length
	result, length := getStringAndLengthFromBuilder()
	str.Reset()
	return []byte(result), length
}

// Generates a random real number string with a given expected length
func getRealNumberObject(expectedLength int) ([]byte, int) {
	// Determine if the number should be negative
	// isNegative := coinFlip()
	isNegative := false

	// Track whether a non-zero digit has been encountered
	var containsNonZeroDigit bool = false

	// Track whether a decimal point has been added
	var hasAddedDecimal bool = false
	var decimalIndex int

	for i := 0; i < expectedLength; i++ {

		// [0] write and negative, should add "-"
		if i == 0 && isNegative {
			str.WriteString("-")
			continue
		}

		// Handle adding decimal/decimal
		if !hasAddedDecimal {
			if isNegative {

				/* Explanation:
				- If we're at the third digit (i == 2) and haven't added a decimal yet, and there's a zero digit before, add a decimal.
				- If we're past the third digit (i >= 2) and a random coin flip is successful, add a decimal.
				- If we're at the second-to-last position (i == expectedLength - 2) and haven't added a decimal yet, add one.
				Why -2 instead of -1? We can't place a decimal at the last place like so: "-012."
				*/
				shouldAddDecimal := (i == 2 && !containsNonZeroDigit) || (i >= 2 && coinFlip()) || (i == expectedLength-2)
				if shouldAddDecimal {
					str.WriteString(".")
					hasAddedDecimal = true
					decimalIndex = i
					continue
				}
			} else {
				// Only add decimal to non-first/last indexes
				if i != 0 && i != expectedLength-1 {

					/* Explanation:
					- If a random coin flip is successful, add a decimal.
					- If we're at the second-to-last position (i == expectedLength - 2) and haven't added a decimal yet, add one.
					Why -2 instead of -1? We can't place a decimal at the last place like so: "1234."
					*/
					shouldAddDecimal := (coinFlip()) || (i == expectedLength-2)
					if shouldAddDecimal {
						str.WriteString(".")
						hasAddedDecimal = true
						decimalIndex = i
						continue
					}
				}
			}
		}

		// Choose a number from 0 - 9
		index := rand.IntN(len(digits))
		str.WriteString(digits[index])

		// Update the containsNonZeroDigit flag if the digit is not zero
		if index != 0 {
			containsNonZeroDigit = true
		}

	}

	var result string

	// If we generated all 0s for expectedLength, like 0.000 or -0.00
	if !containsNonZeroDigit {
		str.Reset()
		str.WriteString("0") // Destined to be an integer
	}

	// Get the current string
	result, _ = getStringAndLengthFromBuilder()

	// If value is positive and the decimal index is not the first possible index
	// Trim left to convert values like 07.08 -> 7.08
	if !isNegative && decimalIndex != 1 {
		result = strings.TrimLeft(result, "0")
	}

	length := len(result)
	str.Reset()

	return []byte(result), length
}

func getIntegerObject(expectedLength int) ([]byte, int) {
	// Determine a maximum value based on expectedLength
	maxValue := int(math.Pow10(expectedLength-1)) - 1

	// Generate a random integer within the range
	randomInt := rand.IntN(2*maxValue+1) - maxValue

	// Convert integer to string
	result := strconv.Itoa(randomInt)

	return []byte(result), len(result)
}

// Returns a boolean value, simulating a coin flip (true or false with equal probability
func coinFlip() bool {
	return rand.IntN(2) == 0
}

// Adds a specified number of spaces to the string builder
func addSpaces(count int) {
	for i := 0; i < count; i++ {
		str.WriteString(" ")
	}
}

// Converts the string builder to a string and returns it along with its length
func getStringAndLengthFromBuilder() (string, int) {
	result := str.String()
	length := str.Len()

	return result, length
}
