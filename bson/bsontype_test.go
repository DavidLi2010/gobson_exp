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

	"github.com/davidli2010/gobson_exp/bson"
)

func TestBsonType(t *testing.T) {
	var tests = []struct {
		bsonType bson.BsonType
		want     byte
	}{
		{bson.BsonTypeEOD, 0x00},
		{bson.BsonTypeFloat64, 0x01},
		{bson.BsonTypeString, 0x02},
		{bson.BsonTypeBson, 0x03},
		{bson.BsonTypeArray, 0x04},
		{bson.BsonTypeBinary, 0x05},
		{bson.BsonTypeUndefined, 0x06},
		{bson.BsonTypeObjectId, 0x07},
		{bson.BsonTypeBool, 0x08},
		{bson.BsonTypeDate, 0x09},
		{bson.BsonTypeNull, 0x0A},
		{bson.BsonTypeRegEx, 0x0B},
		{bson.BsonTypeDBPointer, 0x0C},
		{bson.BsonTypeCode, 0x0D},
		{bson.BsonTypeSymbol, 0x0E},
		{bson.BsonTypeCodeWScope, 0x0F},
		{bson.BsonTypeInt32, 0x10},
		{bson.BsonTypeTimestamp, 0x11},
		{bson.BsonTypeInt64, 0x12},
		{bson.BsonTypeMinKey, 0xFF},
		{bson.BsonTypeMaxKey, 0x7F},
	}

	for _, test := range tests {
		if test.bsonType != bson.BsonType(test.want) {
			t.Errorf("invalid value of bson type: %v, expected: %v", test.bsonType, test.want)
		}
	}
}

func TestBinaryType(t *testing.T) {
	var tests = []struct {
		binaryType bson.BinaryType
		want       byte
	}{
		{bson.BinaryTypeGeneral, 0x00},
		{bson.BinaryTypeFunction, 0x01},
		{bson.BinaryTypeBinaryDeprecated, 0x02},
		{bson.BinaryTypeUUIDDeprecated, 0x03},
		{bson.BinaryTypeUUID, 0x04},
		{bson.BinaryTypeMD5, 0x05},
		{bson.BinaryTypeUser, 0x80},
	}

	for _, test := range tests {
		if test.binaryType != bson.BinaryType(test.want) {
			t.Errorf("invalid value of binary type: %v, expected: %v", test.binaryType, test.want)
		}
	}
}
