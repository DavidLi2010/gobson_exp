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

import "strconv"

func init() {
	itoaCache = make([]string, itoaCacheSize)
	for i := 0; i < itoaCacheSize; i++ {
		itoaCache[i] = strconv.Itoa(i)
	}
}

const itoaCacheSize = 32

var itoaCache []string

func itoa(i int) string {
	if i < itoaCacheSize {
		return itoaCache[i]
	}
	return strconv.Itoa(i)
}

type BsonArrayBuilder struct {
	builder BsonBuilder
	index   int
}

func NewBsonArrayBuilder() *BsonArrayBuilder {
	bsonArrayBuilder := &BsonArrayBuilder{builder: BsonBuilder{raw: make([]byte, 0, initialBufferSize)}}
	bsonArrayBuilder.builder.reserveLength()
	return bsonArrayBuilder
}

func (a *BsonArrayBuilder) Finish() *BsonArrayBuilder {
	a.builder.Finish()
	return a
}

func (child *BsonArrayBuilder) AppendArrayEnd() (parent *BsonBuilder) {
	if child.builder.parent == nil {
		panic("not in child array")
	}
	if !child.builder.finished {
		panic("the child array is not finished")
	}
	if child.builder.raw[len(child.builder.raw)-1] != eod {
		panic("the child array is not finished")
	}
	parent = child.builder.parent
	parent.raw = child.builder.raw
	parent.child = nil
	parent.inChild = false
	child.builder.parent = nil
	return parent
}

func (a *BsonArrayBuilder) Raw() []byte {
	return a.builder.Raw()
}

func (a *BsonArrayBuilder) BsonArray() *BsonArray {
	return &BsonArray{bson: Bson{raw: a.Raw()}}
}

func (a *BsonArrayBuilder) AppendFloat64(value float64) *BsonArrayBuilder {
	a.builder.AppendFloat64(itoa(a.index), value)
	a.index++
	return a
}

func (a *BsonArrayBuilder) AppendString(value string) *BsonArrayBuilder {
	a.builder.AppendString(itoa(a.index), value)
	a.index++
	return a
}

func (a *BsonArrayBuilder) AppendBson(value *Bson) *BsonArrayBuilder {
	a.builder.AppendBson(itoa(a.index), value)
	a.index++
	return a
}

func (a *BsonArrayBuilder) AppendArray(value *BsonArray) *BsonArrayBuilder {
	a.builder.AppendArray(itoa(a.index), value)
	a.index++
	return a
}

func (a *BsonArrayBuilder) AppendBinary(value Binary) *BsonArrayBuilder {
	a.builder.AppendBinary(itoa(a.index), value)
	a.index++
	return a
}

func (a *BsonArrayBuilder) AppendObjectId(value ObjectId) *BsonArrayBuilder {
	a.builder.AppendObjectId(itoa(a.index), value)
	a.index++
	return a
}

func (a *BsonArrayBuilder) AppendBool(value bool) *BsonArrayBuilder {
	a.builder.AppendBool(itoa(a.index), value)
	a.index++
	return a
}

func (a *BsonArrayBuilder) AppendDate(value Date) *BsonArrayBuilder {
	a.builder.AppendDate(itoa(a.index), value)
	a.index++
	return a
}

func (a *BsonArrayBuilder) AppendNull() *BsonArrayBuilder {
	a.builder.AppendNull(itoa(a.index))
	a.index++
	return a
}

func (a *BsonArrayBuilder) AppendRegex(value RegEx) *BsonArrayBuilder {
	a.builder.AppendRegex(itoa(a.index), value)
	a.index++
	return a
}

func (a *BsonArrayBuilder) AppendInt32(value int32) *BsonArrayBuilder {
	a.builder.AppendInt32(itoa(a.index), value)
	a.index++
	return a
}

func (a *BsonArrayBuilder) AppendTimestamp(value Timestamp) *BsonArrayBuilder {
	a.builder.AppendTimestamp(itoa(a.index), value)
	a.index++
	return a
}

func (a *BsonArrayBuilder) AppendInt64(value int64) *BsonArrayBuilder {
	a.builder.AppendInt64(itoa(a.index), value)
	a.index++
	return a
}

func (a *BsonArrayBuilder) AppendMinKey() *BsonArrayBuilder {
	a.builder.AppendMinKey(itoa(a.index))
	a.index++
	return a
}

func (a *BsonArrayBuilder) AppendMaxKey() *BsonArrayBuilder {
	a.builder.AppendMaxKey(itoa(a.index))
	a.index++
	return a
}

func (a *BsonArrayBuilder) Append(value interface{}) *BsonArrayBuilder {
	a.builder.Append(itoa(a.index), value)
	a.index++
	return a
}
