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

type Doc []DocElement

type DocElement struct {
	Name  string
	Value interface{}
}

func (d Doc) toBson(b *Bson) {
	for _, item := range d {
		b.Append(item.Name, item.Value)
	}
}

func (d Doc) Bson() *Bson {
	b := NewBson()
	d.toBson(b)
	b.Finish()
	return b
}

func (d Doc) String() string {
	return d.Bson().String()
}
