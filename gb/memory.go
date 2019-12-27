package gb

/*
   General Memory Map
   ==================
   When booting = 0x00, override with:
   -MMU:
   0x0000-0x00FF   Boot ROM (256 bytes) [DONE]
   When booting != 0x00:
   -Cartridge:
   0x0000-0x3FFF   16KB ROM Bank 00     (in cartridge, fixed at bank 00)
   0x4000-0x7FFF   16KB ROM Bank 01..NN (in cartridge, switchable bank number)
   -GPU:
   0x8000-0x9FFF   8KB Video RAM (VRAM) (switchable bank 0-1 in CGB Mode)
   -Cartridge:
   0xA000-0xBFFF   8KB External RAM     (in cartridge, switchable bank, if any)
   -MMU:
   0xC000-0xCFFF   4KB Work RAM Bank 0 (WRAM)
   0xD000-0xDFFF   4KB Work RAM Bank 1 (WRAM)  (switchable bank 1-7 in CGB Mode)
   0xE000-0xFDFF   Same as C000-DDFF (ECHO)    (typically not used)
   -GPU:
   0xFE00-0xFE9F   Sprite Attribute Table (OAM)
   -N/A:
   0xFEA0-0xFEFF   Not Usable
   -GamePad/SerialData/Timer/Audio/GPU:
   0xFF00-0xFF7F   I/O Ports
   -MMU:
   0xFF80-0xFFFE   High RAM (HRAM)
   0xFFFF          Interrupt Enable Register
*/

type Memory struct {
	//bios 256[byte]
	//cartridge
	HighRAM [0x100]byte  //Used for quick ram acccess. Also contains i/o registers. (256 bytes)
	VRAM    [0x4000]byte // 8kb vram for norma. Switchable bank (0-1) for CGB
	VRAMIdx byte         // For CGB support
	WRAM    [0x9000]byte // 8kb (two banks). CGB has 1-7 banks; 0xC000 - 0xDFFF
	WRAMIdx byte         //	For CGB support.
	OAM     [0x100]byte  // Object attribute memory.

}

func (mem *Memory) Init(gameboy *GameBoy) {

	// Set the default values for PowerUp sequence
	mem.HighRAM[0x04] = 0x1E
	mem.HighRAM[0x05] = 0x00
	mem.HighRAM[0x06] = 0x00
	mem.HighRAM[0x07] = 0xF8
	mem.HighRAM[0x0F] = 0xE1
	mem.HighRAM[0x10] = 0x80
	mem.HighRAM[0x11] = 0xBF
	mem.HighRAM[0x12] = 0xF3
	mem.HighRAM[0x14] = 0xBF
	mem.HighRAM[0x16] = 0x3F
	mem.HighRAM[0x17] = 0x00
	mem.HighRAM[0x19] = 0xBF
	mem.HighRAM[0x1A] = 0x7F
	mem.HighRAM[0x1B] = 0xFF
	mem.HighRAM[0x1C] = 0x9F
	mem.HighRAM[0x1E] = 0xBF
	mem.HighRAM[0x20] = 0xFF
	mem.HighRAM[0x21] = 0x00
	mem.HighRAM[0x22] = 0x00
	mem.HighRAM[0x23] = 0xBF
	mem.HighRAM[0x24] = 0x77
	mem.HighRAM[0x25] = 0xF3
	mem.HighRAM[0x26] = 0xF1
	mem.HighRAM[0x40] = 0x91
	mem.HighRAM[0x41] = 0x85
	mem.HighRAM[0x42] = 0x00
	mem.HighRAM[0x43] = 0x00
	mem.HighRAM[0x45] = 0x00
	mem.HighRAM[0x47] = 0xFC
	mem.HighRAM[0x48] = 0xFF
	mem.HighRAM[0x49] = 0xFF
	mem.HighRAM[0x4A] = 0x00
	mem.HighRAM[0x4B] = 0x00
	mem.HighRAM[0xFF] = 0x00

	mem.WRAMIdx = 1
}

//Reads word from memory
func (mem *Memory) Read(address uint16) byte {
	switch {
	case address < 0x8000:
		// BIOS
		//return mem.Card.Read(address)
	case address < 0xA000:
		//VRAM
		bankOffset := uint16(mem.VRAMIdx) * 0x2000
		return mem.VRAM[address-0x8000+bankOffset]
	case address < 0xC000:
		//return mem.Cart.Read(address)
	case address < 0xD000:
		return mem.WRAM[address-0xC000]
	case address < 0xE000:
		return mem.WRAM[(address-0xC000)+(uint16(mem.WRAMIdx)*0x1000)]
	case address < 0xFE00:
		//Shadow
		return 0xFF
	case address < 0xFEA0:
		return mem.OAM[address-0xFE00]
	case address < 0xFF00:
		//Not used
		return 0xFF

	default:
		return mem.ReadHighMem(address)
	}
}

// Write a value at an address to the relevant location based on the
// current state of the gameboy. This handles banking and side effects
// of writing to certain addresses.
func (mem *Memory) Write(address uint16, value byte) {
	switch {
	case address < 0x8000:
		// Write to the cartridge ROM (banking)
		mem.Cart.WriteROM(address, value)

	case address < 0xA000:
		// VRAM Banking
		bankOffset := uint16(mem.VRAMBank) * 0x2000
		mem.VRAM[address-0x8000+bankOffset] = value

	case address < 0xC000:
		// Cartridge ram
		mem.Cart.WriteRAM(address, value)

	case address < 0xD000:
		// Internal RAM - Bank 0
		mem.WRAM[address-0xC000] = value

	case address < 0xE000:
		// Internal RAM Bank 1-7
		mem.WRAM[(address-0xC000)+(uint16(mem.WRAMBank)*0x1000)] = value

	case address < 0xFE00:
		// Echo RAM
		// TODO: re-enable echo RAM?
		//mem.Data[address] = value
		//mem.Write(address-0x2000, value)

	case address < 0xFEA0:
		// Object Attribute Memory
		mem.OAM[address-0xFE00] = value

	case address < 0xFF00:
		// Unusable memory
		break

	default:
		// High RAM
		mem.WriteHighRam(address, value)
	}
}
