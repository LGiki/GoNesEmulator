package ines

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const (
	INES_MAGIC_NUMBER = "NES\x1a"
	TRAINER_SIZE      = 512 // 512-byte trainer at $7000-$71FF (stored before PRG data)
)

type Header struct {
	MagicNumber [4]byte // Constant $4E $45 $53 $1A ("NES" followed by MS-DOS end-of-file)
	PRGSize     byte    // Size of PRG ROM in 16 KB units
	CHRSize     byte    // Size of CHR ROM in 8 KB units (Value 0 means the board uses CHR RAM)
	Flags6      byte    // Mapper, mirroring, battery, trainer
	Flags7      byte    // Mapper, VS/Playchoice, NES 2.0
	Flags8      byte    // PRG-RAM size (rarely used extension)
	Flags9      byte    // TV system (rarely used extension)
	Flags10     byte    // TV system, PRG-RAM presence (unofficial, rarely used extension)
	_           [5]byte // Unused padding
}

type Rom struct {
	Header     Header
	Trainer    []byte
	PRG        []byte
	CHR        []byte
	Extra      []byte
	MapperType byte // mapper type
	Battery    byte // battery-backed present
}

func LoadRom(reader io.Reader) (*Rom, error) {
	rom := &Rom{}
	header := &rom.Header

	if err := binary.Read(reader, binary.LittleEndian, header); err != nil {
		return nil, err
	}

	if string(header.MagicNumber[:]) != INES_MAGIC_NUMBER {
		return nil, errors.New("not a valid ines file")
	}

	if header.IsTrainerPresent() {
		rom.Trainer = make([]byte, TRAINER_SIZE)
		if _, err := reader.Read(rom.Trainer); err != nil {
			return nil, err
		}
	}

	mapper1 := header.Flags6 >> 4
	mapper2 := header.Flags7 >> 4
	rom.MapperType = mapper1 | mapper2<<4

	rom.Battery = (header.Flags6 >> 1) & 1

	PRG := make([]byte, 16*1024*int(header.PRGSize))
	if _, err := reader.Read(PRG); err != nil {
		return nil, err
	}

	var CHR []byte

	if header.CHRSize > 0 {
		CHR = make([]byte, 8*1024*int(header.CHRSize))
		if _, err := reader.Read(CHR); err != nil {
			return nil, err
		}
	}

	rom.PRG = PRG
	rom.CHR = CHR

	extra := &bytes.Buffer{}
	if _, err := io.Copy(extra, reader); err != nil {
		return nil, err
	}
	return rom, nil
}

func (h *Header) IsTrainerPresent() bool {
	return h.Flags6&(1>>2) == 1
}

func (h *Header) IsHorizontalMirroring() bool {
	return h.Flags6&1 == 0
}

func (h *Header) IsVerticalMirroring() bool {
	return h.Flags6&1 == 1
}
