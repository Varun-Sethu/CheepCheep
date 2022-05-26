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
2 bit integer following the opcode, below is a table of the various addressing modes and their syntax within the assembler  

|  Addressing Mode  | Syntax      | Encoding |
|       ---         |    ---      |   ---    |
| Immediate/Default |    #x       | ```00``` | 
|    Register       |    $rx      | ```01``` |
| Register Relative |    x+$rx    | ```10``` |
|    PC Relative    |   #(x)      | ```11``` |

#### Instructions
Instructs and operands consist of 4 bytes. The first byte is the instruction and the following 3 bytes are the operands. For the instruction byte the first 2 bits is the addressing mode and the final 6 are the the opcode. As an example conider the add instruction with immediate and register addressing.
`add $t1 $t2 $t3` will map to `01 <add op code>` while `add $t1 $t2 #4` maps to `00 <add op code>`.

#### Memory Layout
There are $2^{16}$ unique addresses on this machine hence to have an address thats an argument we require 2 bytes.

#### Sample Code
```x86
add $r1, $r2, $r3
add $r1, $r2, #4

.label
    jmp .label
```
This architecture has 13 registers, 1 register for the stack pointer, 1 register for the output of compare instructions and a zero register that contains the number 0. All of this are integer registers and the machine doesn't support floating point operations.

In code the first 13 registers are addressed with `$r[1 -- 13]` while the stack/cmp/zero register are addressed as `$sp, $cmp, $zero`.