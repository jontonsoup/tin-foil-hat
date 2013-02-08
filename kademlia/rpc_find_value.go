package kademlia

// FIND_VALUE
type FindValueRequest struct {
	Sender Contact
	MsgID  ID
	Key    ID
}

// If Value is nil, it should be ignored, and Nodes means the same as in a
// FindNodeResult.
type FindValueResult struct {
	MsgID ID
	Value []byte
	Nodes []FoundNode
	Err   error
}

func (k *Kademlia) FindValue(req FindValueRequest, res *FindValueResult) error {
	// TODO: Implement.
	return nil
}
