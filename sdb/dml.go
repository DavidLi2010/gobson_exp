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

func buildInsertMsg(cl string, doc bson.Doc) *InsertMsg {
	var msg InsertMsg
	msgLen := msg.FixedSize()

	msg.OpCode = InsertReqMsg
	msg.NameLength = int32(len(cl))
	msg.Name = []byte(cl)
	msgLen += alignedSize(msg.NameLength+1, 4)

	msg.Doc = doc.Bson()
	msgLen += alignedSize(int32(msg.Doc.Length()), 4)

	msg.Length = msgLen
	return &msg
}

func (conn *Conn) Insert(cl string, doc bson.Doc) error {
	msg := buildInsertMsg(cl, doc)

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

/*
func (conn *Conn) Delete(csName, clName string, condition *bson.Doc, hint *bson.Doc) error {

}

func (conn *Conn) Truncate(csName, clName string) error {

}

func (conn *Conn) Update(csName, clName string, rule *bson.Doc, condition *bson.Doc, hint *bson.Doc) error {

}

func (conn *Conn) Upsert(csName, clName string, rule *bson.Doc, condition *bson.Doc, hint *bson.Doc, set *bson.Doc) error {

}

func (conn *Conn) Count(csName, clName string, condition *bson.Doc) (int64, error) {

}

func (conn *Conn) Find(csName, clName string) *Query {

}
*/