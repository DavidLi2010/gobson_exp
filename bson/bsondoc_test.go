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

func TestDoc(t *testing.T) {
	expected := `{"outer":"hello", "obj":{"inner":"world"}, "array":["hello world", 123.456], "uintptr":5000000000}`

	doc := bson.Doc{
		{"outer", "hello"},
		{"obj", bson.Doc{{"inner", "world"}}},
		{"array", []interface{}{"hello world", 123.456}},
		{"uintptr", uintptr(5000000000)},
	}

	if expected != doc.String() {
		t.Errorf("append bson/array error, expected:%s, actual:%s", expected, doc.String())
	}

	doc2 := doc.Bson().Doc()
	if len(doc) != len(doc2) {
		t.Errorf("bson fields num: %d, doc fileds num: %d", len(doc2), len(doc))
	}

	for i, v := range doc {
		if doc2[i].Name != v.Name {
			t.Errorf("bson missing field [%s]", v.Name)
		}
	}

	if expected != doc2.String() {
		t.Errorf("doc convert bson error, expected:%s, actual:%s", expected, doc2.String())
	}
}
