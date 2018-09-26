package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"golang.org/x/text/encoding/charmap"
)

func checksum(data []byte) byte {
	var sum byte = 0

	for _, v := range data {
		sum += v
	}

	return sum
}

type PLUData struct {
	DepartmentNumber       uint16   // 0
	PLUNumber              uint32   // 2
	PLUType                uint8    // 6
	PLUName1               [40]byte // 7
	PLUName2               [40]byte // 47
	PLUName3               [5]byte  // 87
	GroupNumber            uint16   // 92
	LabelNumber            uint16   // 94
	AuxLabelNumber         uint16   // 96
	OriginNumber           uint16   // 98
	UnitWeightNumber       uint8    // 100
	FixedWeight            uint32   // 101
	ItemCode               uint32   // 105
	PCSQuntity             uint16   // 109 !!! keyboard key index !!!
	PCSQuntitySymbolNumber uint8    // 111
	UseFixPriceType        uint8    // 112
	UnitPrice              uint32   // 113
	SpecialPrice           uint32   // 117
	TareWeight             uint32   // 121
	TareNumber             uint8    // 125
	BarcodeNumber          uint16   // 126
	AuxBarcodeNumber       uint16   // 128
	ProducedDate           uint16   // 130
	PackedDate             uint16   // 132
	PackedTime             uint8    // 134
	SellByDate             uint32   // 135
	SellByTime             uint8    // 139
	MessageNumber          uint16   // 140
	Reserved0              uint16   // 142
	Reserved1              uint16   // 144
	SaleMessageNumber      uint8    // 146
}

func encodePacket(address uint32, opcode [2]byte, data []byte) []byte {
	buf := make([]byte, 10)

	buf[0] = opcode[0]
	buf[1] = opcode[1]
	binary.LittleEndian.PutUint32(buf[2:], address)
	buf[6] = ','
	binary.LittleEndian.PutUint16(buf[7:], uint16(len(data)))
	buf[9] = ':'
	buf = append(buf, data...)
	buf = append(buf, ':')
	sum := checksum(buf[2:])
	buf = append(buf, sum, 0x0D)

	return buf
}

func uploadPLU(conn net.Conn, scaleId uint32, number uint32) PLUData {
	opcode := [2]byte{'R', 'L'}

	plu := make([]byte, 4)
	binary.LittleEndian.PutUint32(plu, number)
	buf := encodePacket(0, opcode, plu)

	conn.Write(buf)

	tmp := make([]byte, 512)
	conn.Read(tmp)

	if tmp[0] != 'W' {
		fmt.Println("error")
	}

	if tmp[1] != opcode[1] {
		fmt.Println("error")
	}

	var data PLUData
	binary.Read(bytes.NewReader(tmp[10:]), binary.LittleEndian, &data)

	return data
}

func downloadPLU(conn net.Conn, scaleId uint32, data PLUData) {
	opcode := [2]byte{'W', 'L'}
	var dataBuf bytes.Buffer
	binary.Write(&dataBuf, binary.LittleEndian, data)
	buf := encodePacket(0, opcode, dataBuf.Bytes())

	conn.Write(buf)

	tmp := make([]byte, 512)
	conn.Read(tmp)

	if tmp[0] != 'G' {
		fmt.Printf("error %s %x\n", string(tmp[0:2]), tmp)
	}

	if tmp[1] != opcode[1] {
		fmt.Printf("error %s %x\n", string(tmp[0:2]), tmp)
	}
}

func erasePLU(conn net.Conn, departmentNumber uint16, PLUNumber uint32) {
	opcode := [2]byte{'W', 'L'}
	plu := make([]byte, 6)
	binary.LittleEndian.PutUint16(plu, departmentNumber)
	binary.LittleEndian.PutUint32(plu[2:], PLUNumber)
	buf := encodePacket(0, opcode, plu)

	conn.Write(buf)

	tmp := make([]byte, 512)
	conn.Read(tmp)

	if tmp[0] != 'G' {
		fmt.Printf("error %s %x\n", string(tmp[0:2]), tmp)
	}

	if tmp[1] != opcode[1] {
		fmt.Printf("error %s %x\n", string(tmp[0:2]), tmp)
	}
}

func main() {
	conn, _ := net.Dial("tcp", "192.168.89.231:20304")

	// erasePLU(conn, 0, 0)

	for i := 10; i >= 0; i -= 1 {
		name, _ := charmap.Windows1251.NewEncoder().String(fmt.Sprintf("button %d", i))
		var data PLUData
		copy(data.PLUName1[:], name)
		data.DepartmentNumber = 1
		data.PLUNumber = uint32(i)
		data.PLUType = 1
		downloadPLU(conn, 0, data)
	}

	dec := charmap.Windows1251.NewDecoder()

	for i := 1; i < 10; i += 1 {
		data := uploadPLU(conn, 0, uint32(i))
		out, _ := dec.Bytes(data.PLUName1[0:])
		fmt.Printf("%d %s %d %+v \n", i, string(out), data.PLUNumber, data)
	}
}
