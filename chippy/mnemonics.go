package chippy

import "regexp"

// The OPCODES table is broken up as follows:
/**
	Index 0: The bytecode for the instruction
	Index 1: The required amount of instructions
 */
var OPCODES = map[string][]uint16 {

	// print contents of register to standard out
	"PRINT": {0xC, 1, REGDIRECT},
	"HLT": {0x0, 0, 0},

	// load value into register
	"LDR": {0x1, 2, REGDIRECT, IMMEDIATE | DIRECT | INDIRECT | REGDIRECT | REGINDIRECT | INDEXED | INDEXSCALED},
	// store value in register into memory
	"STR": {0x2, 2, REGDIRECT, DIRECT | INDIRECT | REGDIRECT | REGINDIRECT | INDEXED | INDEXSCALED},

	// add to the value stored in a register
	"ADD": {0x3, 2, REGDIRECT, IMMEDIATE | DIRECT | INDIRECT | REGDIRECT | REGINDIRECT | INDEXED | INDEXSCALED},

	// compare the two values, the one stored in r1 and the second operand
	"CMP": {0x4, 2, REGDIRECT, IMMEDIATE | DIRECT | INDIRECT | REGDIRECT | REGINDIRECT | INDEXED | INDEXSCALED},

	// jump instructions based on conditional flags
	"JMPL": {0x5, 1, IMMEDIATE},
	"JMPG": {0x6, 1, DIRECT},
	"JMPLE": {0x7, 1, DIRECT},
	"JMPGE": {0x8, 1, DIRECT},
	"JMPE": {0x9, 1, DIRECT},
	"JMPNE": {0xA, 1, DIRECT},
	// Normal jumps
	"JMP": {0xB, 0, 0},


}
const OPTBYTE int = 0
const OPTARG  int = 1

/*
	Notes on compiled bytecodes:
		- Most operations are operations on registers
		- The addressing mode is stored along with the opcode
		- Instructions are broken down as follows: [opcode]<-5bits + [addressing mode]<-3bits + ... arguments
 */


// NOTE:
/*
	- Each opcode has a set of supported addressing modes for each of its potential arguments, this set of supported
	- modes is represented as a bit mask for example: 0011001
	- to query if an addressing mode is supported we simply apply it mask bitmask & code, the ith bit represents
      if the opcode supports addressing mode i.
	- Each argument has its own supported set of addressing modes: eg. the first argument of load only supports registers
	  while the second argument can be whatever we want.

 */



// The following set of ADDRMODE values provides a set of regular expressions for determining the type of addressing
// mode an argument to an operation follows
/**


*/
// Encodings for the various addressing modes
const IMMEDIATE uint16 = 0x1
const DIRECT uint16 = 0x2
const INDIRECT uint16 = 0x4
const REGDIRECT uint16 = 0x8
const REGINDIRECT uint16 = 0x10
const INDEXED uint16 = 0x20
const INDEXSCALED uint16 = 0x40

// TODO: For the register values add the restriction on the addressable registers
// TODO: maybe rethink the way addressing modes work?
// TODO: support nested addressing modes: eg: (r1) +3
// immediate values can be a value or a memory location (for jump and store)
var ImmediateR = regexp.MustCompile(`^(#[\d]+|\.[a-zA-z]+)$`)
var DirectR = regexp.MustCompile(`^([\d]+)$`)
var IndirectR = regexp.MustCompile(`^\(([\d]+|\.[a-zA-z]+)\)$`)
var RegDirectR = regexp.MustCompile(`^r([1-9]|1[0-5])$`)
var RegIndirectR = regexp.MustCompile(`^\(r([1-9]|1[0-5])\)$`)

var IndexedR = regexp.MustCompile(`^([\d]+|\.[a-zA-z]+)\+r(\d+)$`)
var IndexScaledR = regexp.MustCompile(`^([\d]+|\.[a-zA-z]+)\+([\d]+)\*r(\d+)$`)