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

package bson

import (
	"encoding/binary"
	"unsafe"
)

var byteOrder binary.ByteOrder
var reverseByteOrder binary.ByteOrder

func init() {
	const N int = int(unsafe.Sizeof(0))
	x := 0x1234
	p := unsafe.Pointer(&x)
	p2 := (*[N]byte)(p)
	if p2[0] == 0 {
		byteOrder = binary.BigEndian
		reverseByteOrder = binary.LittleEndian
	} else {
		byteOrder = binary.LittleEndian
		reverseByteOrder = binary.BigEndian
	}
}

func GetByteOrder() binary.ByteOrder {
	return byteOrder
}

func GetReverseByteOrder() binary.ByteOrder {
	return reverseByteOrder
}

func IsLittleEndian() bool {
	return byteOrder == binary.LittleEndian
}

func IsBigEndian() bool {
	return byteOrder == binary.BigEndian
}

func RevertInt32(v int32) int32 {
	var b [4]byte
	buf := b[:]
	byteOrder.PutUint32(buf, uint32(v))
	return int32(reverseByteOrder.Uint32(buf))
}
