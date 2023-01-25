// Copyright 2022 Listware

package documents

// EdgeDocument is a minimal document for use in edge collection.
type EdgeDocument struct {
	BaseDocument
	From DocumentID `json:"_from,omitempty"`
	To   DocumentID `json:"_to,omitempty"`
	Type string     `json:"_type,omitempty"`
}
