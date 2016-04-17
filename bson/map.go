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

import (
	"fmt"
	"math"
	"reflect"
)

type Map map[string]interface{}

func (m Map) toBson(b *Bson) {
	for name, v := range m {
		switch v.(type) {
		case float32:
			b.AppendFloat64(name, float64(v.(float32)))
		case float64:
			b.AppendFloat64(name, v.(float64))
		case int8:
			b.AppendInt32(name, int32(v.(int8)))
		case int16:
			b.AppendInt32(name, int32(v.(int16)))
		case int32:
			b.AppendInt32(name, v.(int32))
		case int64:
			val := v.(int64)
			if val >= math.MinInt32 && val <= math.MaxInt32 {
				b.AppendInt32(name, int32(val))
			} else {
				b.AppendInt64(name, int64(val))
			}
		case uint8:
			b.AppendInt32(name, int32(v.(uint8)))
		case uint16:
			b.AppendInt32(name, int32(v.(uint16)))
		case uint32:
			val := v.(uint32)
			if int32(val) < 0 {
				b.AppendInt64(name, int64(val))
			} else {
				b.AppendInt32(name, int32(val))
			}
		case uint64:
			val := int64(v.(uint64))
			if val < 0 {
				panic("bson has no uint64 type, and value is too large to fit correctly in an int64")
			}
			if val >= math.MinInt32 && val <= math.MaxInt32 {
				b.AppendInt32(name, int32(val))
			} else {
				b.AppendInt64(name, int64(val))
			}
		case int:
			val := int64(v.(int))
			if val >= math.MinInt32 && val <= math.MaxInt32 {
				b.AppendInt32(name, int32(val))
			} else {
				b.AppendInt64(name, int64(val))
			}
		case uint:
			val := int64(v.(uint))
			if val < 0 {
				panic("bson has no uint64 type, and value is too large to fit correctly in an int64")
			}
			if val >= math.MinInt32 && val <= math.MaxInt32 {
				b.AppendInt32(name, int32(val))
			} else {
				b.AppendInt64(name, int64(val))
			}
		case bool:
			b.AppendBool(name, v.(bool))
		case string:
			b.AppendString(name, v.(string))
		case nil:
			b.AppendNull(name)
		case ObjectId:
			b.AppendObjectId(name, v.(ObjectId))
		case Date:
			b.AppendDate(name, v.(Date))
		case RegEx:
			b.AppendRegex(name, v.(RegEx))
		case Timestamp:
			b.AppendTimestamp(name, v.(Timestamp))
		case Binary:
			b.AppendBinary(name, v.(Binary))
		case orderKey:
			val := v.(orderKey)
			if val == MaxKey {
				b.AppendMaxKey(name)
			} else if val == MinKey {
				b.AppendMinKey(name)
			} else {
				panic("invalid orderkey")
			}
		case Map:
			b.AppendBson(name, v.(Map).Bson())
		default:
			// complex64, complex128
			panic(fmt.Errorf("can't append %s to bson", reflect.TypeOf(v).String()))
		}
	}
}

func (m Map) Bson() *Bson {
	b := NewBson()
	m.toBson(b)
	b.Finish()
	return b
}
