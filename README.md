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
Once the ROMs have been assembled into bytecode they can be run on the emulator by simply calling
```shell script
./emulator.out binaries/rom.chip
```
To clean the current directory run
```shell script
make clean
```

## Details
The assembler's name is "Chippy" :). I am yet to implement any pseudo-ops and the ability to write non operation data
directly into the assembled file (eg. global variables). At the emulator level it would be nice to support some sort of segmentation.