package chippy

import (
	"fmt"
	"os"
)



// Chippy is just a small and hacky assembler
/*
	There are 3 exposed method within the module:
		- ParseAndCompile takes a specified fileName, opens it and saves a compiled bytecode to "saveAs"
		- ParseFile takes a fileName and the contents of the file and parses it, if there is any issues and error will
 		  be returned
		- CompileTokens takes a fileName, io.Writer and a series of parsed tokens, it compiles the tokens and writes them
		  to the io.Writer, if  there are any issues and error will be thrown


 */




// ParseAndCompile takes a source chippy file, parses it and generates the corresponding bytecode
func ParseAndCompile(fileName string, saveAs string, printSymbolTable bool) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0777)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	fileSize, _ := file.Stat()
	buffer := make([]byte, fileSize.Size())
	file.Read(buffer)

	tokens, parseError := ParseFile(fileName, string(buffer))
	if parseError != nil {
		fmt.Print(parseError.Error())
		return
	}

	// Create a file called "saveAs"
	output, _ := os.Create(fmt.Sprintf("%s.chip", saveAs))
	defer output.Close()
	symtable, compilationErr := CompileTokens(fileName, output, tokens)
	if compilationErr != nil {
		fmt.Print(compilationErr.Error())
		return
	}

	if printSymbolTable {
		fmt.Printf("The computed symtable was: %v\n", symtable)
	}

}























