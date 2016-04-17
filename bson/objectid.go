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
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"
)

// ObjectId is a unique ID identifying a BSON value.
// It must be exactly 12 bytes long.
// SequoiaDB objects by default have such a property set in their "_id" property.
type ObjectId string

var machineId = getMachineId()

var objectIdCounter uint32 = getRandomUint32()

func getMachineId() []byte {
	var b [3]byte
	id := b[:]
	hostname, err1 := os.Hostname()
	if err1 != nil {
		_, err2 := io.ReadFull(rand.Reader, id)
		if err2 != nil {
			panic(fmt.Errorf("can't get hostname: %v; %v", err1, err2))
		}
		return id
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	return id
}

func getRandomUint32() uint32 {
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("can't get random uint32: %v", err))
	}
	return uint32((uint32(b[0]) << 0) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24))
}

func NewObjectId() ObjectId {
	var b [12]byte
	// timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// machineId, 3 bytes
	b[4] = machineId[0]
	b[5] = machineId[1]
	b[6] = machineId[2]
	// pid, 2 bytes, big endian
	pid := os.Getpid()
	b[7] = byte(pid >> 8)
	b[8] = byte(pid)
	// counter, 3 bytes, big endian
	i := atomic.AddUint32(&objectIdCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return ObjectId(b[:])
}

func (id ObjectId) IsValid() bool {
	return len(id) == 12
}

func (id ObjectId) String() string {
	return fmt.Sprintf(`{"$oid":"%x"}`, string(id))
}

func (id ObjectId) Hex() string {
	return hex.EncodeToString([]byte(id))
}
