package chippy

import "regexp"

// The OPCODES table is broken up as follows:
/**
	Index 0: The bytecode for the instruction
	Index 1: The required amount of instructions
 */
var OPCODES = map[string][]uint16 {
	"LDR": {0x1, 2, REGDIRECT, IMMEDIATE | DIRECT | INDIRECT | REGDIRECT | REGINDIRECT | RELATIVE | INDEXED},
	"STR": {0x2, 2, REGDIRECT, DIRECT | INDIRECT | REGDIRECT | REGINDIRECT | RELATIVE | INDEXED},
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
const RELATIVE uint16 = 0x20
const INDEXED uint16 = 0x40

// TODO: For the register values add the restriction on the addressable registers
// TODO: maybe rethink the way addressing modes work?
// TODO: support nested addressing modes: eg: (r1) +3
var ImmediateR = regexp.MustCompile(`^#(\d+)$`)
var DirectR = regexp.MustCompile(`^([\d]+|\.[a-zA-z]+)$`)
var IndirectR = regexp.MustCompile(`^\(([\d]+|\.[a-zA-z]+)\)$`)
var RegDirectR = regexp.MustCompile(`^r([1-9]|1[0-5])$`)
var RegIndirectR = regexp.MustCompile(`^\(r([1-9]|1[0-5])\)$`)

var RelativeR = regexp.MustCompile(`^\(r([1-9]|1[0-5])\)\+(\d+)$`)
var IndexedR = regexp.MustCompile(`^([\d]+|\.[a-zA-z]+)\+(\d+)$`)