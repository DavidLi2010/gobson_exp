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
		panic("Four docs at most")
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

const (
	cmdNameCreateCS    = "$create collectionspace"
	cmdNameDropCS      = "$drop collectionspace"
	cmdNameCreateCL    = "$create collection"
	cmdNameDropCL      = "$drop collection"
	cmdNameTruncateCL  = "$truncate"
	cmdNameCreateIndex = "$create index"
	cmdNameDropIndex   = "$drop index"
)

type Cmd interface {
	buildMsg() *QueryMsg
}

type cmdCreateCS struct {
	Name    string
	Options *bson.Doc
}

func (c *cmdCreateCS) buildMsg() *QueryMsg {
	doc := bson.Doc{
		{"Name", c.Name},
	}
	if c.Options != nil {
		doc = append(doc, *c.Options...)
	}

	return buildCmdMsg(cmdNameCreateCS, doc)
}

type cmdDropCS struct {
	Name string
}

func (c *cmdDropCS) buildMsg() *QueryMsg {
	doc := bson.Doc{
		{"Name", c.Name},
	}
	return buildCmdMsg(cmdNameDropCS, doc)
}

type cmdCreateCL struct {
	CSName  string
	CLName  string
	Options *bson.Doc
}

func (c *cmdCreateCL) buildMsg() *QueryMsg {
	fullName := c.CSName + "." + c.CLName
	doc := bson.Doc{
		{"Name", fullName},
	}
	if c.Options != nil {
		doc = append(doc, *c.Options...)
	}

	return buildCmdMsg(cmdNameCreateCL, doc)
}

type cmdDropCL struct {
	CSName string
	CLName string
}

func (c *cmdDropCL) buildMsg() *QueryMsg {
	fullName := c.CSName + "." + c.CLName
	doc := bson.Doc{
		{"Name", fullName},
	}
	return buildCmdMsg(cmdNameDropCL, doc)
}

type cmdCreateIndex struct {
	CSName      string
	CLName      string
	IndexName   string
	IndexDefine bson.Doc
	Options     *bson.Doc
}

func (c *cmdCreateIndex) buildMsg() *QueryMsg {
	fullName := c.CSName + "." + c.CLName

	var m bson.Map
	if c.Options != nil {
		m = c.Options.Map()
	}

	var unique bool
	if v, exist := m["unique"]; exist {
		unique = v.(bool)
	}

	var enforced bool
	if v, exist := m["enforced"]; exist {
		enforced = v.(bool)
	}

	index := bson.Doc{
		{"name", c.IndexName},
		{"key", c.IndexDefine},
		{"unique", unique},
		{"enforced", enforced},
	}

	doc := bson.Doc{
		{"Collection", fullName},
		{"Index", index},
	}

	hint := bson.Doc{}
	if v, exist := m["sortBufferSize"]; exist {
		hint = append(hint, bson.DocElement{"SortBufferSize", v.(int)})
	}

	return buildCmdMsg(cmdNameCreateIndex, doc, bson.Doc{}, bson.Doc{}, hint)
}

type cmdDropIndex struct {
	CSName    string
	CLName    string
	IndexName string
}

func (c *cmdDropIndex) buildMsg() *QueryMsg {
	fullName := c.CSName + "." + c.CLName
	index := bson.Doc{
		{"", c.IndexName},
	}
	doc := bson.Doc{
		{"Collection", fullName},
		{"Index", index},
	}
	return buildCmdMsg(cmdNameDropIndex, doc)
}

type cmdTruncateCL struct {
	CSName string
	CLName string
}

func (c *cmdTruncateCL) buildMsg() *QueryMsg {
	fullName := c.CSName + "." + c.CLName
	doc := bson.Doc{
		{"Collection", fullName},
	}
	return buildCmdMsg(cmdNameTruncateCL, doc)
}
