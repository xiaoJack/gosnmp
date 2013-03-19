// Copyright 2012 Andreas Louca. All rights reserved.
// Use of this source code is goverend by a BSD-style
// license that can be found in the LICENSE file.

package gosnmp

import (
	"encoding/asn1"
	"fmt"
)

type Asn1BER byte

const (
	Integer          Asn1BER = 0x02
	BitString                = 0x03
	OctetString              = 0x04
	Null                     = 0x05
	ObjectIdentifier         = 0x06
	Sequence                 = 0x30
	Counter32                = 0x41
	Gauge32                  = 0x42
	TimeTicks                = 0x43
	Opaque                   = 0x44
	NsapAddress              = 0x45
	Counter64                = 0x46
	Uinteger32               = 0x47
	NoSuchObject             = 0x80
	NoSuchInstance           = 0x81
	GetRequest               = 0xa0
	GetNextRequest           = 0xa1
	GetResponse              = 0xa2
	SetRequest               = 0xa3
	Trap                     = 0xa4
	GetBulkRequest           = 0xa5
)

// Different packet structure is needed during decode, to trick encoding/asn1 to decode the SNMP packet

type Variable struct {
	Name  []int
	Type  Asn1BER
	Size  uint64
	Value interface{}
}

type VarBind struct {
	Name  asn1.ObjectIdentifier
	Value asn1.RawValue
}

type PDU struct {
	RequestId   int32
	ErrorStatus int
	ErrorIndex  int
	VarBindList []VarBind
}
type PDUResponse struct {
	RequestId   int32
	ErrorStatus int
	ErrorIndex  int
	VarBindList []*Variable
}

type Message struct {
	Version   int
	Community []uint8
	Data      asn1.RawValue
}

func decodeValue(valueType Asn1BER, data []byte) (retVal *Variable, err error) {
	retVal = new(Variable)
	retVal.Size = uint64(len(data))

	switch Asn1BER(valueType) {

	// Integer
	case Integer:
		ret, err := parseInt(data)
		if err != nil {
			break
		}
		retVal.Type = Integer
		retVal.Value = ret
	// Octet
	case OctetString:
		retVal.Type = OctetString
		retVal.Value = string(data)
	case ObjectIdentifier:
		retVal.Type = ObjectIdentifier
		retVal.Value, _ = parseObjectIdentifier(data)
	// Counter32
	case Counter32:
		ret, err := parseInt(data)
		if err != nil {
			break
		}
		retVal.Type = Counter32
		retVal.Value = ret
	case TimeTicks:
		ret, err := parseInt(data)
		if err != nil {
			break
		}
		retVal.Type = TimeTicks
		retVal.Value = ret
	// Gauge32
	case Gauge32:
		ret, err := parseInt(data)
		if err != nil {
			break
		}
		retVal.Type = Gauge32
		retVal.Value = ret
	case Counter64:
		ret, err := parseInt64(data)

		// Decode it
		if err != nil {
			break
		}

		retVal.Type = Counter64
		retVal.Value = ret
	case Sequence:
		// NOOP
	case GetResponse:
		// NOOP
		retVal.Value = data
	case NoSuchInstance:
		return nil, fmt.Errorf("No such instance")
	case NoSuchObject:
		return nil, fmt.Errorf("No such object")
	default:
		err = fmt.Errorf("Unable to decode %x - not implemented", valueType)
	}

	return
}

// Parses UINT16
func ParseUint16(content []byte) int {
	number := uint8(content[1]) | uint8(content[0])<<8
	//fmt.Printf("\t%d\n", number)

	return int(number)
}
