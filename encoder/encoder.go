package encoder

// TODO: Optimize the return values, might be redundant return

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
)

type MessageStruct struct {
	matchingBytePattern []byte
	firstStart          int
	firstEnd            int
	secondStart         int
	secondEnd           int
	messageLength       int
	messageBytes        []byte
}

// Read x amount of bytes starting from byte y and return bytearray
func ExtractBytesFromFile(filePath string, bytesToRead int, startByte int) []byte {
	byteArray := make([]byte, bytesToRead) // Create a byteArray of max size byteLengthToHash
	openFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer openFile.Close()

	openFile.Seek(int64(startByte), io.SeekStart) // Jump to startByte
	openFile.Read(byteArray)                      // read the message array
	return byteArray
}

/// Decoding //////
//
//

/*
findAllMatchingPatterns looks for the matching byte pattern in the file
It scans for the first instance of the matching pattern, notes where it starts and ends
Then continues scanning until it finds the second matching pattern and notes where that starts and ends

It then returns the bytearray inbetween the matching patterns.
*/
func FindAllMatchingPatterns(filePath string, matchingBytePattern []byte, maxCaveSize int) MessageStruct {
	Mess := MessageStruct{}
	Mess.matchingBytePattern = matchingBytePattern
	firstStart := 0  // Where first matching pattern starts
	firstEnd := 0    // Where first matching pattern ends
	secondStart := 0 // Where second matching pattern starts
	secondEnd := 0   // Where second matching pattern ends
	readBytes := 0   // Amount of bytes read

	// Find start of codeCave
	firstStart, firstEnd, readBytes = FindMatchingPattern(filePath, matchingBytePattern, readBytes)
	if firstEnd == 0 { // Found no data
		return Mess // Return an empty message struct
	}

	secondStart, secondEnd, readBytes = FindMatchingPattern(filePath, matchingBytePattern, readBytes)

	if secondEnd == 0 { // Found no second matching pattern, data invalid
		return Mess
	}
	Mess.firstStart = firstStart
	Mess.firstEnd = firstEnd
	Mess.secondStart = secondStart
	Mess.secondEnd = secondEnd
	Mess.messageLength = ((secondStart - 1) - firstEnd)
	Mess.messageBytes = ExtractBytesFromFile(filePath, Mess.messageLength, (Mess.firstEnd)) // Extract message of length messageLength starting from the end first matching pattern + 1 byte
	return Mess

}

/*
findMatchingPattern scans a file for a matching pattern

It returns bytes where matching pattern starts and ends
Once pattern start and end is found it stops.

To find a second matching pattern the function is run again but starts reading at the byte number of the end of the last matching pattern.
Returns the index of where the pattern starts and ends and how many bytes were read
*/
func FindMatchingPattern(filePath string, matchingBytePattern []byte, readFromByte int) (int, int, int) {
	openFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer openFile.Close()

	openFile.Seek(int64(readFromByte), io.SeekStart) // Jump to readFromByte in file, if 0 then read from start
	scanner := bufio.NewScanner(openFile)
	scanner.Split(bufio.ScanBytes)

	countBytes := readFromByte // Counts number of bytes read + byte offset in readFromByte
	startMatch := 0
	endMatch := 0

	matchingBytePatternLength := len(matchingBytePattern)
	matchingByteIndex := 0 // Counter for how many bytes in file has matched the byte pattern
	var inCave bool        // Are we inside a potential code cave?
	for scanner.Scan() {
		countBytes += 1 // Increase the read bytes by one
		b := scanner.Bytes()
		if b[0] == matchingBytePattern[matchingByteIndex] { // if byte matches first byte in bytePattern, start comparison for rest of bytes

			matchingByteIndex += 1                              // Since the byte in the file matched the first byte of the pattern, change index to the next byte in the pattern
			if matchingByteIndex == matchingBytePatternLength { // Pattern has matched
				endMatch = countBytes
				return startMatch, endMatch, countBytes
			}
			if !inCave {
				startMatch = countBytes // First byte matched, set the startMatch variable where the beginning of the matching pattern is to number of bytes read
			}
			inCave = true
		} else { // The bytes are no longer the same, end of matching pattern
			inCave = false // Patterns were no longer matching, continuing read
			matchingByteIndex = 0
		}
	}
	return 0, 0, countBytes // Found no pattern
}

func EncodeMessage(inputFile string, outputFile string, inputMessage string, byteLengthToHash int, verbose bool, tryOut bool) {

	bytePatternToMatchLength := 8 // Size of the matching byte pattern before and after message
	printOutput := true
	byteArrayToHash := ReadBytesFromFile(inputFile, byteLengthToHash)                         // Creates an array of bytes of max length byteLengthToHash
	bytePatternToMatch := Md5HashBytesReturnXBytes(byteArrayToHash, bytePatternToMatchLength) // Creates an MD5sum and returns the first bytes defined by the length of bytePatternToMatchLength

	// Create hidden data, Gzip and reverse the inputMessage for obfuscation
	gZbyteMessage, err := GzipData([]byte(inputMessage))
	if err != nil {
		log.Fatal(err.Error())
	}

	revByteMessage := ReverseBytes(gZbyteMessage)
	hiddenMessage := CreateHiddenData(bytePatternToMatch, revByteMessage)

	caveSlice := FindCaves(inputFile, len(hiddenMessage))
	if len(caveSlice) == 0 {

		fmt.Fprintf(os.Stderr, "Error: No cave found in file\n")
		os.Exit(1)

	} else {
		if printOutput {
			for _, item := range caveSlice {
				if item.start < 50 {
					// Dont disturb the first 48 bytes of the file that creates the matching hash pattern
					if verbose {
						fmt.Fprintf(os.Stderr, " Cave in first 48 bytes, continuing search \n")
					}
					continue
				}
				if len(hiddenMessage) > (item.end - item.start) {
					// Data does not fit
					fmt.Fprintf(os.Stderr, " Data does not fit \n")
					continue
				}
				if verbose {
					fmt.Fprintf(os.Stderr, " - Data size:          %d\n", len(hiddenMessage))
					fmt.Fprintf(os.Stderr, " - Cave found size:    %d\n", (item.end - item.start)) //(item.end - item.start))
					fmt.Fprintf(os.Stderr, " - Cave found at byte: %d\n", item.start)
				}
				if (item.end - item.start) > len(hiddenMessage) {
					// Insert codecave and break
					modifiedBytes := InsertCode(inputFile, hiddenMessage, item.start)
					if !tryOut {
						WriteFile(outputFile, modifiedBytes)
					}
					break
				}
			}
		}
	}
}

func DecodeMessage(fileLocation string, maximumCodeCaveSize int, byteLengthToHash int, verbose bool) {
	bytePatternToMatchLength := 8
	byteArrayToHash := ReadBytesFromFile(fileLocation, byteLengthToHash) // Creates an array of bytes of max length byteLengthToHash
	bytePatternToMatch := Md5HashBytesReturnXBytes(byteArrayToHash, bytePatternToMatchLength)

	Mess := FindAllMatchingPatterns(fileLocation, bytePatternToMatch, maximumCodeCaveSize)

	if len(Mess.messageBytes) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Found no data\n")
		os.Exit(1)
	}
	unReversedBytes := ReverseBytes(Mess.messageBytes)
	unzippedMessage, err := GunzipData(unReversedBytes)

	if verbose {
		fmt.Fprintf(os.Stderr, " - Found data at byte:      %d\n", Mess.firstStart)
		fmt.Fprintf(os.Stderr, " - Compressed data size:    %d\n", Mess.messageLength)
		fmt.Fprintf(os.Stderr, " - Uncompressed data size:  %d\n", len(unzippedMessage))
		fmt.Fprintf(os.Stderr, " - Message:\n")
	}

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(unzippedMessage))
}

// Gzip an incoming byte array
func GzipData(inputData []byte) (compressedData []byte, err error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	_, err = gz.Write(inputData)
	if err != nil {
		return compressedData, err
	}

	if err = gz.Flush(); err != nil {
		return compressedData, err
	}

	if err = gz.Close(); err != nil {
		return compressedData, err
	}

	compressedData = b.Bytes()
	return compressedData, err
}

// Unzip an incoming byte array
func GunzipData(inputData []byte) (uncompressedData []byte, err error) {
	b := bytes.NewBuffer(inputData)

	var reader io.Reader
	reader, err = gzip.NewReader(b)
	if err != nil {
		return uncompressedData, err
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(reader)
	if err != nil {
		return uncompressedData, err
	}

	uncompressedData = resB.Bytes()

	return uncompressedData, err
}

// Reads a bytearray and returns a reversed version of it
func ReverseBytes(inputArray []byte) []byte {
	revSlice := make([]byte, 0)
	var b byte
	for i := len(inputArray) - 1; i >= 0; i-- {
		b = inputArray[i]
		revSlice = append(revSlice, b)
	}
	return revSlice
}

func CreateHiddenData(patternArrayToMatch []byte, message []byte) []byte {
	hiddenMessage := make([]byte, 0)
	hiddenMessage = append(hiddenMessage, patternArrayToMatch...)
	hiddenMessage = append(hiddenMessage, message...)
	hiddenMessage = append(hiddenMessage, patternArrayToMatch...)
	return hiddenMessage
}
