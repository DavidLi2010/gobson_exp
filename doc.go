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

package bson

type Doc struct {
	raw []byte
}

type Date struct {
	ms  int32
	inc int32
}

type Timestamp int64

type Binary struct {
	subtype byte
	data    []byte
}

type ObjectId struct {
}

func (doc *Doc) appendDouble(name string, value float64) {

}

func (doc *Doc) appendString(name string, value string) {

}

func (doc *Doc) appendDoc(name string, value *Doc) {

}

func (doc *Doc) appendArray(name string, value []Doc) {

}

func (doc *Doc) appendBinary(name string, value Binary) {

}

func (doc *Doc) appendObjectId(name string, value ObjectId) {

}

func (doc *Doc) appendBoolean(name string, value bool) {

}

func (doc *Doc) appendDate(name string, value Date) {

}

func (doc *Doc) appendNull(name string) {

}

func (doc *Doc) appendRegex(name string, pattern string, options string) {

}

func (doc *Doc) appendInt32(name string, value int32) {

}

func (doc *Doc) appendTimestamp(name string, value Timestamp) {

}

func (doc *Doc) appendInt64(name string, value int64) {

}

func (doc *Doc) appendMinKey(name string) {

}

func (doc *Doc) appendMaxKey(name string) {

}
