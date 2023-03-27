package encoder

import (
	"bufio"
	"log"
	"os"
)

type caveStruct struct {
	start int
	end   int
}

// Look for consecutive zeros
func FindCaves(filePath string, minCaveSize int) []caveStruct {

	openFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer openFile.Close()
	//openFile.Seek(int64(50), 0) // Skip 50 bytes since 48 bytes in the beginning of the files create the matching pattern hash

	scanner := bufio.NewScanner(openFile)
	scanner.Split(bufio.ScanBytes)

	count := 0
	cStruct := caveStruct{start: 0, end: 0}
	caveSlice := make([]caveStruct, 0)
	start_match := 0
	end_match := 0
	var inCave bool
	for scanner.Scan() {
		count += 1
		b := scanner.Bytes()
		if b[0] == 0 {
			// Since last byte was not 0, assume this is the beginning of the code cave
			if !inCave {
				start_match = count
			}
			inCave = true
		} else {
			// If the character is not 0 but inCave is true, then the last byte was the last 0.
			if inCave {
				end_match = count - 1
				// Add the values to caveSlice
				if end_match-start_match >= minCaveSize {
					cStruct.start = start_match
					cStruct.end = end_match
					caveSlice = append(caveSlice, cStruct)
				}
			}
			inCave = false
		}
	}
	return caveSlice
}
