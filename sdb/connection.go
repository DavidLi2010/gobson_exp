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
	"bytes"
	"encoding/binary"
	"errors"
	"net"

	"time"

	"fmt"

	"github.com/davidli2010/gobson_exp/bson"
)

type Conn struct {
	host      string
	conn      net.Conn
	order     binary.ByteOrder
	osType    int32
	buf       bytes.Buffer
	requestId uint64
}

func Connect(host string) (*Conn, error) {
	addr, addrErr := net.ResolveTCPAddr("tcp", host)
	if addrErr != nil {
		return nil, addrErr
	}

	conn, dialErr := net.DialTCP("tcp", nil, addr)
	if dialErr != nil {
		return nil, dialErr
	}

	if err := conn.SetNoDelay(true); err != nil {
		return nil, err
	}

	if err := conn.SetKeepAlive(true); err != nil {
		return nil, err
	}

	if err := conn.SetKeepAlivePeriod(5 * time.Second); err != nil {
		return nil, err
	}

	sysInfoReq := NewSysInfoRequest()
	if err := sysInfoReq.Encode(conn, bson.GetByteOrder()); err != nil {
		return nil, err
	}

	sysInfoReply := SysInfoReply{}
	if err := sysInfoReply.Decode(conn, bson.GetByteOrder()); err != nil {
		return nil, err
	}

	var order binary.ByteOrder
	var osType int32
	if sysInfoReply.EyeCatcher == sysInfoEyeCatcher {
		order = bson.GetByteOrder()
		osType = sysInfoReply.OSType
	} else if sysInfoReply.EyeCatcher == sysInfoEyeCatcherRevert {
		order = bson.GetReverseByteOrder()
		osType = bson.RevertInt32(sysInfoReply.OSType)
	} else {
		return nil, errors.New("Invalid eyecatcher")
	}

	return &Conn{
		host:   host,
		conn:   conn,
		order:  order,
		osType: osType,
		buf:    bytes.Buffer{},
	}, nil
}

func (conn *Conn) Close() error {
	return conn.conn.Close()
}

func (conn *Conn) CreateCS(name string, options *bson.Doc) error {
	cmd := cmdCreateCS{name, options}
	msg := cmd.buildMsg()
	if err := msg.Encode(conn.conn, conn.order); err != nil {
		return err
	}

	var rsp ReplyMsg
	if err := rsp.Decode(conn.conn, conn.order); err != nil {
		return err
	}

	if rsp.Flags != 0 {
		return fmt.Errorf("failed to create cs, rc=%d", rsp.Flags)
	}

	return nil
}
