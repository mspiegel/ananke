package ananke

const NoSuchNodeError = "node %d does not exist"

type Storage interface {
	Write(node Node) error
	Read(id int64) (Node, error)
}
