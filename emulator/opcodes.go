package emulator



// just some of the basic opcodes
const HLT 	uint8 = 0x0
const PRINT uint8 = 0xC
const LDR uint8 = 0x1
const STR uint8 = 0x2

const ADD uint8 = 0x3
const CMP uint8 = 0x4
const JMPL uint8 = 0x5



// addressing modes
const immediate uint8 = 0
const direct uint8 = 1
const indirect uint8 = 2
const registerDirect uint8 = 3
