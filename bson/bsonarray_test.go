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

	"github.com/DavidLi2010/gobson_exp/bson"
)

func TestBsonArray(t *testing.T) {
	var tests = []struct {
		bsonType bson.BsonType
		value    interface{}
		want     []byte
	}{
		{bson.BsonTypeFloat64, 5.05, []byte("\x10\x00\x00\x00\x010\x00333333\x14@\x00")},
	}

	for _, test := range tests {
		doc := bson.NewBsonArray()
		switch test.bsonType {
		case bson.BsonTypeFloat64:
			doc.AppendFloat64(test.value.(float64))
		case bson.BsonTypeString:
			doc.AppendString(test.value.(string))
		case bson.BsonTypeBinary:
			doc.AppendBinary(test.value.(bson.Binary))
		case bson.BsonTypeObjectId:
			doc.AppendObjectId(test.value.(bson.ObjectId))
		case bson.BsonTypeBool:
			doc.AppendBool(test.value.(bool))
		case bson.BsonTypeDate:
			doc.AppendDate(test.value.(bson.Date))
		case bson.BsonTypeNull:
			doc.AppendNull()
		case bson.BsonTypeRegEx:
			doc.AppendRegex(test.value.(bson.RegEx))
		case bson.BsonTypeInt32:
			doc.AppendInt32(test.value.(int32))
		case bson.BsonTypeInt64:
			doc.AppendInt64(test.value.(int64))
		case bson.BsonTypeTimestamp:
			doc.AppendTimestamp(test.value.(bson.Timestamp))
		case bson.BsonTypeMaxKey:
			doc.AppendMaxKey()
		case bson.BsonTypeMinKey:
			doc.AppendMinKey()
		default:
			panic("invalid bson type")
		}
		doc.Finish()

		data := doc.Raw()

		if len(data) != len(test.want) {
			t.Errorf("type: %v\nexpected: %v\n  actual: %v", test.bsonType, test.want, data)
		} else {
			for i, b := range data {
				if b != test.want[i] {
					t.Errorf("type: %v\nexpected: %v\n  actual: %v", test.bsonType, test.want, data)
				}
			}
		}
	}
}
