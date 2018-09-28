package cas

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// Cas lib context
type Cas struct {
	conn net.Conn
}

func checksum(data []byte) byte {
	var sum byte

	for _, v := range data {
		sum += v
	}

	return sum
}

// PLUName1String is a special type of Name1
type PLUName1String [40]byte

// PLUName2String is a special type of Name2
type PLUName2String [40]byte

// PLUName3String is a special type of Name3
type PLUName3String [5]byte

// PLUData contains all product info
type PLUData struct {
	DepartmentNumber       uint16         // 0
	PLUNumber              uint32         // 2
	PLUType                uint8          // 6
	PLUName1               PLUName1String // 7
	PLUName2               [40]byte       // 47
	PLUName3               [5]byte        // 87
	GroupNumber            uint16         // 92
	LabelNumber            uint16         // 94
	AuxLabelNumber         uint16         // 96
	OriginNumber           uint16         // 98
	UnitWeightNumber       uint8          // 100
	FixedWeight            uint32         // 101
	ItemCode               uint32         // 105
	PCSQuntity             uint16         // 109 !!! keyboard key index !!!
	PCSQuntitySymbolNumber uint8          // 111
	UseFixPriceType        uint8          // 112
	UnitPrice              uint32         // 113
	SpecialPrice           uint32         // 117
	TareWeight             uint32         // 121
	TareNumber             uint8          // 125
	BarcodeNumber          uint16         // 126
	AuxBarcodeNumber       uint16         // 128
	ProducedDate           uint16         // 130
	PackedDate             uint16         // 132
	PackedTime             uint8          // 134
	SellByDate             uint32         // 135
	SellByTime             uint8          // 139
	MessageNumber          uint16         // 140
	Reserved0              uint16         // 142
	Reserved1              uint16         // 144
	SaleMessageNumber      uint8          // 146
}

// Status Protocol Structure
type Status struct {
	LoadFlag           uint8 // 0: Zero 1: Non zero 2: Overload
	StableFlag         uint8 // 0: Unstable 1: Stable
	TareFlag           uint8
	DualRage           uint8
	WeightUnit         uint8
	WeightDecimalPoint uint8
	PriceDecimalPoint  uint8
	Reserved           uint8
	Tare               uint32
	Weight             int32
	UnitPrice          uint32
	TotalPrice         uint32
	PLUNumber          uint32
	DepartmentNumber   uint16
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

// UploadPLU requests data from scale
func (c Cas) UploadPLU(scaleID uint32, number uint32) (PLUData, error) {
	var data PLUData
	var err error
	var tmp []byte

	opcode := [2]byte{'R', 'L'}

	plu := make([]byte, 4)
	binary.LittleEndian.PutUint32(plu, number)
	buf := encodePacket(scaleID, opcode, plu)

	if _, err = c.conn.Write(buf); err != nil {
		goto End
	}

	tmp = make([]byte, 512)

	if _, err = c.conn.Read(tmp); err != nil {
		goto End
	}

	if tmp[0] != 'W' {
		return data, fmt.Errorf("get %#x opcode[0]", tmp[0])
	}

	if tmp[1] != opcode[1] {
		return data, fmt.Errorf("get %#x opcode[1]", tmp[1])
	}

	// TODO: checksum

	// 10 header + 4 room number or 4 DeptPLU
	if err = binary.Read(bytes.NewReader(tmp[14:]), binary.LittleEndian, &data); err != nil {
		goto End
	}

End:
	return data, nil
}

// DownloadPLU send PLUData to scale
func (c Cas) DownloadPLU(scaleID uint32, data PLUData) error {
	opcode := [2]byte{'W', 'L'}
	var dataBuf bytes.Buffer

	if err := binary.Write(&dataBuf, binary.LittleEndian, data); err != nil {
		return err
	}

	buf := encodePacket(scaleID, opcode, dataBuf.Bytes())

	if _, err := c.conn.Write(buf); err != nil {
		return err
	}

	tmp := make([]byte, 512)

	if _, err := c.conn.Read(tmp); err != nil {
		return err
	}

	if tmp[0] != 'G' || tmp[1] != opcode[1] {
		return fmt.Errorf("DownloadPLU %s %x", string(tmp[0:2]), tmp)
	}

	return nil
}

// ErasePLU delete one PLU or all if departmentNumber = 0 and PLUNumber = 0
func (c Cas) ErasePLU(scaleID uint32, departmentNumber uint16, PLUNumber uint32) error {
	opcode := [2]byte{'W', 'L'}
	plu := make([]byte, 6)
	binary.LittleEndian.PutUint16(plu, departmentNumber)
	binary.LittleEndian.PutUint32(plu[2:], PLUNumber)
	buf := encodePacket(scaleID, opcode, plu)

	if _, err := c.conn.Write(buf); err != nil {
		return err
	}

	tmp := make([]byte, 512)

	if _, err := c.conn.Read(tmp); err != nil {
		return err
	}

	if tmp[0] != 'G' || tmp[1] != opcode[1] {
		return fmt.Errorf("ErasePLU %s %x", string(tmp[0:2]), tmp)
	}

	return nil
}

// GetStatus returns current scale state
func (c Cas) GetStatus(scaleID uint32) (Status, error) {
	var status Status

	opcode := [2]byte{'R', 'N'}
	buf := encodePacket(scaleID, opcode, []byte{})

	if _, err := c.conn.Write(buf); err != nil {
		return status, err
	}

	tmp := make([]byte, 512)

	if _, err := c.conn.Read(tmp); err != nil {
		return status, err
	}

	if tmp[0] != 'W' || tmp[1] != 'N' {
		return status, fmt.Errorf("GetStatus %s %x", string(tmp[0:2]), tmp)
	}

	return status, nil
}

// Connect to scale
func Connect(addr string) (Cas, error) {
	var cas Cas
	var err error

	cas.conn, err = net.DialTimeout("tcp", addr, 5*time.Second)

	if err != nil {
		return cas, err
	}

	return cas, nil
}

// Close connection
func (c Cas) Disconnect() error {
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}
