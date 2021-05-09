# CheepCheep

Just a reasonably small emulator + assembler for a custom ISA, not entirely a serious project. The ISA is loosely based on the 6502 and the Chip8 system.
It's a 8bit system with a 16 bit address bus, see "chippy" for more details.

### Chippy
The assembler is called "chippy", usage is rather simple; just build the source file and call the compiled source
```shell script
go build chippy.go
./chippy sourcefile.chippy outfile

hexdump -C outfile.chip
```
Once again more information regarding the assembler can be found within "chippy", it also contains a cursory glance at the instruction set of the ISA as well as supported addressing modes.