package ananke

func (n *Node) MarshalBinary() (data []byte, err error) {
	return n.marshalBinaryLeaf()
}

func (n *Node) UnmarshalBinary(data []byte) error {
	return n.unmarshalBinaryLeaf(data)
}
