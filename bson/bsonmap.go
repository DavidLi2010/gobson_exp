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

type Map map[string]interface{}

func (m Map) toBsonBuilder(b *BsonBuilder) {
	for name, v := range m {
		b.Append(name, v)
	}
}

func (m Map) Bson() *Bson {
	b := NewBsonBuilder()
	m.toBsonBuilder(b)
	b.Finish()
	return b.Bson()
}

func (m Map) String() string {
	return m.Bson().String()
}
