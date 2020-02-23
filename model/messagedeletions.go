package model

// MessageDeletions holds information about deleted messages.
type MessageDeletions struct {
	Messages []*Message
}

// ToExternal converts the event into an external representation.
func (del MessageDeletions) ToExternal() interface{} {
	res := &MessageDeletionsExternal{
		make([]uint, len(del.Messages)),
	}
	for key, value := range del.Messages {
		res.IDs[key] = value.ID
	}
	return res
}

// MessageDeletionsExternal Model
//
// MessageDeletionsExternal holds information about deleted messages.
//
// swagger:model MessageDeletions
type MessageDeletionsExternal struct {
	// The IDs to be deleted.
	//
	// read only: true
	// required: true
	// example: [14,15,16]
	IDs []uint `json:"ids"`
}
