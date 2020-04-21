package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

const LogFileName = "txnlog.dat"
const HeaderMagicValue = "MPS7"
const SpecialUserID = 2456938384156277127

type RecordType byte

const (
	DebitType        RecordType = 0x00
	CreditType                  = 0x01
	StartAutopayType            = 0x02
	EndAutopayType              = 0x03
)

// Header is the first line of the transaction log
type Header struct {
	Magic       [4]byte
	Version     byte
	RecordCount uint32
}

// MinRecord is the set of fields that all records have in common; some also have an amount
type MinRecord struct {
	Type      RecordType
	Timestamp uint32
	UserID    uint64
}

// Sum is used to track the summary info of the log
type Sum struct {
	Credits       float64
	Debits        float64
	AutopayStarts int64
	AutopayEnds   int64
	UserBalance   float64
}

func main() {
	logFile, err := os.Open(LogFileName)
	if err != nil {
		log.Panic("Unable to read log file", err)
	}
	defer logFile.Close()

	// read the header first
	h := Header{}
	err = binary.Read(logFile, binary.BigEndian, &h)
	if err != nil {
		log.Panic("Bad header read", err)
	}

	if string(h.Magic[:]) != HeaderMagicValue {
		log.Panic("Bad Magic")
	}

	// read each record and store summary info
	sum := Sum{}
	for i := uint32(0); i < h.RecordCount; i++ {
		if err := readRecord(logFile, &sum); err != nil {
			log.Panic("Bad read", err)
		}
	}

	fmt.Printf(`
total credit amount=%.2f
total debit amount=%.2f
autopays started=%d
autopays ended=%d
balance for user %d=%.2f
`,
		sum.Credits,
		sum.Debits,
		sum.AutopayStarts,
		sum.AutopayEnds,
		SpecialUserID,
		sum.UserBalance)
}

// ReadRecord consumes a record in the provided buf and captures summary info from that record
// Returns an error if it can't read the record or doesn't recognize the transaction type
func readRecord(buf io.Reader, sum *Sum) (err error) {
	rec := MinRecord{}
	if err = binary.Read(buf, binary.BigEndian, &rec); err != nil {
		return err
	}

	switch rec.Type {
	case DebitType, CreditType:
		// then we also need to read the amount for this transaction
		var amount float64
		binary.Read(buf, binary.BigEndian, &amount)

		if rec.Type == DebitType {
			sum.Debits += amount
			amount = -1 * amount
		} else {
			sum.Credits += amount
		}

		if rec.UserID == SpecialUserID {
			sum.UserBalance += amount
		}

	case StartAutopayType:
		sum.AutopayStarts++

	case EndAutopayType:
		sum.AutopayEnds++
	default:
		return fmt.Errorf("Unknown tran type: %v", rec.Type)
	}

	return nil
}
