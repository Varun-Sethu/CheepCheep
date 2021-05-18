# Chippy

Chippy is just be a small and experimental assembler for a random ISA loosely based on the 6502 and Chip8, it exists purely for experimental reasons :).
Outline below are some opcode mnemonics for the assembler as well a description of their function. 
## Opcodes
| Opcode | Params | Description |
|  ---   |   ---  |     ---     |
|  LDR   |  rx, n | Loads n into the register rx, note the value of x is dictated by its addressing mode | 
|  STR   |  rx, n | Stores the value in rx to memory location x |
|  JMP   |    n   | Jump to memory location n |

TODO: implement 

#### Addressing modes
Memory locations within the system can be addressed in several ways, in bytecode the addressing mode is specified by a
3 bit integer following the opcode, below is a table of the various addressing modes and their syntax within the assembler  

|  Addressing Mode  | Syntax      | Encoding |
|       ---         |    ---      |    ---   |
|    Immediate      |   #x        |```000``` | 
|    Direct         |   x         |```001``` |
|    Indirect       |  (x)        |```010``` |
| Register Direct   |   rx        |```011``` |
| Register Indirect |  (rx)       |```100``` |
| Indexed Register  |  x+rx       |```101``` |
| Indexed Scaled    |  x+k*rx     |```110``` |

There are 7 independent states requiring 3 bits of information to encode.

#### Instructions
Instructions have a relatively simple encoding in memory, the first 5 bits refers to the opcode, the following 3 refers to the addressing mode and every byte to follow are
the arguments for the instruction. For example the instruction: ```ldr r1 r15``` is encoded as: ```0x0b010f```
in memory, in binary this is: ```0b 00001011 00000001 00001111```.