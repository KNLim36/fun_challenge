# Fun challenge 

This project consists of two Go programs: `challengeA.go` and `challengeB.go`. The purpose of these programs is to generate a large file with random data and then process that file to classify the data.

## Overview

### `challengeA.go`

This program generates a text file containing random data objects, including alphabetical strings, numbers (both integers and real numbers), and alphanumeric strings with spaces. 

#### Key Features
- **File Generation**: Creates a file named `challengeA.txt` with random data.
- **Data Composition**: Includes objects like:
  - Random alphabetical strings (e.g., `abcde`)
  - Random integers (e.g., `12345`)
  - Random real numbers (e.g., `12.34`)
  - Alphanumeric strings with spaces at the front or back (e.g., ` a1b2  `)
- **File Size**: The file is approximately 10 MB in size.

### `challengeB.go`

This program reads the generated file `challengeA.txt`, processes its content, and classifies each data object into categories such as alphabetical, numerical, or alphanumeric. The results are saved to a new file named `challengeB.txt`.

#### Key Features
- **File Reading**: Reads data from `challengeA.txt`.
- **Data Classification**: Determines the type of each data object based on its content.
  - **Alphabetical**: Contains letters (either lowercase or uppercase).
  - **Integral**: Contains digits, possibly representing integers.
  - **Real**: Contains a decimal point.
  - **Alphanumerical**: Includes letters and digits, potentially with spaces.
- **File Writing**: Outputs classified data to `challengeB.txt` and optionally to the console.

## Usage

### Running `challengeA.go` and `challengeB.go`
1. Ensure you have Go installed on your system.
2. `cd` into respective directories.
   ```sh
   cd .\challengeA\
   ```
   or
   ```sh
   cd .\challengeB\
   ```
4. Run the script using these commands:
   ```sh
   go run challengeA.go
   ```
   or
   ```sh
   go run challengeB.go
   ```

### Note
For `challengeB.go`, a generated `challengeA.txt` file has to be in the `challengeB` directory, you can obtain that from copying the generated file from the `challengeA` directory.
