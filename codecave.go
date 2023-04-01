package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cpuboi/codecave/encoder"
)

var Version = "development"

func getArgs() (string, string, string, bool, bool, bool, bool) {
	var inputFile = flag.String("f", "./example_input_file.bin", "Input file")
	var outputFile = flag.String("o", "./example_output_file.bin", "Output file")
	var doEncodeMessage = flag.String("e", "", "Encode message into file, provide message")
	var doDecode = flag.Bool("d", false, "Decode message from file")
	var verbose = flag.Bool("v", false, "Verbose output")
	var tryOut = flag.Bool("t", false, "Test run, dont actually write output, needs to be combined with -e")
	var analyzeFile = flag.Bool("a", false, "Analyze file, show all codecaves in csv format, combine with -f")

	flag.Parse()
	if *inputFile == "./example_input_file.bin" { // User did not change inputFile variable, print default message
		fmt.Println("\n Codecave - Encodes or decodes hidden messages in files \n The file is loaded into RAM so proceed with care \n")
		flag.PrintDefaults()
		os.Exit(7)
	}
	if !*analyzeFile { // If analyze file is false, check other variables
		if ((len(*doEncodeMessage) > 0) && (*doDecode == true)) || ((*doEncodeMessage == "") && (*doDecode == false)) {
			fmt.Println("Can either encode or decode file, not both at the same time")
			os.Exit(7)
		}
	}
	return *inputFile, *outputFile, *doEncodeMessage, *doDecode, *verbose, *tryOut, *analyzeFile
}

func main() {
	//inputFile, outputFile, inputMessage, doEncode, doDecode, maxMessageSize, verbose := getArgs() // Old
	byteLengthToHash := 48    // How many bytes in the beginning of the file to base the matching pattern hash on
	maxMessageSize := 200_000 // Message is set to max 200kb
	inputFile, outputFile, doEncodeMessage, doDecode, verbose, tryOut, analyzeFile := getArgs()
	if analyzeFile {
		if verbose {
			fmt.Fprintf(os.Stderr, "Analyzing file: %s\n", inputFile)
		}
		minCaveSize := 32 // Cave must be at least 32 bytes
		caveSlice := encoder.FindCaves(inputFile, minCaveSize)
		encoder.CaveSliceParser(inputFile, caveSlice)

	} else if doEncodeMessage != "" {
		if verbose {
			fmt.Fprintf(os.Stderr, "Encoding data in file: %s\n", inputFile)
		}
		encoder.EncodeMessage(inputFile, outputFile, doEncodeMessage, byteLengthToHash, verbose, tryOut)
	} else if doDecode {
		if verbose {
			fmt.Fprintf(os.Stderr, "Decoding data from file: %s\n", inputFile)
		}
		encoder.DecodeMessage(inputFile, maxMessageSize, byteLengthToHash, verbose)
	} else {
		fmt.Println("Ok? ")
	}
}
