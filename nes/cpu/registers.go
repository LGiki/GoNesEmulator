package cpu

type Registers struct {
	A    uint8           // Accumulator
	X, Y uint8           // Indexes
	SP   uint8           // Stack Pointer
	P    ProcessorStatus // Status Register
	PC   uint16          // Program Counter
}

type ProcessorStatus uint8

const (
	P_FLAG_CARRY        ProcessorStatus = 1 << iota // Carry flag
	P_FLAG_ZERO                                     // Zero flag
	P_FLAG_INTERUPT                                 // Interrupt disable flag
	P_FLAG_DECIMAL_MODE                             // Decimal mode flag. On the NES, this flag has no effect.
	P_FLAG_BREAK                                    // Break flag
	P_FLAG_UNUSED                                   // Unused flag, always 1
	P_FLAG_OVERFLOW                                 // Overflow flag
	P_FLAG_NEGATIVE                                 // Negative flag
)

func (r *Registers) Reset() {
	r.P = 0x34
	r.A = 0
	r.X = 0
	r.Y = 0
	r.SP = 0xfd
}

func (p *ProcessorStatus) Set(flag ProcessorStatus, value bool) {
	if value {
		*p |= flag
	} else {
		*p &= ^flag
	}
}
