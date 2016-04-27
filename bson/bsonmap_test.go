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

func TestMap(t *testing.T) {
	m := bson.Map{
		"int":     int(1),
		"int8":    int8(2),
		"float64": float64(123.456),
		"bson": map[string]interface{}{
			"i_int": int(10),
		},
		"array":    []interface{}{int(1), float64(123.456), "hello", bson.NewObjectId()},
		"objectid": bson.NewObjectId(),
		"bool":     true,
		"null":     nil,
		"maxkey":   bson.MaxKey,
		"minkey":   bson.MinKey,
		"string":   "hello",
	}

	m2 := m.Bson().Map()

	if len(m2) != len(m) {
		t.Errorf("bson fields num: %d, map fileds num: %d", len(m2), len(m))
	}

	for k, _ := range m {
		_, exist := m2[k]
		if !exist {
			t.Errorf("bson missing field [%s]", k)
		}
	}
}
