package encoder

import (
	"log"
	"os"
)

func WriteFile(filePath string, fileArray []byte) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Write(fileArray)
}

// Read defined amount (length) of bytes from file
func ReadBytesFromFile(filePath string, byteLengthToHash int) []byte {
	byteArray := make([]byte, byteLengthToHash) // Create a byteArray of max size byteLengthToHash
	openFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer openFile.Close()

	openFile.Read(byteArray)
	return byteArray
}

// Open file, load data to an array, at cave.start insert data
func InsertCode(filePath string, hiddenMessage []byte, startByte int) []byte {

	openFile, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}
	defer openFile.Close()
	stat, err := openFile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fileSlice := make([]byte, stat.Size()) // Create a slice of size filesize.
	openFile.Read(fileSlice)

	for _, b := range hiddenMessage {
		fileSlice[startByte] = b
		startByte += 1
	}
	return fileSlice
}
