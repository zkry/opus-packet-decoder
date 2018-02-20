package main

import (
	"bytes"
	"testing"
)

func TestInt16ToByteSlice(t *testing.T) {
	tests := []struct {
		input []int16
		want  []byte
	}{
		{
			input: []int16{0x0000, 0x0000, 0x0000, 0x0000},
			want:  []byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			input: []int16{0x0A0B, 0x0C0D, 0x0E0F, 0x0102},
			want:  []byte{0x0B, 0x0A, 0x0D, 0x0C, 0x0F, 0x0E, 0x02, 0x01},
		},
		{
			//			   0x8F8F  0x8182  0x8E81  0x8F88}
			input: []int16{-28785, -32382, -29055, -28792},
			want:  []byte{0x8F, 0x8F, 0x82, 0x81, 0x81, 0x8E, 0x88, 0x8F},
		},
	}

	for _, test := range tests {
		got := int16ToByteSlice(test.input)
		if !bytes.Equal(got, test.want) {
			t.Errorf("Error with input:%v\nGot: %v\nWanted: %v\n", test.input, got, test.want)
		}
	}
}
