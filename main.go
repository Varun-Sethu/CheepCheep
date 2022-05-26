package main

import (
	"bufio"
	"cheepcheep/chippy"
	"os"
)

// Simple assembler for the .chippy assembly language, also kinda hacky :(
// see the chippy directory for more information

func main() {
	sourceFile := os.Args[1]
	outputFile := os.Args[2]

	f, _ := os.Open(sourceFile)
	o, _ := os.Create(outputFile)
	defer f.Close()
	defer o.Close()

	nodes := chippy.Parse(*bufio.NewReader(f))
	compiledStream := chippy.Compile(nodes)

	// write the compiled stream out
	w := bufio.NewWriter(o)
	for c, err := compiledStream.ReadByte(); err == nil; c, err = compiledStream.ReadByte() {
		w.WriteByte(c)
	}
	w.Flush()
}
