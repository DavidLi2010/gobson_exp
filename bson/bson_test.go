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

	"github.com/DavidLi2010/gobson_exp/bson"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBson(t *testing.T) {
	Convey("test double type", t, func() {
		doc := bson.NewBson()
		doc.AppendDouble("_", 5.05)
		doc.Finish()

		data := doc.Raw()
		if bson.IsLittleEndian() {
			So(string(data), ShouldEqual, ("\x10\x00\x00\x00\x01_\x00333333\x14@\x00"))
		} else {

		}
	})

	Convey("test string type", t, func() {
		doc := bson.NewBson()
		doc.AppendString("_", "hello!")
		doc.Finish()

		data := doc.Raw()
		if bson.IsBigEndian() {
			So(string(data), ShouldEqual, ("\x10\x00\x00\x00\x02_\x00\x05hello\x00\x00"))
		} else {

		}
	})
}
