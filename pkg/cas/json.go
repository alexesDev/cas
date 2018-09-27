package cas

import (
	"bytes"
	"encoding/json"
	"golang.org/x/text/encoding/charmap"
)

func encode(out []byte, data []byte) error {
	var str string
	json.Unmarshal(data, &str)
	encoder := charmap.Windows1251.NewEncoder()
	val, err := encoder.String(str)

	if err == nil {
		copy(out, val)
	}

	return err
}

func decode(data []byte) ([]byte, error) {
	decoder := charmap.Windows1251.NewDecoder()
	out, err := decoder.Bytes(bytes.Trim(data, "\x00"))

	if err != nil {
		return nil, err
	}

	return json.Marshal(string(out))
}

func (s *PLUName1String) UnmarshalJSON(data []byte) error {
	return encode(s[:], data)
}

func (s *PLUName2String) UnmarshalJSON(data []byte) error {
	return encode(s[:], data)
}

func (s *PLUName3String) UnmarshalJSON(data []byte) error {
	return encode(s[:], data)
}

func (s PLUName1String) MarshalJSON() ([]byte, error) {
	return decode(s[:])
}

func (s PLUName2String) MarshalJSON() ([]byte, error) {
	return decode(s[:])
}

func (s PLUName3String) MarshalJSON() ([]byte, error) {
	return decode(s[:])
}
