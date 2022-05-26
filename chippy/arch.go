package chippy

// The OPCODES table is broken up as follows:
/**
Index 0: The bytecode for the instruction
Index 1: The number of arguments
Index 2: The type of the second argument
Index 3: The type of the third argument
...
*/
const (
	BYTECODE = iota
	NUMARGS
)

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

var OPCODES = map[string][]uint16{
	// halt
	"HLT": {0x0, 0},

	// move value into register, either from a register or a direct value
	"MOV": {0x1, 2, RegisterValue, Addr | RegisterValue | RegisterRelativeValue},

	// load value from memory into a register
	"LDR": {0x2, 2, RegisterValue, ImmediateValue | RegisterRelativeValue},

	// print contents of register to standard out
	"PRINT": {0x3, 1, RegisterValue},

	// add to the value stored in a register
	"ADD": {0x8, 2, RegisterValue, ImmediateValue | RegisterValue},
	"SUB": {0x9, 2, RegisterValue, ImmediateValue | RegisterValue},
	"MUL": {0xA, 2, RegisterValue, ImmediateValue | RegisterValue},
	"DIV": {0xB, 2, RegisterValue, ImmediateValue | RegisterValue},
	"XOR": {0xC, 2, RegisterValue, ImmediateValue | RegisterValue},
	"AND": {0xD, 2, RegisterValue, ImmediateValue | RegisterValue},
	"OR":  {0xE, 2, RegisterValue, ImmediateValue | RegisterValue},
	"NOT": {0xF, 2, RegisterValue, ImmediateValue | RegisterValue},

	// compare the two values, the one stored in r1 and the second operand
	"CMP": {0x4, 2, RegisterValue, ImmediateValue | RegisterValue},

	// jump instructions based on conditional flags
	"JMPL":  {0x5, 1, Addr | Label | RegisterRelativeValue | PCRelativeValue},
	"JMPG":  {0x6, 1, Addr | Label | RegisterRelativeValue | PCRelativeValue},
	"JMP":   {0x7, 1, Addr | Label | RegisterRelativeValue | PCRelativeValue},
	"JMPLE": {0x10, 1, Addr | Label | RegisterRelativeValue | PCRelativeValue},
	"JMPGE": {0x11, 1, Addr | Label | RegisterRelativeValue | PCRelativeValue},
}

// The register table maps register value to their appropriate numeric value on the arch
var REGISTERS = map[string]uint8{
	"r1":   1,
	"r2":   2,
	"r3":   3,
	"r4":   4,
	"r5":   5,
	"r6":   6,
	"r7":   7,
	"r8":   8,
	"r9":   9,
	"r10":  10,
	"r11":  11,
	"r12":  12,
	"r14":  13,
	"sp":   14,
	"cmp":  15,
	"zero": 16,
}

// ADDRMODES maps nodeTypes to integers representing
// the addr mode
var ADDRMODES = map[nodeType]uint8{
	ImmediateValue:        0,
	Label:                 0,
	Addr:                  0,
	RegisterValue:         1,
	RegisterRelativeValue: 2,
	PCRelativeValue:       3,
}
