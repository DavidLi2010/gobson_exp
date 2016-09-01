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
	"fmt"
	"testing"
)

func TestNewConnection(t *testing.T) {
	conn, err := NewConnection("192.168.100.54:11810")
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("conn.order=%v\n", conn.order)
	fmt.Printf("conn.osType=%v\n", conn.osType)

	if err := conn.CreateCS("foo", nil); err != nil {
		t.Error(err)
	}

	if err := conn.Close(); err != nil {
		t.Error(err)
	}
}
