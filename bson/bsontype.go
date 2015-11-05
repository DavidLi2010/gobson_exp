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

type BsonType byte

const (
	BsonTypeEOD BsonType = iota // end of doc
	BsonTypeDouble
	BsonTypeString
	BsonTypeDoc
	BsonTypeArray
	BsonTypeBinary
	BsonTypeUndefined // deprecated
	BsonTypeObjectId
	BsonTypeBool
	BsonTypeDate
	BsonTypeNull
	BsonTypeRegEx
	BsonTypeDBPointer // deprecated
	BsonTypeCode
	BsonTypeSymbol // deprecated
	BsonTypeCodeWScope
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

// end of cstring
const eos byte = 0x00

type Date int64

type RegEx struct {
	Pattern string
	Options string
}

type Timestamp struct {
	Second    int32
	Increment int32
}

type Binary struct {
	Subtype BinaryType
	Data    []byte
}
