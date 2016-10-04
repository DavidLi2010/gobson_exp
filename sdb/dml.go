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

func buildDeleteMsg(cl string, condition *bson.Doc, hint *bson.Doc) *DeleteMsg {
	var msg DeleteMsg
	msgLen := msg.FixedSize()

	msg.OpCode = DeleteReqMsg
	msg.NameLength = int32(len(cl))
	msg.Name = []byte(cl)
	msgLen += alignedSize(msg.NameLength+1, 4)

	if condition != nil {
		msg.Condition = condition.Bson()
	} else {
		msg.Condition = emptyBson
	}

	msgLen += alignedSize(int32(msg.Condition.Length()), 4)

	if hint != nil {
		msg.Hint = hint.Bson()
	} else {
		msg.Hint = emptyBson
	}

	msgLen += alignedSize(int32(msg.Hint.Length()), 4)

	msg.Length = msgLen
	return &msg
}

func (conn *Conn) Delete(cl string, condition *bson.Doc, hint *bson.Doc) error {
	msg := buildDeleteMsg(cl, condition, hint)

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

func buildUpdateMsg(cl string, flag int32, rule bson.Doc, condition, hint *bson.Doc) *UpdateMsg {
	var msg UpdateMsg
	msgLen := msg.FixedSize()

	msg.OpCode = UpdateReqMsg
	msg.Flags = flag
	msg.NameLength = int32(len(cl))
	msg.Name = []byte(cl)
	msgLen += alignedSize(msg.NameLength+1, 4)

	msg.Rule = rule.Bson()

	msgLen += alignedSize(int32(msg.Rule.Length()), 4)

	if condition != nil {
		msg.Condition = condition.Bson()
	} else {
		msg.Condition = emptyBson
	}

	msgLen += alignedSize(int32(msg.Condition.Length()), 4)

	if hint != nil {
		msg.Hint = hint.Bson()
	} else {
		msg.Hint = emptyBson
	}

	msgLen += alignedSize(int32(msg.Hint.Length()), 4)

	msg.Length = msgLen
	return &msg
}

func (conn *Conn) update(cl string, flag int32, rule bson.Doc, condition, hint *bson.Doc) error {
	msg := buildUpdateMsg(cl, flag, rule, condition, hint)

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

func (conn *Conn) Update(cl string, rule bson.Doc, condition, hint *bson.Doc) error {
	return conn.update(cl, 0, rule, condition, hint)
}

const updateFlagUpsert = 0x00000001

func (conn *Conn) Upsert(cl string, rule bson.Doc, condition, hint, setOnInsert *bson.Doc) error {
	var newHint bson.Doc

	if hint != nil {
		newHint = *hint
	} else {
		newHint = bson.Doc{}
	}

	if setOnInsert != nil {
		newHint = append(newHint, bson.DocElement{"$SetOnInsert", *setOnInsert})
	}

	return conn.update(cl, updateFlagUpsert, rule, condition, &newHint)
}

/*
func (conn *Conn) Count(csName, clName string, condition *bson.Doc) (int64, error) {

}

func (conn *Conn) Find(csName, clName string) *Query {

}
*/
