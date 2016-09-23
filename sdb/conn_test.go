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

package sdb

import (
	"testing"

	"github.com/davidli2010/gobson_exp/bson"
)

func TestNewConnection(t *testing.T) {
	conn, err := Connect("192.168.100.53:11810")
	if err != nil {
		t.Error(err)
	}

	cs := "foo"
	cl := "bar"

	if err := conn.CreateCS(cs, nil); err != nil {
		t.Error(err)
	}

	if err := conn.CreateCL(cs, cl, nil); err != nil {
		t.Error(err)
	}

	indexOptions := bson.Doc{
		{"unique", true},
		{"enforced", true},
		{"sorBufferSize", 128},
	}
	if err := conn.CreateIndex(cs, cl, "a_idx", bson.Doc{{"a", 1}}, &indexOptions); err != nil {
		t.Error(err)
	}

	if err := conn.Insert(cs+"."+cl, bson.Doc{{"a", 123}}); err != nil {
		t.Error(err)
	}

	if err := conn.Insert(cs+"."+cl, bson.Doc{{"a", 456}}); err != nil {
		t.Error(err)
	}

	if err := conn.Delete(cs+"."+cl, nil, nil); err != nil {
		t.Error(err)
	}

	if err := conn.DropIndex(cs, cl, "a_idx"); err != nil {
		t.Error(err)
	}

	if err := conn.DropCL(cs, cl); err != nil {
		t.Error(err)
	}

	if err := conn.DropCS(cs); err != nil {
		t.Error(err)
	}

	if err := conn.Close(); err != nil {
		t.Error(err)
	}
}
