// Copyright 2015 David Li
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
	. "github.com/smartystreets/goconvey/convey"
)

func TestByteOrder(t *testing.T) {
	Convey("test byte order", t, func() {
		So(BigEndian.IsBigEndian(), ShouldBeTrue)
		So(BigEndian.IsLittleEndian(), ShouldBeFalse)
		So(LittleEndian.IsBigEndian(), ShouldBeFalse)
		So(LittleEndian.IsLittleEndian(), ShouldBeTrue)

		byteOrder := GetByteOrder()
		t := byteOrder.IsLittleEndian() || byteOrder.IsBigEndian()
		So(t, ShouldBeTrue)
		So(byteOrder.IsLittleEndian(), ShouldNotEqual, byteOrder.IsBigEndian())
	})
}

func TestAppend(t *testing.T) {
	Convey("test int32", t, func() {
		x := int32(0x12345678)
		b := make([]byte, 0, 4)
		byteOrder := GetByteOrder()
		if byteOrder.IsLittleEndian() {
			p := *(*byte)(unsafe.Pointer(&x))
			So(p, ShouldEqual, 0x78)

			b = byteOrder.AppendInt32(b, x)
			So(b[0], ShouldEqual, 0x78)
			So(b[1], ShouldEqual, 0x56)
			So(b[2], ShouldEqual, 0x34)
			So(b[3], ShouldEqual, 0x12)

			b = make([]byte, 0, 4)
			b = BigEndian.AppendInt32(b, x)
			So(b[0], ShouldEqual, 0x12)
			So(b[1], ShouldEqual, 0x34)
			So(b[2], ShouldEqual, 0x56)
			So(b[3], ShouldEqual, 0x78)
		} else {
			So(byteOrder.IsBigEndian(), ShouldBeTrue)

			p := *(*byte)(unsafe.Pointer(&x))
			So(p, ShouldEqual, 0x12)

			b = byteOrder.AppendInt32(b, x)
			So(b[0], ShouldEqual, 0x12)
			So(b[1], ShouldEqual, 0x34)
			So(b[2], ShouldEqual, 0x56)
			So(b[3], ShouldEqual, 0x78)

			b = make([]byte, 0, 4)
			b = LittleEndian.AppendInt32(b, x)
			So(b[0], ShouldEqual, 0x78)
			So(b[1], ShouldEqual, 0x56)
			So(b[2], ShouldEqual, 0x34)
			So(b[3], ShouldEqual, 0x12)
		}
	})
}
