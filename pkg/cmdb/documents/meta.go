// Copyright 2023 NJWS Inc.
// Copyright 2022 Listware

package documents

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// DocumentID is a document ID
// Consists of two parts - collection name and key
type DocumentID string

// Topic representation of DocumentID
func (id DocumentID) Topic() string {
	return strings.Replace(id.String(), "/", ".", -1)
}

// Parent - return parent's Document ID
func (id DocumentID) Parent() DocumentID {
	return DocumentID(path.Dir(id.String()))
}

// Validate validates the given id
func (id DocumentID) Validate() error {
	if id == "" {
		return fmt.Errorf("DocumentID is empty")
	}
	parts := strings.Split(string(id), "/")
	if len(parts) < 2 {
		return fmt.Errorf("Expected 'collection/key[/profile_name]', got '%s'", string(id))
	}
	if parts[0] == "" {
		return fmt.Errorf("Collection part of '%s' is empty", string(id))
	}
	if parts[1] == "" {
		return fmt.Errorf("Key part of '%s' is empty", string(id))
	}
	return nil
}

// pathUnescape unescapes the given value for use in a URL path.
func pathUnescape(s string) string {
	r, _ := url.QueryUnescape(s)
	return r
}

// Collection returns the collection part of the ID.
func (id DocumentID) Collection() string {
	parts := strings.Split(string(id), "/")
	return pathUnescape(parts[0])
}

// Key returns the key part of the ID.
func (id DocumentID) Key() string {
	parts := strings.Split(string(id), "/")
	if len(parts) >= 2 {
		return pathUnescape(parts[1])
	}
	return ""
}

// ProfileName returns the profile name part of the ID.
func (id DocumentID) ProfileName() string {
	parts := strings.Split(string(id), "/")
	if len(parts) >= 3 {
		return pathUnescape(parts[2])
	}
	return ""
}

func (id DocumentID) String() string {
	return string(id)
}
