package gb

/*
	Represents memory for z80 cpu.
*/

type Memory struct {
	//bios 256[byte]
	//cartridge
	bios              [256]byte  // 0x0000 -> 0x00FF
	vRAM              [8192]byte // 0x8000 -> 9FFF
	internalRAM       [8192]byte // 0xC000 -> 0xDFFF
	internalRAMShadow [7680]byte // 0xE000 -> 0xFDFF (same as internal - 512 bytes)
	OAM               [160]byte  //Object Attribute Memory (used for sprites)
	ioRegisters       [128]byte  //I/O registers
	zeroPageRAM       [128]byte  // High speed RAM at top of memory
}
