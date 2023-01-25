// Copyright 2022 Listware

package documents

import (
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
)

// BaseDocument is a minimal document for use in non-edge collection.
type BaseDocument struct {
	Key string     `json:"_key,omitempty"`
	ID  DocumentID `json:"_id,omitempty"`
	Rev string     `json:"_rev,omitempty"`
}

func NewBaseDocument(meta *pbcmdb.Meta) *BaseDocument {
	return &BaseDocument{
		Key: meta.GetKey(),
		ID:  DocumentID(meta.GetId()),
		Rev: meta.GetRev(),
	}
}
