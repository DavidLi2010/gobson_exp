package bson_test

import (
	"testing"

	. "github.com/DavidLi2010/gobson_exp/bson"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBsonType(t *testing.T) {
	Convey("test value of BsonType", t, func() {
		So(BsonTypeEOD, ShouldEqual, BsonType(0x00))
		So(BsonTypeDouble, ShouldEqual, BsonType(0x01))
		So(BsonTypeString, ShouldEqual, BsonType(0x02))
		So(BsonTypeDoc, ShouldEqual, BsonType(0x03))
		So(BsonTypeArray, ShouldEqual, BsonType(0x04))
		So(BsonTypeBinary, ShouldEqual, BsonType(0x05))
		So(BsonTypeUndefined, ShouldEqual, BsonType(0x06))
		So(BsonTypeObjectId, ShouldEqual, BsonType(0x07))
		So(BsonTypeBool, ShouldEqual, BsonType(0x08))
		So(BsonTypeDate, ShouldEqual, BsonType(0x09))
		So(BsonTypeNull, ShouldEqual, BsonType(0x0A))
		So(BsonTypeRegEx, ShouldEqual, BsonType(0x0B))
		So(BsonTypeDBPointer, ShouldEqual, BsonType(0x0C))
		So(BsonTypeCode, ShouldEqual, BsonType(0x0D))
		So(BsonTypeSymbol, ShouldEqual, BsonType(0x0E))
		So(BsonTypeCodeWScope, ShouldEqual, BsonType(0x0F))
		So(BsonTypeInt32, ShouldEqual, BsonType(0x10))
		So(BsonTypeTimestamp, ShouldEqual, BsonType(0x11))
		So(BsonTypeInt64, ShouldEqual, BsonType(0x12))
		So(BsonTypeMinKey, ShouldEqual, BsonType(0xFF))
		So(BsonTypeMaxKey, ShouldEqual, BsonType(0x7F))
	})
}

func TestBinaryType(t *testing.T) {
	Convey("test value of BinaryType", t, func() {
		So(BinaryTypeGeneral, ShouldEqual, BinaryType(0x00))
		So(BinaryTypeFunction, ShouldEqual, BinaryType(0x01))
		So(BinaryTypeBinaryDeprecated, ShouldEqual, BinaryType(0x02))
		So(BinaryTypeUUIDDeprecated, ShouldEqual, BinaryType(0x03))
		So(BinaryTypeUUID, ShouldEqual, BinaryType(0x04))
		So(BinaryTypeMD5, ShouldEqual, BinaryType(0x05))
		So(BinaryTypeUser, ShouldEqual, BinaryType(0x80))
	})
}
