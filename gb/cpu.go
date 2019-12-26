package gb

import (
	"github.com/asadaizaz/gameboyemu/bits"
)

type register struct {
	value uint16
	mask  uint16 //a mask over the register, to get lower bits. Only for AF, because lower bits of F are always 0

}

func (reg *register) Hi() byte {
	return byte(reg.value >> 8)
}
func (reg *register) Lo() byte {
	return byte(reg.value & 0xFF)
}

// HiLo gets the 2 byte value of the register.
func (reg *register) HiLo() uint16 {
	return reg.value
}
func (reg *register) SetHi(val byte) {
	reg.value = uint16(val)<<8 | (uint16(reg.value) & 0xFF)
	reg.updateMask()
}

// SetLog sets the lower byte of the register.
func (reg *register) SetLo(val byte) {
	reg.value = uint16(val) | (uint16(reg.value) & 0xFF00)
	reg.updateMask()
}

// Set the value of the register.
func (reg *register) Set(val uint16) {
	reg.value = val
	reg.updateMask()
}

// Mask the value if one is set on this register.
func (reg *register) updateMask() {
	if reg.mask != 0 {
		reg.value &= reg.mask
	}
}

// CPU is definition of Z80 CPU
type CPU struct {
	//Registers
	AF register //Accumulators and flags
	BC register
	DE register
	HL register

	PC uint16   //Program counter
	SP register //stack pointer
	//Clock
	// memory?
	// Instructions

}

/*
Flags:
	Z - Zero (Set if value is zero)
	H - Half carry (carry from bit 3 to bit 4)
	N - Subtract (Set if last operation was sub)
	C - Carry	(Set if result did not fit in register)
*/
func (cpu *CPU) Init() {
	cpu.PC = 0x100 //Might need to change to 0x0?
	cpu.AF.Set(0x0000)
	cpu.BC.Set(0x0000)
	cpu.DE.Set(0x0000)
	cpu.HL.Set(0x0000)
	cpu.SP.Set(0x0000)

	cpu.AF.mask = 0xFFF0
}

func (cpu *CPU) setFlag(index byte, on bool) {
	if on {
		cpu.AF.SetLo(bits.Set(cpu.AF.Lo(), index))
	} else {
		cpu.AF.SetLo(bits.Reset(cpu.AF.Lo(), index))
	}
}

// SetZ sets value of Z flag.
func (cpu *CPU) SetZ(on bool) {
	cpu.setFlag(7, on)
}

// SetN sets N flag.
func (cpu *CPU) SetN(on bool) {
	cpu.setFlag(6, on)
}

// SetH sets H flag.
func (cpu *CPU) SetH(on bool) {
	cpu.setFlag(5, on)
}

// SetC sets C flag.
func (cpu *CPU) SetC(on bool) {
	cpu.setFlag(4, on)
}

// Z returns value of Z flag
func (cpu *CPU) Z() bool {
	return cpu.AF.HiLo()>>7&1 == 1
}

// N returns value of N flag
func (cpu *CPU) N() bool {
	return cpu.AF.HiLo()>>6&1 == 1
}

// H returns value of H flag
func (cpu *CPU) H() bool {
	return cpu.AF.HiLo()>>5&1 == 1
}

// C returns value of C flag
func (cpu *CPU) C() bool {
	return cpu.AF.HiLo()>>4&1 == 1
}
