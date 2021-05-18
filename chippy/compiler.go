package chippy

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

/*
	compiler.go
		- The purpose of the compiler function is to read the token array generated by the parser
		  and produce a compiled byte array

		- The compiler performs no real syntax validation
		- The compiler additionally generates a symbol table for all the labels declared within the assembly program
		- The compiler may however thrown an error if a requested label is not defined

		- Everything is in big endian mode :)
 */


// ROMSTART defines what address in memory all compiled ROMs start at
const ROMSTART uint16 = 0x0


// CompileTokens takes a series of tokens and produces a compiled array of bytes representing the machine code
// for the source file, it then writes this compiled series of bytes into the requested file, its a separate function
// for debugging reasons :)
func CompileTokens(fileName string, fileHandle io.Writer ,tokens []token) (map[string]uint16, error) {
	bytecode, symbolTable, compilationFailure  := compileInstructions(tokens)
	if compilationFailure != nil {
		return nil, fmt.Errorf("Compilation Error:\n%s:%s\n", fileName, compilationFailure.Error())
	}

	fileHandle.Write(bytecode)
	return symbolTable, nil
}




// validateAddressingMode takes a token representing an opcode
func validateAddressingMode(operation token) (bool, uint8) {

	// If the operation requires less than 1 argument then just give up and try to return an addressing mode
	if len(operation.tokenParams) == 0 {
		return true, 0
	} else if len(operation.tokenParams) == 1 {
		return true, uint8(math.Log2(float64(operation.tokenParams[0].addressingMode)))
	}

	// Iterate over all parameters within the token and verify correctness
	var expectedAddressingMode = operation.tokenParams[1].addressingMode
	for _, parameter := range operation.tokenParams[2:] {
		// If the addressing mode is invalid just give up
		if parameter.addressingMode != expectedAddressingMode {
			return false, 0
		}
	}

	// hmmm... yes
	return true, uint8(math.Log2(float64(expectedAddressingMode)))
}




// resolveAddress takes a requested address in string format and determines its physical location
// note: this requires the computed symbol table
func resolveAddress(addr string, symbolTable map[string]uint16) uint16 {
	representation, err := strconv.Atoi(addr)

	// an error is only thrown if the value is an immediate value;
	// there are two cases: the value is #nnn or .label, in the first case we just call Atoi on addr[1:]
	// in the second case we search the value in the lookup table
	if err != nil {
		if addr[0] == '#' {
			// integer
			representation, _ = strconv.Atoi(addr[1:])
		} else {
			// label
			representation = int(symbolTable[addr[1:]])}
	}
	return uint16(representation)
}





// computeRepresentation takes an argument to an opcode and determines how to represent it within memory
// Note: function is a bit of a mess, sorry :(
func computeRepresentation(arg param, symbolTable map[string]uint16) []uint8 {
	// Each argument can be broken down as follows:
	//	baseLocation: exists for all arguments
	//	offset:	exists for relative and indexed arguments
	var baseLocation = resolveAddress(arg.addressData[0], symbolTable)
	// Based off the computed representation figure out how to pack it into memory
	// Note: Registers require 1 byte of memory (16 unique registers -> 4 bits -> packed to 8 bits)
	//		 Actual values require 1 byte as this is an 8 bit machine, with the exception of addresses
	//		 Memory locations require 2 bytes as 64kb of memory are addressable
	var computedMemoryBlock []uint8
	var isImmediateValue  bool = arg.addressingMode == IMMEDIATE && arg.addressData[0][0] == '#'

	// register and immediate value commands
	if arg.addressingMode == REGDIRECT || arg.addressingMode == REGINDIRECT || isImmediateValue {
		computedMemoryBlock = append(computedMemoryBlock, uint8(baseLocation))
	// regular blocks of memory
	} else {
		computedMemoryBlock = append(computedMemoryBlock, uint8((baseLocation & 0xff00) >> 8))
		computedMemoryBlock = append(computedMemoryBlock, uint8(baseLocation & 0xff))
	}

	// if the offset is defined then append that too (note offsets are 1 byte long as they are registers)
	if len(arg.addressData) >= 2 {
		// Now there are two situations where the arg.addressData length would be greater than or equal to 2
		// the command is either at: indexed register or an indexed scaled addressing mode

		baseRegister, _ := strconv.Atoi(arg.addressData[1])
		computedMemoryBlock = append(computedMemoryBlock, uint8(baseRegister))
		// Special case for scaled indexes as they have 3 arguments
		if arg.addressingMode == INDEXSCALED {
			// Just append on the register as previous we tacked on the constant value
			offsetRegister, _ := strconv.Atoi(arg.addressData[2])
			computedMemoryBlock = append(computedMemoryBlock, uint8(offsetRegister))
		}
	}

	return computedMemoryBlock
}




// compileInstructions takes a symbol table and a series of tokens and produces a compiled byte array corresponding
// to the bytecode of the assembled file, the function throws an error if there were any issues during compilation
func compileInstructions(tokens []token) ([]byte, map[string]uint16, error) {

	var currentMemoryOffset uint16 = 0
	var compiledByteCode []byte
	var symbolTable  = map[string]uint16{}

	// Iterate over each token in our token set and only take action if it is an OPERATION token
	for _, token := range tokens {
		if token.tokenType == OPERATION {
			// The token is an operation, there are a few key stages when processing an operation
			/*
				Step 1: Determine the opcode for the operation (lookup in the opcode table)
				Step 2: Determine the addressing mode for the command (this is the addressing mode of the parameters)
						that follow after the first parameter
				Step 3: Resolve the arguments of the operation
							- If it is a label, find the correct address
							- Otherwise, determine how we are going to pack he data into the correct nibble format
			 */

			// keeps track of the size of the instruction within the binary source file
			var instructionSize uint16 = 1

			// Steps 1 and 2
			validMode, requestedMode := validateAddressingMode(token)

			if !validMode {
				return nil, nil, fmt.Errorf("%d: inconsistent addressing modes for \"%s\"",
					token.lineNumber, token.tokenData)
			}
			// finally compute the final instruction (opcode + addressing mode included)
			instructionOpcode := OPCODES[strings.ToUpper(token.tokenData)][OPTBYTE]
			var fullInstruction uint8 = uint8(instructionOpcode << 3) | requestedMode

			// Step 3: compute and resolve the arguments by iterating, appending them to a finally computed array
			var argumentArray []uint8
			for _, arg := range token.tokenParams {
				rep := computeRepresentation(arg, symbolTable)
				instructionSize += uint16(len(rep))
				argumentArray = append(argumentArray, rep...)
			}


			// Append everything to the compiled bytecode and increase the current memory offset :)
			compiledByteCode = append(compiledByteCode, fullInstruction)
			compiledByteCode = append(compiledByteCode, argumentArray...)
			currentMemoryOffset += instructionSize


		} else if token.tokenType == LABEL {
			// to resolve labels we simply just need to update its location in our computed symbol table
			if _, ok := symbolTable[token.tokenData]; !ok {
				// We dont increment the currentMemoryOffset as the label points to the NEXT instruction
				symbolTable[token.tokenData] = ROMSTART + currentMemoryOffset
			} else {
				return nil, nil, fmt.Errorf("%d: label declaration \"%s\" shadows a previous declaration",
					token.lineNumber, token.tokenData)
			}

		}
	}

	return compiledByteCode, symbolTable, nil

}

