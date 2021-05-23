package emulator

// just some of the basic opcodes
const HLT uint8 = 0x0
const PRINT uint8 = 0x3
const LDR uint8 = 0x1
const STR uint8 = 0x2

const ALU uint8 = 0x8
	const ADD uint8 = 0x0
	const SUB uint8 = 0x1
	const MUL uint8 = 0x2
	const DIV uint8 = 0x3
	const XOR uint8 = 0x4
	const AND uint8 = 0x5
	const OR  uint8 = 0x6
	const NOT uint8 = 0x7

// comparison and jump operations
const CMP uint8 = 0x4
const JMPL uint8 = 0x5
const JMPG uint8 = 0x6
const JMP uint8 = 0x7
const JMPLE uint8 = 0x10
const JMPGE uint8 = 0x11

// addressing modes
const immediate uint8 = 0
const direct uint8 = 1
const indirect uint8 = 2
const registerDirect uint8 = 3
