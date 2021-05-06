package emulator



// defines the opcodes as constants
// "System operations"
const SYS	  uint8 = 0x0
  const RET   uint8 = 0xee // return from a subroutine
  const CLR	  uint8 = 0xe0  // clear the screen

// http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#00E0

// Control flow operations
const JMP 	  uint8 = 0x10 // jump
const CALL	  uint8 = 0x20 // calls a sub routine
const SE      uint8 = 0x30 // skip next instruction if Vx = kk (3xKK)
const SNEK    uint8 = 0x40 // skip next instruction if Vx != kk
const SER     uint8 = 0x50 // skip next instruction if Vx == Vy
const SNE     uint8 = 0x90 // skip next instruction if Vx != Vy
const JMPN    uint8 = 0xB0 // jump to location nnn + V0
// Skip instructions
const SKP    uint8 = 0xE
  const SKPE  uint8 = 0x9E // Skip the next instruction if the inputted key = the value at Vx
  const SKPNE uint8 = 0xA1 // skip the next instruction if the inputted key != the value at Vx


// State management functions (Loading and in place adding)
const LD 	  uint8 = 0x60 // set the register value
const LDI     uint8 = 0xA0 // set the index register
const DISP	  uint8 = 0xD0 // display to the screen
// Loading into chip8 state
const LDC     uint8 = 0xF0 // load into important registers within the chip8
  const LFDT   uint8 = 0x7 // load contents from  delay timer into VX
  const PSE    uint8 = 0xA // wait for a key press and load into Vx
  const LDT    uint8 = 0x15 // load the value of Vx into the DT
  const LDS    uint8 = 0x18 // load the value of Vx into the sound timer
  const ADDI   uint8 = 0x1E // set IR += Vx
  


// Arithmetic operations
const ADDK    uint8 = 0x70 // add KK to Register Vx: 7xxk
const ALU     uint8 = 0x80 // Indicates a regular arithmetic operation
  const LDR   uint8 = 0x0 // load Vy into Vx: 8xy0
  const OR    uint8 = 0x1
  const AND   uint8 = 0x2
  const XOR   uint8 = 0x3
  const ADD   uint8 = 0x4
  const SUB   uint8 = 0x5 // vx = vx - vy
  const SHR   uint8 = 0x6 // performs a bitwise shift right
  const SUBN  uint8 = 0x7 // vx = vy - vx
  const SHL   uint8 = 0xE // performs a bitwise shift left
const RND     uint8 = 0xC0 // Generate a random number between 0 and 255 which is ANDed against kk


