package main


import (
	"chip8/chippy"
	"os"
	"strconv"
)




// Simple assembler for the .chippy assembly language, also kinda hacky :(
// see the chippy directory for more information


func main() {
	sourceFile := os.Args[1]
	requestedOutputFile := os.Args[2]
	printSymbolTable, err := strconv.ParseBool(os.Args[3])
	if err != nil {
		printSymbolTable = false
	}

	chippy.ParseAndCompile(sourceFile, requestedOutputFile,  printSymbolTable)
}