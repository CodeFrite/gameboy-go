package gameboy

/*
 * Print the current instruction being executed in the following format (examples):
 * PC: 0x00A0, Bytes: 00 			, ASM: NOP
 * PC: 0x00A1, Bytes: 40 			, ASM: LD B, $40
 * PC: 0x00A2, Bytes: 3E 01 	, ASM: LD A, $01
 * PC: 0x00A4, Bytes: F8 4E 	, ASM: LD HL, SP + $4E
 * PC: 0x00A6, Bytes: EA AB 01, ASM: LD [$01AB], HL
 */
