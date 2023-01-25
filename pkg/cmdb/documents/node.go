// Copyright 2022 Listware

package documents

import (
	"encoding/json"

	"git.fg-tech.ru/listware/proto/sdk/pbcmdb/pbqdsl"
)

type Node struct {
	Id     DocumentID      `json:"id,omitempty"`
	LinkId DocumentID      `json:"link_id,omitempty"`
	Key    string          `json:"key,omitempty"`
	Name   string          `json:"name,omitempty"`
	Type   string          `json:"type,omitempty"`
	Object json.RawMessage `json:"object,omitempty"`
	Link   json.RawMessage `json:"link,omitempty"`
	Path   json.RawMessage `json:"path,omitempty"`
}

type Nodes []*Node

func NewNodes() Nodes {
	return make([]*Node, 0)
}

func (nodes *Nodes) AddElement(element *pbqdsl.Element) {
	*nodes = append(*nodes, NewNode(element))
}
func (nodes *Nodes) Add(node ...*Node) {
	*nodes = append(*nodes, node...)
}

func NewNode(element *pbqdsl.Element) *Node {
	return &Node{
		Id:     DocumentID(element.Id),
		Key:    element.Key,
		Name:   element.Name,
		Type:   element.Type,
		Object: element.Object,
		LinkId: DocumentID(element.LinkId),
		Link:   element.Link,
		Path:   element.Path,
	}
}

func (node *Node) ToElement() *pbqdsl.Element {
	return &pbqdsl.Element{
		Id:     node.Id.String(),
		Key:    node.Key,
		Name:   node.Name,
		Type:   node.Type,
		Object: node.Object,
		Link:   node.Link,
		LinkId: node.LinkId.String(),
		Path:   node.Path,
	}
}
