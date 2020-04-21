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
const SpecialUser = uint64(2456938384156277127)

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

type Sum struct {
	UserID        uint64
	Credits       float64
	Debits        float64
	AutopayStarts int64
	AutopayEnds   int64
	UserBalance   float64
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

	recCount := binary.BigEndian.Uint32(h.RecordCount[:])
	fmt.Printf("Header: %v, %v, %v\n", string(h.Magic[:]),
		h.Version,
		recCount)

	sum := Sum{UserID: SpecialUser}
	for i := uint32(0); i < recCount; i++ {
		err := ReadRecord(logFile, &sum)
		if err != nil {
			log.Panic("Bad read", err)
		}
	}

	fmt.Printf("total credit amount=%.2f\ntotal debit amount=%.2f\nautopays started=%d\nautopays ended=%d\nbalance for user %d=%.2f",
		sum.Credits,
		sum.Debits,
		sum.AutopayStarts,
		sum.AutopayEnds,
		sum.UserID,
		sum.UserBalance)
}

func ReadRecord(buf io.Reader, sum *Sum) (err error) {
	rec := Record{}
	err = binary.Read(buf, binary.BigEndian, &rec)

	if err != nil {
		return err
	}

	var amount float64

	switch rec.Type {
	case Debit:
		binary.Read(buf, binary.BigEndian, &amount)
		sum.Debits += amount

		if binary.BigEndian.Uint64(rec.UserID[:]) == SpecialUser {
			sum.UserBalance -= amount
		}

	case Credit:
		binary.Read(buf, binary.BigEndian, &amount)
		sum.Credits += amount

		if binary.BigEndian.Uint64(rec.UserID[:]) == SpecialUser {
			sum.UserBalance += amount
		}

		fmt.Println(amount)
	case StartAutopay:
		sum.AutopayStarts++

	case EndAutopay:
		sum.AutopayEnds++
	default:
		return fmt.Errorf("Unknown tran type: %v", rec.Type)
	}

	/*
		fmt.Printf("%v, %v, %v, %v\n",
			rec.Type,
			binary.BigEndian.Uint32(rec.Timestamp[:]),
			binary.BigEndian.Uint64(rec.UserID[:]),
			amount)
	*/

	return nil
}
