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

import "github.com/davidli2010/gobson_exp/bson"

var emptyBson = bson.NewBsonBuilder().Finish().Bson()

func buildCmdMsg(cmd string, objs ...bson.Doc) *QueryMsg {
	var msg QueryMsg
	msgLen := msg.FixedSize()

	msg.OpCode = QueryReqMsg
	msg.NameLength = int32(len(cmd))
	msg.Name = []byte(cmd)
	msg.SkipNum = -1
	msg.ReturnNum = -1

	msgLen += alignedSize(msg.NameLength+1, 4)

	if len(objs) > 4 {
		panic("four docs at most")
	}

	if len(objs) > 0 {
		msg.Where = objs[0].Bson()
	} else {
		msg.Where = emptyBson
	}
	msgLen += alignedSize(int32(msg.Where.Length()), 4)

	if len(objs) > 1 {
		msg.Select = objs[1].Bson()
	} else {
		msg.Select = emptyBson
	}
	msgLen += alignedSize(int32(msg.Select.Length()), 4)

	if len(objs) > 2 {
		msg.OrderBy = objs[2].Bson()
	} else {
		msg.OrderBy = emptyBson
	}
	msgLen += alignedSize(int32(msg.OrderBy.Length()), 4)

	if len(objs) > 3 {
		msg.Hint = objs[3].Bson()
	} else {
		msg.Hint = emptyBson
	}
	msgLen += alignedSize(int32(msg.Hint.Length()), 4)

	msg.Length = msgLen
	return &msg
}

func alignedSize(original, bytes int32) int32 {
	size := original + bytes - 1
	size -= size % bytes
	return size
}

const cmdNameCreateCS = "$create collectionspace"

type cmdCreateCS struct {
	Name    string
	Options *bson.Doc
}

func (c *cmdCreateCS) buildMsg() *QueryMsg {
	var doc bson.Doc
	doc = append(doc, bson.DocElement{"Name", c.Name})
	if c.Options != nil {
		doc = append(doc, bson.DocElement{"Options", *c.Options})
	}

	return buildCmdMsg(cmdNameCreateCS, doc)
}
