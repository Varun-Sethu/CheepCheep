# CheepCheep

Just a reasonably small emulator + assembler for a custom ISA, not entirely a serious project. The ISA is loosely based on the 6502 and the Chip8 system.
It's a 8bit system with a 16 bit address bus, see "chippy" for more details.

### Compilation and setup
A makefile is provided to ease development, to compile all the ROMs within the ROMs/ directory run:
```shell script
make roms
```
Likewise, to build the emulator and assembler target run
```shell script
make all
```
To clean the current directory run
```shell script
make clean
```