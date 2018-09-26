package main

import "testing"

func TestChecksum(t *testing.T) {
	in := []byte{0x01, 0, 0, 0, 0x2C, 0x04, 0, 0x3A, 0x41, 0x42, 0x0F, 0, 0x3A}
	out := checksum(in)

	if out != 0x37 {
		t.Errorf("wrong checksum %x %x", out, 0x37)
	}
}
