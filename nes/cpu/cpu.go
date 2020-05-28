package cpu

import "GoNesEmulator/nes/memory"

type CPU struct {
	Registers Registers
	Memory memory.Memory
}