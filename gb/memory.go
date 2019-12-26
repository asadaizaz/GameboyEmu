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
	internalRAMShadow [7680]byte // 0xE000 -> 0xFDFF (same as internal - 512 bytes) Copy of internal RAM.
	OAM               [160]byte  //Object Attribute Memory (used for sprites)
	ioRegisters       [128]byte  //I/O registers
	zeroPageRAM       [128]byte  // High speed RAM at top of memory
	inBIOS            bool
}

//Reads word from memory
func (mem *Memory) Read(address uint16) byte {
	//TODO: Figure out proper way to read addresses
	switch {
	case address < 0x0100:
		// BIOS
		return mem.bios[address]
	case address < 0x8000:
		//return mem.Cart.Read(address)

	case address < 0xA000:
		//VRAM
		return mem.vRAM[address&(0x9FFF-0x8000)]
	case address < 0xC000:
		//return mem.Cart.Read(address)
	case address < 0xE000:
		return mem.internalRAM[address]
	case address < 0xFE00:
		return mem.internalRAMShadow[address]
	case address < 0xFEA0:
		return mem.OAM[address]
	case address < 0xFF00:
		//Unuseable
		return 0xFF
	case address < 0xFF80:
		return mem.ioRegisters[address]

	default:
		return mem.ReadHighMem(address)
	}

}
