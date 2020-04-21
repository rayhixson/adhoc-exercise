package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

const logFileName = "txnlog.dat"
const MagicValue = "MPS7"

type RecordType byte

const (
	Debit        = 0
	Credit       = 1
	StartAutopay = 2
	EndAutopay   = 3
)

type Header struct {
	Magic       [4]byte
	Version     byte
	RecordCount [4]byte
}

type Record struct {
	Type      byte
	Timestamp [4]byte
	UserID    [8]byte
}

func main() {
	logFile, err := os.Open(logFileName)
	if err != nil {
		log.Panic("Unable to read log file", err)
	}
	defer logFile.Close()

	h := Header{}
	err = binary.Read(logFile, binary.BigEndian, &h)
	if err != nil {
		log.Panic("Bad header read", err)
	}

	if string(h.Magic[:]) != MagicValue {
		log.Panic("Bad Magic")
	}

	fmt.Printf("%v, %v, %v\n", string(h.Magic[:]),
		h.Version,
		binary.BigEndian.Uint32(h.RecordCount[:]))

	rec, amount, err := ReadRecord(logFile)
	if err != nil {
		log.Panic("Bad read", err)
	}

	fmt.Printf("%v, %v, %v, %v\n",
		rec.Type,
		binary.BigEndian.Uint32(rec.Timestamp[:]),
		binary.BigEndian.Uint64(rec.UserID[:]),
		amount)
}

func ReadRecord(buf io.Reader) (rec Record, amount float64, err error) {
	err = binary.Read(buf, binary.BigEndian, &rec)

	if err != nil {
		return rec, 0, err
	}

	switch rec.Type {
	case Debit:
		binary.Read(buf, binary.BigEndian, &amount)
		return rec, (-1 * amount), nil
	case Credit:
		binary.Read(buf, binary.BigEndian, &amount)
		return rec, amount, nil
	}

	return rec, 0, nil
}
