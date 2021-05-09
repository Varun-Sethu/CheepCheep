# Chippy

Chippy is designed to just be a small and experimental assembler for a random ISA loosely based on the 6502 and Chip8, it exists purely for experimental reasons :).

## Opcodes
| Opcode | Params |
|  ---   |   ---  |
|  JMP   |   Address (See addressing modes) |


#### Addressing modes
Memory locations within the system can be addressed in several ways, in bytecode the addressing mode is specified by a
3 bit integer following the opcode, below is a table of the various addressing modes and their syntax within the assembler  

|  Addressing Mode  | Syntax   | Encoding |
|       ---         |    ---   |    ---   |
|    Immediate      |   #x     |          | 
|    Direct         |   x      |          |
|    Indirect       |  (x)     |          |
| Register Direct   |   rx     |          |
| Register Indirect |  (rx)    |          |
| Indexed Register  | (rx)+nnn |          |
|   Relative        |  x+nnn   |          |

There are 7 independent states requiring 3 bits of information to encode.
