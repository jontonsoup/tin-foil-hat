package kademlia

// STORE
type StoreRequest struct {
	Sender Contact
	MsgID  ID
	Key    ID
	Value  []byte
}

type StoreResult struct {
	MsgID ID
	Err   error
}

func (k *Kademlia) Store(req StoreRequest, res *StoreResult) error {
	k.updateContact(&req.Sender)
	res.MsgID = req.MsgID
	k.Table[req.MsgID] = req.Value
	res.Err = nil
	return nil
}
