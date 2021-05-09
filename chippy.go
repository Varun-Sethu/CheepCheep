package main


import (
	"chip8/chippy"
	"os"
)




// Simple assembler for the .chippy assembly language, also kinda hacky :(
// see the chippy directory for more information


func main() {
	sourceFile := os.Args[1]
	requestedOutputFile := os.Args[2]
	chippy.ParseAndCompile(sourceFile, requestedOutputFile)
}