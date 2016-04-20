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

import "fmt"

type BsonType byte

const (
	BsonTypeEOD BsonType = iota // end of doc
	BsonTypeFloat64
	BsonTypeString
	BsonTypeBson
	BsonTypeArray
	BsonTypeBinary
	BsonTypeUndefined // deprecated
	BsonTypeObjectId
	BsonTypeBool
	BsonTypeDate
	BsonTypeNull
	BsonTypeRegEx
	BsonTypeDBPointer  // deprecated
	BsonTypeCode       // not support
	BsonTypeSymbol     // deprecated
	BsonTypeCodeWScope // not support
	BsonTypeInt32
	BsonTypeTimestamp
	BsonTypeInt64
	BsonTypeMaxKey BsonType = 0x7F
	BsonTypeMinKey BsonType = 0xFF
)

type BinaryType byte

const (
	BinaryTypeGeneral BinaryType = iota
	BinaryTypeFunction
	BinaryTypeBinaryDeprecated
	BinaryTypeUUIDDeprecated
	BinaryTypeUUID
	BinaryTypeMD5
	BinaryTypeUser BinaryType = 0x80
)

type Date int64

func (d Date) String() string {
	return fmt.Sprintf(`{"$date":%v}`, d)
}

type RegEx struct {
	Pattern string
	Options string
}

func (re RegEx) String() string {
	return fmt.Sprintf(`{"$regex":"%s", "$options":"%s"}`, re.Pattern, re.Options)
}

type Timestamp struct {
	Second    int32
	Increment int32
}

func (t Timestamp) String() string {
	return fmt.Sprintf(`{"$timestamp":"%d %d"}`, t.Second, t.Increment)
}

type Binary struct {
	Subtype BinaryType
	Data    []byte
}

func (b Binary) String() string {
	return fmt.Sprintf(`{"$binary":"%s", "$type":"%d"}`, string(b.Data), b.Subtype)
}

type orderKey int64

func (o orderKey) String() string {
	if o == MaxKey {
		return `{"$maxKey":1}`
	} else if o == MinKey {
		return `{"$minKey":1}`
	} else {
		panic("invalid order key")
	}
}

var MaxKey = orderKey(1<<63 - 1)

var MinKey = orderKey(-1 << 63)
