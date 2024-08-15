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

var str strings.Builder
var alphabets = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var digits = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var alphanumerics = append(alphabets, digits...)

func main() {
	// Hard limit for file size: 10 MB = 10 * 1024 KB = 10 * 1024 * 1024 bytes
	const fileSizeLimit = 10 * 1024 * 1024
	const bufferSize = 65536
	// const fileSizeLimit = 1024 * 30
	// const bufferSize = 4096
	const writeFrequencyNeeded = fileSizeLimit / bufferSize
	var currentByteCount int = 0
	var currentWriteFrequency int = 0
	var data []byte

	// Creates a file first
	// fileName := filepath.Join("../", "challengeA.txt")
	fileName := "challengeA.txt"
	file, errs := os.Create(fileName)
	if errs != nil {
		fmt.Println("Failed to create file:", errs)
		return
	}
	defer file.Close()

	for {
		// https://www.kelche.co/blog/go/golang-bufio/
		writer := bufio.NewWriterSize(file, bufferSize)

		// We build the random object as string
		object, byteSize := getRandomObject()

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

	fmt.Println("Wrote to file ", fileName, ".")
}

// Get the random objects in bytes
func getRandomObject() ([]byte, int) {
	randMap := map[int]func(expectedLength int) ([]byte, int){
		0: getAlphabeticalObject,
		1: getRealNumberObject,
		2: getIntegerObject,
		3: getAlphanumericObjectWithSpace,
	}

	n := rand.IntN(4)

	// We expect the average string length to be 5
	// e.g. absde, 12345, 1.234, a123d
	return randMap[n](5)
}

func getAlphabeticalObject(expectedLength int) ([]byte, int) {
	// https://www.educative.io/answers/how-to-use-string-builder-in-golang
	for i := 0; i < expectedLength; i++ {
		index := rand.IntN(len(alphabets))

		// Add in randomizing between small and capital letters
		if rand.IntN(1) == 1 {
			str.WriteString(strings.ToUpper(alphabets[index]))
		} else {
			str.WriteString(alphabets[index])
		}
	}

	resultString := str.String()
	resultLength := str.Len()
	str.Reset()
	return []byte(resultString), resultLength
}

// Requirement: should contain a random number of spaces before and after it (not exceeding 10 spaces)
// We constraint the alphanumerical string length to 5 with the front back spaces
func getAlphanumericObjectWithSpace(expectedLength int) ([]byte, int) {
	// minimum 1 space, maximum 10 spaces
	spaceCount := rand.IntN(10) + 1

	// front space should be less than or equal to spaceCount
	frontCount := rand.IntN(spaceCount + 1)
	backCount := spaceCount - frontCount

	for i := 0; i < frontCount; i++ {
		str.WriteString(" ")
	}

	for i := 0; i < expectedLength; i++ {
		index := rand.IntN(len(alphanumerics))
		if rand.IntN(1) == 1 {
			str.WriteString(strings.ToUpper(alphanumerics[index]))
		} else {
			str.WriteString(alphanumerics[index])
		}
	}

	for i := 0; i < backCount; i++ {
		str.WriteString(" ")
	}

	resultString := str.String()
	resultLength := str.Len()

	str.Reset()
	return []byte(resultString), resultLength
}

// Real digits can be positive, negative, zero, fractions, decimals and irrational
// Skipping irrational digits since they have infinite length
// Skipping zero since it's possible from integers
// Skipping fractions since they can be represented with decimals
func getRealNumberObject(expectedLength int) ([]byte, int) {
	// 50% chance to be negative
	isNegative := rand.IntN(2) == 0

	// Check if non-zero showed up
	var hasNonZero bool = false
	var hasAddedPeriod bool = false
	var periodIndex int

	for i := 0; i < expectedLength; i++ {
		// choose a number from 0 - 9
		index := rand.IntN(len(digits))

		// [0] write and negative, should add "-"
		if i == 0 && isNegative {
			str.WriteString("-")
			continue
		}

		if !hasAddedPeriod {
			if isNegative {
				// If we reach -0* and has no period yet, should add "."
				if (i == 2 && !hasNonZero) ||
					// If we reach allowed indexes and has no period yet, roll 50% to add "."
					(i >= 2 && rand.IntN(2) == 0) ||
					// If we reach last possible index and has no period yet, should add "."
					(i == expectedLength-1) {
					str.WriteString(".")
					hasAddedPeriod = true
					periodIndex = i
					continue
				}
			} else {
				// Only add period to non-first/last indexes
				if i != 0 && i != expectedLength-1 {

					// If dice roll successful, or it's already last possible index
					if (i == 3) || (rand.IntN(2) == 0) {
						str.WriteString(".")
						hasAddedPeriod = true
						periodIndex = i
						continue
					}
				}
			}
		}

		str.WriteString(digits[index])

		if index != 0 {
			hasNonZero = true
		}

	}

	var resultString string
	// This is failing for values like 0.086 and -00.7

	// If we generated all 0s, write 0 instead
	if !hasNonZero {
		str.Reset()
		str.WriteString("0")
	}

	if isNegative {
		resultString = str.String()
		if periodIndex != 2 {
			resultString = strings.TrimLeft(str.String(), "0")
		}
	} else {
		resultString = str.String()

		// If period is at index 1, shouldn't trim left (or 0.2 will be .2)
		if periodIndex != 1 {
			resultString = strings.TrimLeft(str.String(), "0")
		}
	}

	resultLength := len(resultString)
	str.Reset()

	return []byte(resultString), resultLength
}

func getIntegerObject(expectedLength int) ([]byte, int) {
	// Determine a maximum value based on expectedLength
	maxValue := int(math.Pow10(expectedLength-1)) - 1

	// Generate a random integer within the range
	randomInt := rand.IntN(2*maxValue+1) - maxValue

	// Convert integer to string
	resultString := strconv.Itoa(randomInt)

	return []byte(resultString), len(resultString)
}
