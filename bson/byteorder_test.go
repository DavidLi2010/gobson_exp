// Copyright 2015-2016 David Li
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package bson_test

import (
	"testing"
	"unsafe"

	. "github.com/DavidLi2010/gobson_exp/bson"
)

func TestByteOrder(t *testing.T) {
	if BigEndian.IsBigEndian() == false {
		t.Errorf("BigEndian.IsBigEndian() == false")
	}

	if BigEndian.IsLittleEndian() == true {
		t.Errorf("BigEndian.IsLittleEndian() == true")
	}

	if LittleEndian.IsBigEndian() == true {
		t.Errorf("LittleEndian.IsBigEndian() == true")
	}

	if LittleEndian.IsLittleEndian() == false {
		t.Errorf("LittleEndian.IsLittleEndian() == false")
	}

	byteOrder := GetByteOrder()
	if byteOrder.IsLittleEndian() {
		if byteOrder.IsBigEndian() {
			t.Errorf("invalid byte order: little endian")
		}
	} else {
		if byteOrder.IsLittleEndian() {
			t.Errorf("invalid byte order: big endian")
		}
	}
}

func TestAppend(t *testing.T) {
	x := int32(0x12345678)
	b := make([]byte, 0, 4)
	byteOrder := GetByteOrder()
	if byteOrder.IsLittleEndian() {
		p := *(*byte)(unsafe.Pointer(&x))
		if p != 0x78 {
			t.Errorf("invalid endian")
		}

		b = byteOrder.AppendInt32(b, x)
		if b[0] != 0x78 ||
			b[1] != 0x56 ||
			b[2] != 0x34 ||
			b[3] != 0x12 {
			t.Errorf("LittleEndian.AppendInt32() error")
		}

		b = make([]byte, 0, 4)
		b = BigEndian.AppendInt32(b, x)
		if b[0] != 0x12 ||
			b[1] != 0x34 ||
			b[2] != 0x56 ||
			b[3] != 0x78 {
			t.Errorf("BigEndian.AppendInt32() error")
		}
	} else {
		if byteOrder.IsBigEndian() != true {
			t.Errorf("invalid endian")
		}

		p := *(*byte)(unsafe.Pointer(&x))
		if p != 0x12 {
			t.Errorf("invalid endian")
		}

		b = byteOrder.AppendInt32(b, x)
		if b[0] != 0x12 ||
			b[1] != 0x34 ||
			b[2] != 0x56 ||
			b[3] != 0x78 {
			t.Errorf("BigEndian.AppendInt32() error")
		}

		b = make([]byte, 0, 4)
		b = LittleEndian.AppendInt32(b, x)
		if b[0] != 0x78 ||
			b[1] != 0x56 ||
			b[2] != 0x34 ||
			b[3] != 0x12 {
			t.Errorf("LittleEndian.AppendInt32() error")
		}
	}
}

func TestLittleEndianReadValue(t *testing.T) {
	in := int32(0x123456)
	b := []byte{}
	b = LittleEndian.AppendInt32(b, in)
	out := LittleEndian.Int32(b)
	if in != out {
		t.Errorf("LittleEndian.Int32(%v)=%v", in, out)
	}

	in64 := int64(0x1234567890123456)
	b = []byte{}
	b = LittleEndian.AppendInt64(b, in64)
	out64 := LittleEndian.Int64(b)
	if in != out {
		t.Errorf("LittleEndian.Int64(%v)=%v", in64, out64)
	}

	inf64 := float64(12345678.0123456)
	b = []byte{}
	b = LittleEndian.AppendFloat64(b, inf64)
	outf64 := LittleEndian.Float64(b)
	if in != out {
		t.Errorf("LittleEndian.Float64(%v)=%v", inf64, outf64)
	}
}

func TestBigEndianReadValue(t *testing.T) {
	in := int32(0x123456)
	b := []byte{}
	b = BigEndian.AppendInt32(b, in)
	out := BigEndian.Int32(b)
	if in != out {
		t.Errorf("BigEndian.Int32(%v)=%v", in, out)
	}

	in64 := int64(0x1234567890123456)
	b = []byte{}
	b = BigEndian.AppendInt64(b, in64)
	out64 := BigEndian.Int64(b)
	if in != out {
		t.Errorf("BigEndian.Int64(%v)=%v", in64, out64)
	}

	inf64 := float64(12345678.0123456)
	b = []byte{}
	b = BigEndian.AppendFloat64(b, inf64)
	outf64 := BigEndian.Float64(b)
	if in != out {
		t.Errorf("BigEndian.Float64(%v)=%v", inf64, outf64)
	}
}
