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

	"github.com/davidli2010/gobson_exp/bson"
)

func (conn *Conn) runCmd(cmd Cmd) error {
	msg := cmd.buildMsg()

	if err := msg.Encode(conn.conn, conn.order); err != nil {
		return err
	}

	var rsp ReplyMsg
	if err := rsp.Decode(conn.conn, conn.order); err != nil {
		return err
	}

	if rsp.Flags != 0 {
		return fmt.Errorf("error=%s,rc=%d",
			rsp.Error, rsp.Flags)
	}

	return nil
}

func (conn *Conn) CreateCS(name string, options *bson.Doc) error {
	cmd := &cmdCreateCS{name, options}
	return conn.runCmd(cmd)
}

func (conn *Conn) DropCS(name string) error {
	cmd := &cmdDropCS{name}
	return conn.runCmd(cmd)
}

func (conn *Conn) CreateCL(csName, clName string, options *bson.Doc) error {
	cmd := &cmdCreateCL{csName, clName, options}
	return conn.runCmd(cmd)
}

func (conn *Conn) DropCL(csName, clName string) error {
	cmd := &cmdDropCL{csName, clName}
	return conn.runCmd(cmd)
}

func (conn *Conn) CreateIndex(csName, clName, indexName string, indexDefine bson.Doc, options *bson.Doc) error {
	cmd := &cmdCreateIndex{csName, clName, indexName, indexDefine, options}
	return conn.runCmd(cmd)
}

func (conn *Conn) DropIndex(csName, clName, indexName string) error {
	cmd := &cmdDropIndex{csName, clName, indexName}
	return conn.runCmd(cmd)
}
