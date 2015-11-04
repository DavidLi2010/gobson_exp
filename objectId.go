package bson

// ObjectId is a unique ID identifying a BSON value. It must be exactly 12 bytes
// long. SequoiaDB objects by default have such a property set in their "_id"
// property.
type ObjectId string

// Valid returns true if id is valid. A valid id must contain exactly 12 bytes.
func (id ObjectId) Valid() bool {
	return len(id) == 12
}